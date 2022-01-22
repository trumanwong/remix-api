package logs

import (
	"encoding/json"
	"fmt"
	"remix-api/internal/mq"
	"remix-api/internal/util"
	"runtime"
)

type Logger interface {
	IndexName() string
}

const (
	LogLevelInfo  = "info"
	LogLevelWarn  = "Warn"
	LogLevelError = "Error"
)

// Write 推入日志队列
func Write(esLog Logger) {
	m := map[string]interface{}{
		"index_name": esLog.IndexName(),
		"content":    esLog,
	}
	_, file, line, ok := runtime.Caller(1)
	if ok {
		m["file"], m["line"] = file, line
	}
	data, _ := json.Marshal(m)
	err := mq.RabbitMQ.NewWorkQueue(
		mq.QueueLog,
		data,
	)
	if err != nil {
		util.PushWeChatRobot(LogLevelError, fmt.Sprintf("Push write log queue fail, %s, %s", err, string(data)))
	}
}
