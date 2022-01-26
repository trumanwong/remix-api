package main

import (
	"remix-api/api/routers"
	"remix-api/configs"
	"remix-api/internal/cache"
	"remix-api/internal/logs"
	"remix-api/internal/mq"
	"remix-api/models"
	"fmt"
	"time"
)

func init()  {
	configs.Setup()
	cache.Setup()
	mq.Setup()
	models.Setup()
}

func main() {
	router := routers.InitRouter()
	err := router.Run(fmt.Sprintf(":%d", configs.Config.HttpPort))
	if err != nil {
		logs.Write(logs.Log{
			Message:   fmt.Sprintf("启动失败：%s", err),
			Level:     logs.LogLevelError,
			CreatedAt: time.Now(),
		})
	}
}