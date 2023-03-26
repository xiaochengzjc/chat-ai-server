package controller

import (
	"gin-gorm/internal/service"
	"gin-gorm/pkg/app"
	"gin-gorm/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type UserCtr struct {
	userSvc *service.UserSvc
}

func NewUserCtr(userSvc *service.UserSvc) *UserCtr {
	return &UserCtr{userSvc: userSvc}
}

func (c *UserCtr) SaveUser(ctx *gin.Context) {
	userRequest := service.UserRequest{}
	valid, errs := app.BindAndValid(ctx, &userRequest)
	if !valid {
		errRsp := errcode.InvalidParams.WithDetails(errs.Errors()...)
		app.NewResponse(ctx).ToErrResponse(errRsp)
		return
	}

	err := c.userSvc.AddOrUpdateUser(&userRequest)
	if err != nil {
		app.NewResponse(ctx).ToErrResponse(errcode.ErrorAddData)
		return
	}
	app.NewResponse(ctx).ToResponse(nil)
	return
}
