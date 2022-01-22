package main

import (
	"fmt"
	"remix-api/api/routers"
	"remix-api/configs"
	"remix-api/internal/logs"
	"time"
)

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