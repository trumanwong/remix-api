package models

import (
	"gorm.io/gorm"
)

type User struct {
	Model
	Name               string         `gorm:"column:name" json:"name"`
	Nickname           string         `gorm:"column:nickname" json:"nickname"`
	Avatar             string         `gorm:"column:avatar" json:"avatar"`
	Email              string         `gorm:"column:email" json:"email"`
	ConfirmCode        string         `gorm:"column:confirm_code" json:"confirm_code"`
	Status             int8           `gorm:"column:status" json:"status"`
	IsAdmin            int8           `gorm:"column:is_admin" json:"is_admin"`
	Password           string         `gorm:"column:password" json:"password"`
	GithubID           string         `gorm:"column:github_id" json:"github_id"`
	GithubName         string         `gorm:"column:github_name" json:"github_name"`
	GithubUrl          string         `gorm:"column:github_url" json:"github_url"`
	WeiboName          string         `gorm:"column:weibo_name" json:"weibo_name"`
	WeiboLink          string         `gorm:"column:weibo_link" json:"weibo_link"`
	Website            string         `gorm:"column:website" json:"website"`
	Description        string         `gorm:"column:description" json:"description"`
	EmailNotifyEnabled uint8          `gorm:"column:email_notify_enable" json:"email_notify_enabled"`
	RememberToken      string         `gorm:"column:remember_token" json:"remember_token"`
	DeletedAt          gorm.DeletedAt `json:"deleted_at"`
}
