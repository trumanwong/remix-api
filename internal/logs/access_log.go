package logs

import (
	"net/http"
	"time"
)

type AccessLog struct {
	Uri         string        `json:"uri"`
	ClientIp    string        `json:"client_ip"`
	Method      string        `json:"method"`
	Header      http.Header   `json:"header"`
	StatusCode  int           `json:"status_code"`
	ExecuteTime time.Duration `json:"execute_time"`
	CreatedAt   time.Time     `json:"created_at"`
}

func (this AccessLog) IndexName() string {
	return "qnw_api_access_logs"
}
