package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/trumanwong/go-internal/util"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Model
	UserID      uint64         `gorm:"column:user_id" json:"user_id"`
	UUID        string         `gorm:"column:uuid" json:"uuid"`
	Serial      string         `gorm:"column:serial" json:"serial"`
	Status      uint8          `gorm:"column:status" json:"status"`
	FileName    string         `gorm:"column:file_name" json:"file_name"`
	Size        uint64         `gorm:"column:size" json:"size"`
	Path        string         `gorm:"column:path" json:"path"`
	ConvertPath string         `gorm:"column:convert_path" json:"convert_path"`
	IP          uint32         `gorm:"column:ip" json:"ip"`
	Type        uint8          `gorm:"column:type" json:"type"`
	DeletedAt   gorm.DeletedAt `json:"-"`
	Other       TaskOther      `gorm:"column:other;type:json" json:"-"`
}

type TaskOther map[string]interface{}

func (this TaskOther) Value() (driver.Value, error) {
	b, err := json.Marshal(this)
	return string(b), err
}

func (this *TaskOther) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), this)
}

const (
	// TaskStatusProcessing 处理中
	TaskStatusProcessing uint8 = iota
	// TaskStatusSuccess 处理成功
	TaskStatusSuccess
	// TaskStatusFail 处理失败
	TaskStatusFail
)

const (
	// TaskTypeGifRevert gif倒放
	TaskTypeGifRevert = 0
	// TaskTypeNokiaSms 诺基亚短信
	TaskTypeNokiaSms = 1
	// TaskTypeRemix 鬼畜
	TaskTypeRemix = 2
)

func (this Task) scopeID(params map[string]interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if _, ok := params["id"]; ok {
			id, err := strconv.ParseInt(fmt.Sprintf("%v", params["id"]), 10, 64)
			if err == nil && id > 0 {
				return db.Where("id = ?", id)
			}
		}
		return db
	}
}

func (this Task) scopeUUID(params map[string]interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if _, ok := params["uuid"]; ok {
			uuid := strings.Trim(fmt.Sprintf("%v", params["uuid"]), " ")
			if len(uuid) > 0 {
				return db.Where("uuid = ?", uuid)
			}
		}
		return db
	}
}

func (this Task) scopeSerial(params map[string]interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if _, ok := params["serial"]; ok {
			serial := strings.Trim(fmt.Sprintf("%v", params["serial"]), " ")
			if len(serial) > 0 {
				return db.Where("serial = ?", serial)
			}
		}
		return db
	}
}

func (this Task) scopeUserID(params map[string]interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if _, ok := params["user_id"]; ok {
			userId, err := strconv.ParseInt(fmt.Sprintf("%v", params["user_id"]), 10, 64)
			if err == nil && userId > 0 {
				return db.Where("user_id = ?", userId)
			}
		}
		return db
	}
}

func (this Task) scopeCreatedAt(params map[string]interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if _, ok := params["start_at"]; ok {
			startAt, err := time.ParseInLocation("2006-01-02", fmt.Sprintf("%v", params["start_at"]), time.Local)
			if err == nil {
				db = db.Where("created_at >= ?", startAt)
			}
		}
		if _, ok := params["end_at"]; ok {
			endAt, err := time.ParseInLocation("2006-01-02", fmt.Sprintf("%v", params["end_at"]), time.Local)
			if err == nil {
				db = db.Where("created_at <= ?", endAt)
			}
		}
		return db
	}
}

// Create 新增
func (this *Task) Create() error {
	return db.Create(&this).Error
}

// NewTask 查找task
func NewTask(params map[string]interface{}) (*Task, error) {
	task := new(Task)
	err := db.Model(Task{}).Scopes(
		task.scopeID(params),
		task.scopeUserID(params),
		task.scopeSerial(params),
	).First(&task).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}

// Update 更新
func (this *Task) Update(values interface{}) error {
	return db.Model(Task{}).Where("id = ?", this.ID).Updates(values).Error
}

func (this Task) scopeType(params map[string]interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if _, ok := params["type"]; ok {
			value, err := strconv.ParseUint(fmt.Sprintf("%v", params["type"]), 10, 8)
			if err == nil {
				return db.Where("type = ?", value)
			}
		}
		return db
	}
}

// GetList 获取列表
func (this *Task) GetList(params map[string]interface{}) map[string]interface{} {
	list := make([]Task, 0)
	var total int64
	db := db.Model(Task{}).Scopes(
		this.scopeUUID(params),
		this.scopeSerial(params),
		this.scopeType(params),
		this.scopeCreatedAt(params),
	)
	if params["export"] != "1" {
		db = db.Scopes(util.ScopePaginate(params))
	}
	db.Count(&total)
	if total > 0 {
		db.Order("id desc").Find(&list)
	}
	return map[string]interface{}{
		"list":  list,
		"total": total,
	}
}

// Delete 删除
func (this Task) Delete() error {
	return db.Delete(&this).Error
}
