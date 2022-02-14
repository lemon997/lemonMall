package global

import (
	"github.com/lemon997/lemonMall/common/logger"
	"github.com/lemon997/lemonMall/common/setting"
)

var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
	Logger          *logger.Logger
	JWTSetting      *setting.JWTSettingS
	RedisSetting    *setting.RedisSettingS
	RabbitMQSetting *setting.RabbitMQSettingS
)

//雪花算法节点
var Node1 *Node
var Node2 *Node
