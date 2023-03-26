package setting

import "time"

type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type AppSettingS struct {
	DefaultPageSize int
	MaxPageSize     int
	LogSavePath     string
	LogFileName     string
	LogFileExt      string
}

type DatabaseSettingS struct {
	DBType       string
	UserName     string
	Password     string
	Host         string
	DBName       string
	TablePrefix  string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
}

type MiniProgramSettingS struct {
	AppID         string `required:"true" env:"AppID"`
	Secret        string `required:"true" env:"Secret"`
	RedisAddr     string `env:"RedisAddr"`
	MessageToken  string `env:"MessageToken"`
	MessageAesKey string `env:"MessageAesKey"`
}

type JWTSettingS struct {
	Secret string
	Issuer string
	Expire time.Duration
}

type OpenAISettingS struct {
	BaseURL     string
	AuthToken   string
	EnableProxy bool
	ProxyUrl    string
	MaxTokens   int
}

type ChatSettingS struct {
	MaxReq int64
}

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}
