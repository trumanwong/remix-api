package main

import (
	"remix-api/api/console/commands"
	"github.com/robfig/cron/v3"
)

func main() {
	crontab := cron.New(cron.WithSeconds())
	// 每分钟执行一次
	crontab.AddFunc("0 */10 * * * *", commands.ClearFiles)
	crontab.Start()
	select {}
}
