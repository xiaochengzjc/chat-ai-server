// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package routers

import (
	"gin-gorm/internal/controller"
	"gin-gorm/internal/dao"
	"gin-gorm/internal/service"
)

// Injectors from wire.go:

func GetUserCtr() *controller.UserCtr {
	daoDao := dao.NewDao()
	userDao := dao.NewUserDao(daoDao)
	userSvc := service.NewUserSvc(userDao)
	userCtr := controller.NewUserCtr(userSvc)
	return userCtr
}

func GetWechatCtr() *controller.WechatCtr {
	daoDao := dao.NewDao()
	authDao := dao.NewAuthDao(daoDao)
	authSvc := service.NewAuthSvc(authDao)
	userDao := dao.NewUserDao(daoDao)
	wechatSvc := service.NewWechatSvc(authSvc, userDao)
	wechatCtr := controller.NewWechatCtr(wechatSvc)
	return wechatCtr
}

func GetWSCtr() *controller.WSCtr {
	daoDao := dao.NewDao()
	userDao := dao.NewUserDao(daoDao)
	userSvc := service.NewUserSvc(userDao)
	chatSvc := service.NewChatSvc()
	wsCtr := controller.NewWSCtr(userSvc, chatSvc)
	return wsCtr
}