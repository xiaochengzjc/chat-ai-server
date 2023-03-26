package controller

import (
	"gin-gorm/internal/service"
	"gin-gorm/pkg/app"
	"gin-gorm/pkg/errcode"
	"github.com/gin-gonic/gin"
)

type WechatCtr struct {
	wechatSvc *service.WechatSvc
}

func NewWechatCtr(wechatSvc *service.WechatSvc) *WechatCtr {
	return &WechatCtr{
		wechatSvc: wechatSvc,
	}
}

func (c *WechatCtr) Login(ctx *gin.Context) {
	wechatReq := service.WechatLoginReq{}
	response := app.NewResponse(ctx)
	valid, errs := app.BindAndValid(ctx, &wechatReq)
	if !valid {
		errRsp := errcode.InvalidParams.WithDetails(errs.Errors()...)
		response.ToErrResponse(errRsp)
		return
	}
	rsp, err := c.wechatSvc.Login(ctx, &wechatReq)
	if err != nil {
		response.ToErrResponse(errcode.ErrorWechatLogin)
		return
	}
	response.ToResponse(rsp)
}
