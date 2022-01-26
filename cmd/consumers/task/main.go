package main

import (
	"remix-api/configs"
	"remix-api/internal/cache"
	"remix-api/internal/logs"
	"remix-api/internal/mq"
	"remix-api/models"
	"bufio"
	"bytes"
	"fmt"
	"github.com/trumanwong/go-internal/util"
	"html/template"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func init()  {
	configs.Setup()
	cache.Setup()
	mq.Setup()
	models.Setup()
}

func main() {
	ch, err := mq.RabbitMQ.NewChannel()
	if err != nil {
		logs.Write(logs.Log{
			Message:   fmt.Sprintf("Failed to open a channel, %s", err),
			Level:     logs.LogLevelError,
			CreatedAt: time.Now(),
		})
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		mq.QueueTask,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logs.Write(logs.Log{
			Message:   fmt.Sprintf("Failed to declare a queue, %s", err),
			Level:     logs.LogLevelError,
			CreatedAt: time.Now(),
		})
		return
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		logs.Write(logs.Log{
			Message:   fmt.Sprintf("Failed to set QoS, %s", err),
			Level:     logs.LogLevelError,
			CreatedAt: time.Now(),
		})
		return
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logs.Write(logs.Log{
			Message:   fmt.Sprintf("Failed to register a consumer, %s", err),
			Level:     logs.LogLevelError,
			CreatedAt: time.Now(),
		})
		return
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)
			msg.Ack(false)
			id, _ := strconv.ParseInt(string(msg.Body), 10, 64)
			params := map[string]interface{}{
				"id": uint64(id),
			}
			task, err := models.NewTask(params)
			if err != nil {
				logs.Write(logs.Log{
					Message:   fmt.Sprintf("[%v]任务不存在, %s", task.ID, err),
					Level:     logs.LogLevelError,
					CreatedAt: time.Now(),
				})
				continue
			}

			arg := ""
			scriptPath := configs.Config.Command.Python.GifRevertScript
			if task.Type == models.TaskTypeGifRevert {
				arg = fmt.Sprintf("%s/%s", configs.Config.Other.StorePath, task.Path)
				if !util.PathExists(arg) {
					logs.Write(logs.Log{
						Message:   fmt.Sprintf("[%v]文件不存在, %s", task.ID, task.Path),
						Level:     logs.LogLevelError,
						CreatedAt: time.Now(),
					})
					continue
				}
			} else if task.Type == models.TaskTypeNokiaSms {
				arg = fmt.Sprintf("%v", task.Other["text"])
				scriptPath = configs.Config.Command.Python.NokiaSmsScript
			}

			realConvertPath := fmt.Sprintf(
				"%s%s",
				configs.Config.Other.StorePath,
				task.ConvertPath,
			)

			if task.Type == models.TaskTypeRemix {
				scriptPath = configs.Config.Command.Python.RemixScript
				tpl := fmt.Sprintf("%v", task.Other["remix_type"])
				sentences := strings.Split(fmt.Sprintf("%v", task.Other["text"]), ",")
				err = runRemixCommand(tpl, realConvertPath, sentences)
			} else {
				cmd := exec.Command(
					configs.Config.Command.Python.Command,
					scriptPath,
					arg,
					realConvertPath,
				)
				err = cmd.Run()
			}

			status := models.TaskStatusSuccess
			if err != nil || !util.PathExists(realConvertPath) {
				status = models.TaskStatusFail
				logs.Write(logs.Log{
					Message:   fmt.Sprintf("[%v]转换失败, %s", task.ID, err),
					Level:     logs.LogLevelError,
					CreatedAt: time.Now(),
				})
			}
			err = task.Update(models.Task{
				Status: status,
			})
			if err != nil {
				logs.Write(logs.Log{
					Message:   fmt.Sprintf("[%v]更新数据失败, %s", task.ID, err),
					Level:     logs.LogLevelError,
					CreatedAt: time.Now(),
				})
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func runRemixCommand(tpl, realConvertPath string, sentences []string) error {
	tplPath := fmt.Sprintf(
		"%s/%s/template.tpl",
		configs.Config.Task.Remix.TemplatePath,
		tpl,
	)
	videoPath := fmt.Sprintf(
		"%s/%s/template.mp4",
		configs.Config.Task.Remix.TemplatePath,
		tpl,
	)

	// 生成ass
	arr := strings.Split(realConvertPath, "/")
	assPath := strings.Join(arr[0:len(arr)-1], "/") + "/template.ass"
	err := generateAss(tplPath, assPath, sentences)

	if err != nil {
		return err
	}

	cmd := exec.Command(
		configs.Config.Command.Ffmpeg,
		"-i",
		videoPath,
		"-r",
		"8",
		"-vf",
		fmt.Sprintf("ass=%s,scale=300:-1", assPath),
		"-y",
		realConvertPath,
	)
	err = cmd.Run()
	// 删除ass
	os.Remove(assPath)
	return err
}

func generateAss(tplPath, assPath string, sentences []string) error {
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		return err
	}
	var b bytes.Buffer
	err = tpl.Execute(&b, sentences)
	if err != nil {
		return err
	}
	f, err := os.Create(assPath)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	_, err = w.WriteString(b.String())
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}
