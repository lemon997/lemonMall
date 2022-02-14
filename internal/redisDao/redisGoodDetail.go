package redisDao

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/lemon997/lemonMall/global"
)

type RedisGoodDetail struct {
	ctx context.Context
}

func NewMethodsRedisGoodDetail(c context.Context) RedisGoodDetail {
	return RedisGoodDetail{ctx: c}
}

func (r RedisGoodDetail) HSetNX(key, field string, date []byte) error {
	txf := func(tx *redis.Tx) error {
		tx.TxPipelined(r.ctx, func(pipe redis.Pipeliner) error {
			pipe.HSetNX(r.ctx, key, field, date)
			return nil
		})
		return nil
	}
	rdb := global.RedisEngine
	err := rdb.Watch(r.ctx, txf, key)
	if err == redis.TxFailedErr {
		return nil
	}
	return err
}

func (r RedisGoodDetail) HGet(key, field string) ([]byte, error) {
	rdb := global.RedisEngine
	return rdb.HGet(r.ctx, key, field).Bytes()
}
