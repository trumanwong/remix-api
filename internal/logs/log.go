package logs

import "time"

type Log struct {
	Message   string    `json:"message"`
	Level     string    `json:"level"`
	CreatedAt time.Time `json:"created_at"`
}

func (this Log) IndexName() string {
	return "qnw_api_logs"
}
