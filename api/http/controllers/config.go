package controllers

import (
	"remix-api/models"
	"github.com/gin-gonic/gin"
	"github.com/trumanwong/go-internal/util"
	"net/http"
)

type ConfigController struct {
	Controller
}

// Show 获取配置
func (this *ConfigController) Show(ctx *gin.Context) {
	config := new(models.Config)
	params := this.getAllQueryParams(ctx)
	params["id"] = ctx.Param("id")
	c, err := config.First(params)
	if err != nil {
		util.Response(ctx, nil, http.StatusNotFound, "Not Found")
		return
	}
	util.Response(ctx, c, http.StatusOK, "操作成功")
}