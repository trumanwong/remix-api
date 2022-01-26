package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"remix-api/configs"
	"remix-api/internal/logs"
	"strconv"
	"time"
)

var db *gorm.DB
var err error

type Model struct {
	ID        uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func Setup() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", configs.Config.MySQL.User, configs.Config.MySQL.Password, configs.Config.MySQL.Host+":"+configs.Config.MySQL.Port, configs.Config.MySQL.Database)
	newLogger := logs.NewLogger(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			Colorful:                  false,       // Disable color
			IgnoreRecordNotFoundError: false,
		},
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   configs.Config.MySQL.TablePrefix, // table name prefix, table for `User` would be `t_users`
			SingularTable: false,                            // use singular table name, table for `User` would be `users` with this option enabled
		},
		Logger: newLogger,
	})
	db = db.Debug()

	if err != nil {
		// logs.New().Error("数据库连接失败", err)
		return
	}
}

// Paginate 分页查询
func Paginate(params map[string]interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if _, ok := params["export"]; ok && params["export"] == "1" {
			return db
		}
		page := 1
		if _, ok := params["page"]; ok {
			page, _ = strconv.Atoi(params["page"].(string))
		}
		if page <= 0 {
			page = 1
		}

		perPage := 10
		if _, ok := params["per_page"]; ok {
			perPage, _ = strconv.Atoi(params["per_page"].(string))
		}
		switch {
		case perPage > 100:
			perPage = 100
		case perPage <= 0:
			perPage = 10
		}

		offset := (page - 1) * perPage
		return db.Offset(offset).Limit(perPage)
	}
}

func GetDb() *gorm.DB {
	return db
}
