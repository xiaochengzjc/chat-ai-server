package model

import "gorm.io/gorm"

type Auth struct {
	*gorm.Model
	AppKey    string `json:"app_key" gorm:"column:app_key"`
	AppSecret string `json:"app_secret" gorm:"column:app_secret"`
}

func (a Auth) TableName() string {
	return "ai_auth"
}
