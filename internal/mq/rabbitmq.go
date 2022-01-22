package mq

import (
	"remix-api/configs"
	"remix-api/internal/util"
	"github.com/trumanwong/go-internal/mq/rabbitmq"
)

const (
	QueueLog  = "blog_api_queue_log"
	QueueTask = "blog_api_queue_task"
)

var RabbitMQ *rabbitmq.RabbitMQ
var err error

func init() {
	RabbitMQ, err = rabbitmq.NewRabbitMQ(configs.Config.RabbitMQ.Url)
	if err != nil {
		util.PushWeChatRobot(
			"Error",
			err.Error(),
		)
	}
}
