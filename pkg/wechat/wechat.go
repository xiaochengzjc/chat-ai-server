package wechat

import (
	"gin-gorm/global"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/miniProgram"
)

func NewMiniProgram() *miniProgram.MiniProgram {
	MiniProgramApp, err := miniProgram.NewMiniProgram(&miniProgram.UserConfig{
		AppID:     global.MiniProgramSetting.AppID,  // 小程序appid
		Secret:    global.MiniProgramSetting.Secret, // 小程序app secret
		HttpDebug: true,
		Log: miniProgram.Log{
			Level: "debug",
			File:  "./wechat.log",
		},
		// 可选，不传默认走程序内存
		//Cache: kernel.NewRedisClient(&kernel.RedisOptions{
		//	Addr:     "127.0.0.1:6379",
		//	Password: "",
		//	DB:       0,
		//}),
	})
	if err != nil {
		global.Logger.Fatalf("NewMiniProgram调用失败，%v", err)
	}

	return MiniProgramApp
}
