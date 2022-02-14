package global

import (
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/lemon997/lemonMall/common/setting"
)

var RedisEngine *redis.Client

var redisOnce sync.Once

func NewRedis(redisSetting *setting.RedisSettingS) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisSetting.Addr,
		Password: redisSetting.Password,
		DB:       redisSetting.DB,
	})
	redisOnce.Do(func() {
		RedisEngine = rdb
	})
	return rdb
}
