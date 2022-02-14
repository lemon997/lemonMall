package redisDao

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/lemon997/lemonMall/global"
)

type Stock struct {
	ctx context.Context
}

func NewStockMethod(c context.Context) Stock {
	return Stock{ctx: c}
}

func (s Stock) SetNX(key string, value int64) (bool, error) {
	rdb := global.RedisEngine
	return rdb.SetNX(s.ctx, key, value, 0).Result()
}

func (s Stock) Del(key string) error {
	rdb := global.RedisEngine

	txf := func(tx *redis.Tx) error {
		n, err := tx.Exists(s.ctx, key).Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if n != 1 {
			return DestNotExists
		}

		_, err = tx.TxPipelined(s.ctx, func(pipe redis.Pipeliner) error {
			res, err := pipe.Del(s.ctx, key).Result()
			if res == 0 {
				return DestNotExists
			}
			return err
		})
		return err
	}
	err := rdb.Watch(s.ctx, txf, key)
	if err == redis.TxFailedErr {
		return nil
	}
	return err
}
