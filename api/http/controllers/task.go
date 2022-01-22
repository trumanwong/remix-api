package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/trumanwong/go-internal/util"
	"net/http"
	"path"
	"remix-api/configs"
	"remix-api/internal/logs"
	"remix-api/internal/mq"
	"remix-api/models"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type TaskController struct {
	Controller
}

// Create 创建任务
func (this *TaskController) Create(ctx *gin.Context) {
	err := this.construct(ctx)
	if err != nil {
		util.Response(ctx, nil, http.StatusForbidden, "Token Invalid")
		return
	}
	task := new(models.Task)
	task.Type = this.Type
	task.Other = map[string]interface{}{}

	// 文件存储路径
	filePath := fmt.Sprintf(
		"/upload/hotlink-ok/%s/%s",
		util.MD5(fmt.Sprintf("%v", this.User.ID))[4:8],
		time.Now().Format("2006/01/02"),
	)
	storePath := fmt.Sprintf("%s%s", configs.Config.Other.StorePath, filePath)
	if !util.PathExists(storePath) {
		err = util.CreatePath(storePath)
		if err != nil {
			logs.Write(logs.Log{
				Message:   err.Error(),
				Level:     logs.LogLevelError,
				CreatedAt: time.Now(),
			})
			util.Response(ctx, nil, http.StatusInternalServerError, "Server Error")
			return
		}
	}
	if this.Type == models.TaskTypeGifRevert {
		file, err := ctx.FormFile("file")
		if err != nil {
			util.Response(ctx, nil, http.StatusBadRequest, "No file")
			return
		}
		ext := path.Ext(file.Filename)
		if file.Size > 10*1024*1024 ||
			!util.InArray(
				ext,
				strings.Split(configs.Config.Other.AllowPictureExt, ","),
			) {
			util.Response(ctx, nil, http.StatusBadRequest, "Invalid file")
			return
		}
		task.FileName = file.Filename
		task.Size = uint64(file.Size)
		fileName := util.MD5(util.GenerateSnowFlakeID(1))
		task.Path = fmt.Sprintf("%s/%s.%s", filePath, fileName, ext)
		task.ConvertPath = fmt.Sprintf("%s/%s_covert.%s", filePath, fileName, ext)
		err = ctx.SaveUploadedFile(file, fmt.Sprintf("%s%s", configs.Config.Other.StorePath, task.Path))
		if err != nil {
			logs.Write(logs.Log{
				Message:   err.Error(),
				Level:     logs.LogLevelError,
				CreatedAt: time.Now(),
			})
			util.Response(ctx, nil, http.StatusInternalServerError, "Server Error")
			return
		}
	} else if task.Type == models.TaskTypeNokiaSms {
		text := strings.Trim(ctx.PostForm("text"), " ")
		if len(text) == 0 {
			util.Response(ctx, nil, http.StatusForbidden, "请输入短信文字")
			return
		}
		if utf8.RuneCountInString(text) > 30 {
			util.Response(ctx, nil, http.StatusBadRequest, "短信文字最大长度不能超过30")
			return
		}
		task.ConvertPath = fmt.Sprintf("%s/%s_covert.png", filePath, util.MD5(util.GenerateSnowFlakeID(1)))
		task.Other = map[string]interface{}{
			"text": text,
		}
	} else if task.Type == models.TaskTypeRemix {
		// 鬼畜
		texts := strings.Split(ctx.PostForm("texts"), ",")
		remixType, err := strconv.ParseUint(ctx.PostForm("remix_type"), 10, 8)
		arr := append([]string{}, configs.Config.Task.Remix.Sentences[remixType]...)
		if err != nil || int(remixType) >= len(configs.Config.Task.Remix.Sentences) || len(texts) != len(arr) {
			util.Response(ctx, nil, http.StatusBadRequest, "Invalid Params")
			return
		}
		for i := 0; i < len(arr); i++ {
			text := strings.Trim(texts[i], " ")
			text = strings.ReplaceAll(text, ",", " ")
			if text == "" {
				continue
			}
			if utf8.RuneCountInString(text) > 20 {
				util.Response(ctx, nil, http.StatusBadRequest, "每条最多只能输入20个字符")
				return
			}
			arr[i] = text
		}
		task.ConvertPath = fmt.Sprintf("%s/%s_covert.gif", filePath, util.MD5(util.GenerateSnowFlakeID(1)))
		task.Other = map[string]interface{}{
			"text":       strings.Join(arr, ","),
			"remix_type": remixType,
		}
	}
	task.UUID = this.UUID
	task.Serial = util.GenerateUUID()
	task.Status = models.TaskStatusProcessing
	task.UserID = this.User.ID
	task.IP = util.IP2long(ctx.ClientIP())
	err = task.Create()
	if err != nil {
		logs.Write(logs.Log{
			Message:   err.Error(),
			Level:     logs.LogLevelError,
			CreatedAt: time.Now(),
		})
		util.Response(ctx, nil, http.StatusInternalServerError, "Server Error")
		return
	}
	err = mq.RabbitMQ.NewWorkQueue(mq.QueueTask, []byte(fmt.Sprintf("%v", task.ID)))
	if err != nil {
		util.Response(ctx, nil, http.StatusInternalServerError, "Server Error")
		return
	}
	util.Response(ctx, task, http.StatusOK, "Success")
}

