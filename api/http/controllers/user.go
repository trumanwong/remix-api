package controllers

import (
	"remix-api/configs"
	"remix-api/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/trumanwong/go-internal/util"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UserController struct {
	Controller
}

// MiniProgramLogin 小程序登录
func (this *UserController) MiniProgramLogin(ctx *gin.Context) {
	code := strings.Trim(ctx.PostForm("code"), " ")
	taskType, err := strconv.ParseUint(ctx.PostForm("type"), 10, 8)
	ctx.Request.ParseForm()
	if len(code) == 0 || err != nil || uint8(taskType) < models.TaskTypeGifRevert || uint8(taskType) > models.TaskTypeRemix {
		util.Response(ctx, nil, http.StatusBadRequest, "Invalid params")
		return
	}
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=" + configs.Config.Wechat.MiniProgram.AppId + "&secret=" + configs.Config.Wechat.MiniProgram.AppSecret + "&js_code=" + code + "&grant_type=authorization_code"
	resp, err := http.Get(url)
	if err != nil {
		util.Response(ctx, nil, http.StatusBadRequest, "Invalid params")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		util.Response(ctx, nil, http.StatusBadRequest, "Invalid params")
		return
	}
	result := make(map[string]interface{})
	json.Unmarshal(body, &result)

	if errCode, ok := result["errcode"]; ok && int(errCode.(float64)) != 0 {
		util.Response(ctx, result, http.StatusBadRequest, "Invalid params")
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uuid": util.GenerateUUID(),
		"type": taskType,
		"exp":  time.Now().Unix() + 3600,
	})
	tokenString, err := token.SignedString([]byte(configs.Config.Other.JwtKey))
	util.Response(ctx, tokenString, http.StatusOK, "操作成功")
}
