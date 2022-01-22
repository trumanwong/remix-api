package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/trumanwong/go-internal/util"
	"net/http"
	"os"
	"remix-api/configs"
	"remix-api/models"
)

type Controller struct {
	User *models.User
	UUID string
	Type uint8
}

func (this *Controller) construct(ctx *gin.Context) error {
	value, exists := ctx.Get("uuid")
	if !exists {
		return errors.New("invalid token")
	}
	this.User = new(models.User)
	this.UUID = value.(string)
	value, exists = ctx.Get("type")
	if exists {
		this.Type = uint8(value.(float64))
	}
	return nil
}

// 获取所有get请求参数
func (this *Controller) getAllQueryParams(ctx *gin.Context) map[string]interface{} {
	res := make(map[string]interface{})
	params := ctx.Request.URL.Query()
	for key, param := range params {
		res[key] = param[0]
	}
	res["uuid"] = this.UUID
	return res
}

// GetFiles 文件服务器
func (this *Controller) GetFiles(ctx *gin.Context) {
	fileName := fmt.Sprintf("%s/%s", configs.Config.Other.StorePath, ctx.Query("file"))
	fileInfo, err := os.Stat(fileName)
	if err != nil || fileInfo.IsDir() {
		util.Response(ctx, nil, http.StatusNotFound, "Not Found")
		return
	}
	ctx.File(fileName)
}