// Show 查询任务
func (this *TaskController) Show(ctx *gin.Context) {
	err := this.construct(ctx)
	if err != nil {
		util.Response(ctx, nil, http.StatusForbidden, "Token Invalid")
		return
	}

	serial := strings.Trim(ctx.Query("serial"), " ")
	if len(serial) == 0 {
		util.Response(ctx, nil, http.StatusForbidden, "Invalid Params")
		return
	}
	task, err := models.NewTask(map[string]interface{}{
		"serial":  serial,
		"user_id": this.User.ID,
		"uuid":    this.UUID,
	})
	if err != nil {
		util.Response(ctx, nil, http.StatusNotFound, "Not found")
		return
	}
	util.Response(ctx, task, http.StatusOK, "Success")
}

// Download 任务文件下载
func (this *TaskController) Download(ctx *gin.Context) {
	err := this.construct(ctx)
	if err != nil {
		util.Response(ctx, nil, http.StatusForbidden, "Token Invalid")
		return
	}

	serial := strings.Trim(ctx.Query("serial"), " ")
	if len(serial) == 0 {
		util.Response(ctx, nil, http.StatusForbidden, "Invalid Params")
		return
	}
	task, err := models.NewTask(map[string]interface{}{
		"serial":  serial,
		"user_id": this.User.ID,
		"uuid":    this.UUID,
	})
	if err != nil {
		util.Response(ctx, nil, http.StatusNotFound, "Not found")
		return
	}
	realConvertPath := fmt.Sprintf(
		"%s%s",
		configs.Config.Other.StorePath,
		task.ConvertPath,
	)
	if task.Status != models.TaskStatusSuccess || len(task.ConvertPath) == 0 || !util.PathExists(realConvertPath) {
		util.Response(ctx, nil, http.StatusNotFound, "Not Found")
		return
	}
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", "attachment; filename="+path.Base(realConvertPath))
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Cache-Control", "no-cache")
	ctx.File(realConvertPath)
}

// GetList 获取任务列表
func (this *TaskController) GetList(ctx *gin.Context) {
	err := this.construct(ctx)
	if err != nil {
		util.Response(ctx, nil, http.StatusForbidden, "Token Invalid")
		return
	}

	params := this.getAllQueryParams(ctx)
	params["type"] = this.Type
	task := new(models.Task)
	util.Response(ctx, util.PaginateData(task.GetList(params), params), http.StatusOK, "Success")
}

// Delete 删除任务
func (this *TaskController) Delete(ctx *gin.Context) {
	err := this.construct(ctx)
	if err != nil {
		util.Response(ctx, nil, http.StatusForbidden, "Token Invalid")
		return
	}

	serial := strings.Trim(ctx.Query("serial"), " ")
	if len(serial) == 0 {
		util.Response(ctx, nil, http.StatusForbidden, "Invalid Params")
		return
	}
	task, err := models.NewTask(map[string]interface{}{
		"serial":  serial,
		"user_id": this.User.ID,
		"uuid":    this.UUID,
	})
	if err != nil {
		util.Response(ctx, nil, http.StatusNotFound, "Not found")
		return
	}
	err = task.Delete()
	if err != nil {
		util.Response(ctx, nil, http.StatusInternalServerError, "Server Error")
		return
	}
	util.Response(ctx, nil, http.StatusOK, "Success")
}
