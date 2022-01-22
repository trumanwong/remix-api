package commands

import (
	"remix-api/configs"
	"remix-api/internal/logs"
	"remix-api/models"
	"fmt"
	"github.com/trumanwong/go-internal/util"
	"log"
	"os"
	"time"
)

func ClearFiles() {
	log.Println("clear file processing...")
	task := new(models.Task)
	params := map[string]interface{}{
		"end_at": time.Now().Add(-time.Hour),
		"export": "1",
	}

	result := task.GetList(params)
	tasks := result["list"].([]models.Task)
	var err error
	for i := 0; i < len(tasks); i++ {
		filePath := fmt.Sprintf("%s%s", configs.Config.Other.StorePath, tasks[i].Path)
		if tasks[i].Path != "" && task.Type == models.TaskTypeGifRevert && util.PathExists(filePath) {
			err = os.Remove(filePath)
			if err != nil {
				logs.Write(logs.Log{
					Message:   fmt.Sprintf("[%v]删除本地文件失败, %s, %s", tasks[i].ID, filePath, err),
					Level:     logs.LogLevelError,
					CreatedAt: time.Now(),
				})
			}
		}

		convertPath := fmt.Sprintf("%s%s", configs.Config.Other.StorePath, tasks[i].ConvertPath)
		if tasks[i].ConvertPath != "" && util.PathExists(convertPath) {
			err = os.Remove(convertPath)
			if err != nil {
				logs.Write(logs.Log{
					Message:   fmt.Sprintf("[%v]删除本地文件失败, %s, %s", tasks[i].ID, convertPath, err),
					Level:     logs.LogLevelError,
					CreatedAt: time.Now(),
				})
			}
		}
	}
	log.Println("clear file processed.")
}
