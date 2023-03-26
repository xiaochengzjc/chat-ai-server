package model

import (
	"gorm.io/gorm"
)

type User struct {
	*gorm.Model
	Name     string `gorm:"column:name" json:"name"`         // 用户名
	Openid   string `gorm:"column:openid" json:"openid"`     // openid
	UnionId  string `gorm:"column:unionid" json:"unionid"`   // unionid
	Password string `gorm:"column:password" json:"password"` // 密码
	Phone    string `gorm:"column:phone" json:"phone"`       // 手机号
	Avatar   string `gorm:"column:avatar" json:"avatar"`     // 头像
	Gender   uint8  `gorm:"column:gender" json:"gender"`     //性别
}

func (u User) TableName() string {
	return "ai_user"
}
