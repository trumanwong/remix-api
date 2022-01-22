package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strconv"
)

type Config struct {
	Model
	Content ConfigContent `gorm:"content;type:json" json:"content"`
}

type ConfigContent map[string]interface{}

func (this ConfigContent) Value() (driver.Value, error) {
	b, err := json.Marshal(this)
	return string(b), err
}

func (this *ConfigContent) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), this)
}

func (this Config) scopeType(params map[string]interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if _, ok := params["type"]; ok {
			configType, err := strconv.ParseUint(fmt.Sprintf("%v", params["type"]), 10, 8)
			if err == nil {
				return db.Where("type = ?", configType)
			}
		}
		return db
	}
}

func (this Config) scopeId(params map[string]interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if _, ok := params["id"]; ok {
			id, err := strconv.ParseUint(fmt.Sprintf("%v", params["id"]), 10, 64)
			if err == nil {
				return db.Where("id = ?", id)
			}
		}
		return db
	}
}

func (this Config) First(params map[string]interface{}) (*Config, error) {
	config := new(Config)
	err := db.Scopes(this.scopeId(params), this.scopeType(params)).First(&config).Error
	return config, err
}
