package global

import (
	"gin-gorm/pkg/logger"
	"gin-gorm/pkg/setting"
)

var (
	ServerSetting      *setting.ServerSettingS
	AppSetting         *setting.AppSettingS
	DatabaseSetting    *setting.DatabaseSettingS
	MiniProgramSetting *setting.MiniProgramSettingS

	Logger        *logger.Logger
	JWTSetting    *setting.JWTSettingS
	OpenAISetting *setting.OpenAISettingS
	ChatSetting   *setting.ChatSettingS
)
