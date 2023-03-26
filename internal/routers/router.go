package routers

import (
	"gin-gorm/internal/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Translations())
	r.Use(middleware.Cors())

	wechatCtr := GetWechatCtr()
	wsCtr := GetWSCtr()
	userCtr := GetUserCtr()

	api := r.Group("/api")
	{
		//需要token校验
		auth := api.Group("")
		auth.Use(middleware.JWT())
		{

		}

		//不需要token 校验
		api.POST("/wechat/codeLogin", wechatCtr.Login)
		api.POST("/user/saveUser", userCtr.SaveUser)
	}

	api.Any("/ws", wsCtr.HandShake)

	return r
}
