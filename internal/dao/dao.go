package dao

import (
	"gin-gorm/global"
	"gorm.io/gorm"
)

type Dao struct {
	db *gorm.DB
}

func NewDao() *Dao {
	return &Dao{db: global.DBEngine}
}
