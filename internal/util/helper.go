package util

import (
	"github.com/trumanwong/go-internal/notice"
	"log"
	"remix-api/configs"
)

// PushWeChatRobot 推送WechatRobot
func PushWeChatRobot(level, message string) {
	_, err := notice.PushWeChatRobot(level, message, configs.Config.Other.WeChatRobotUrl)
	if err != nil {
		log.Println(err)
	}
}
