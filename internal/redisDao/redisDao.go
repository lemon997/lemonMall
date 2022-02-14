package redisDao

import (
	"context"

	"github.com/lemon997/lemonMall/global"
)

type RedisDao struct {
	ctx context.Context
}

func NewRedisDao(c context.Context) RedisDao {
	return RedisDao{ctx: c}
}

func (r RedisDao) DelAll() error {
	return global.RedisEngine.FlushDB(r.ctx).Err()
}
