package setting

import (
	"time"
)

//配置文件处理

type ServerSettingS struct {
	RunMode      string
	HttpPort     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type AppSettingS struct {
	DefaultPageSize      int
	MaxPageSize          int
	LogSavePath          string
	LogFileName          string
	LogFileExt           string
	UploadSavePath       string
	UploadServerUrl      string
	UploadImageMaxSize   int
	UploadImageAllowExts []string
}

type DatabaseSettingS struct {
	DBType       string
	UserName     string
	Password     string
	Host         string
	DBName       string
	DBName2      string
	Charset      string
	ParseTime    bool
	MaxIdleConns int
	MaxOpenConns int
}

// TablePrefix  string

type JWTSettingS struct {
	Secret string
	Issuer string
	Expire time.Duration
}

type RedisSettingS struct {
	Addr     string
	DB       int
	Password string
}

type RabbitMQSettingS struct {
	Addr     string
	Username string
	Password string
	Head     string
}

func (s *Setting) ReadSection(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}

	return nil
}
