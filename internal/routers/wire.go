//go:build wireinject
// +build wireinject

package routers

import (
	"gin-gorm/internal/controller"
	"gin-gorm/internal/dao"
	"gin-gorm/internal/service"
	"github.com/google/wire"
)

func GetUserCtr() *controller.UserCtr {
	wire.Build(controller.NewUserCtr, service.NewUserSvc, dao.NewUserDao, dao.NewDao)
	return &controller.UserCtr{}
}

func GetWechatCtr() *controller.WechatCtr {
	wire.Build(controller.NewWechatCtr, service.NewWechatSvc, service.NewAuthSvc, dao.NewUserDao, dao.NewAuthDao, dao.NewDao)
	return &controller.WechatCtr{}
}

func GetWSCtr() *controller.WSCtr {
	wire.Build(controller.NewWSCtr, service.NewChatSvc, service.NewUserSvc, dao.NewUserDao, dao.NewDao)
	return &controller.WSCtr{}
}
