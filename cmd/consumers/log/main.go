package main

import (
	"remix-api/configs"
	"remix-api/internal/cache"
	"remix-api/internal/logs"
	"remix-api/internal/mq"
	util2 "remix-api/internal/util"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/sirupsen/logrus"
	"github.com/trumanwong/go-internal/es"
	"github.com/trumanwong/go-internal/util"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func init()  {
	configs.Setup()
	cache.Setup()
	mq.Setup()
}

func main() {
	ch, err := mq.RabbitMQ.NewChannel()
	if err != nil {
		util2.PushWeChatRobot("error", fmt.Sprintf("Failed to open a channel, %s", err))
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		mq.QueueLog,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		util2.PushWeChatRobot("error", fmt.Sprintf("Failed to declare a queue, %s", err))
		return
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		util2.PushWeChatRobot("error", fmt.Sprintf("Failed to set QoS, %s", err))
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
		util2.PushWeChatRobot("error", fmt.Sprintf("Failed to register a consumer, %s", err))
		return
	}

	forever := make(chan bool)

	var esClient *elasticsearch.Client
	if configs.Config.Log.Driver == "es" {
		esClient, err = es.NewElasticClient(
			configs.Config.Elasticsearch.UserName,
			configs.Config.Elasticsearch.Password,
			configs.Config.Elasticsearch.Addresses,
		)

		if err != nil {
			util2.PushWeChatRobot("error", fmt.Sprintf("Failed to declare es client, %s", err))
			return
		}
	}
	go func() {
		for msg := range msgs {
			msg.Ack(false)
			m := make(map[string]interface{})
			_ = json.Unmarshal(msg.Body, &m)
			content := m["content"].(map[string]interface{})
			if _, ok := m["file"]; ok {
				content["file"], content["line"] = m["file"], m["line"]
			}
			m["content"] = content
			b, _ := json.Marshal(m["content"])
			level := "info"
			if _, ok := content["level"]; ok {
				level = content["level"].(string)
			}
			if level == logs.LogLevelError {
				util2.PushWeChatRobot("error", string(b))
			}
			switch configs.Config.Log.Driver {
			case "es":
				func() {
					documentId := util.GenerateSnowFlakeID(1)
					// Set up the request object
					req := esapi.IndexRequest{
						Index:      m["index_name"].(string),
						DocumentID: documentId,
						Body:       strings.NewReader(string(b)),
						Refresh:    "true",
					}

					// Perform the request with the client
					res, err := req.Do(context.Background(), esClient)
					if err != nil {
						util2.PushWeChatRobot("error", fmt.Sprintf("Error getting response: %s", err))
						return
					}
					defer res.Body.Close()

					if res.IsError() {
						util2.PushWeChatRobot("error", fmt.Sprintf("[%s] Error indexing document ID=%s, %s,", res.Status(), documentId, res.String()))
					}
				}()
			case "file":
				filename := configs.Config.Log.Path + "/" + time.Now().Format("2006-01-02") + ".log"
				file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
				if err != nil {
					panic(err)
				}
				logger := logrus.New()
				logger.Formatter = &logrus.JSONFormatter{}
				logger.SetReportCaller(true)
				logger.SetOutput(file)
				value, _ := cache.Cache.Get("clear_expire_log_" + filename)
				if len(value) == 0 {
					clearExpireFile()
					cache.Cache.Put("clear_expire_log_"+filename, 1, 86400)
				}
				switch level {
				case logs.LogLevelInfo:
					logger.Info(m["content"])
				case logs.LogLevelWarn:
					logger.Warn(m["content"])
				case logs.LogLevelError:
					logger.Error(m["content"])
				}
			}

		}
	}()

	log.Printf(" [*] Waiting for log. To exit press CTRL+C")
	<-forever
}

// clearExpireFile 清除过期日志文件
func clearExpireFile() {
	layout := time.Now().AddDate(0, 0, -configs.Config.Log.MaxDay).Format("2006-01-02")
	files, err := ioutil.ReadDir(configs.Config.Log.Path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.Name() < layout {
			err := os.Remove(configs.Config.Log.Path + "/" + file.Name())
			if err != nil {
				panic(err)
			}
		}
	}
}
