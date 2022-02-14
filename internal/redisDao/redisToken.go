package redisDao

import (
	"context"
	"time"

	"github.com/lemon997/lemonMall/global"
)

type RedisToken struct{}

func NewMethodsRedisToken() RedisToken {
	return RedisToken{}
}

func (r RedisToken) CreateSetnxKV(ctx context.Context, key, token string, expire time.Duration) (bool, error) {
	//Redis数据库为0,key采取"id:name"得出的字符串为key,
	rdb := global.RedisEngine
	createSuccess, err := rdb.SetNX(ctx, key, token, expire).Result()
	if err != nil {
		global.Logger.Errorf(ctx, "redisDao.CreateSetnxKV报错err= %v", err)
		return false, err
	}
	//field存在则返回false,field不存在则返回true
	return createSuccess, nil
}

func (r RedisToken) QueryGetKV(ctx context.Context, key string) (string, error) {
	rdb := global.RedisEngine
	res, err := rdb.Get(ctx, key).Result()
	if err != nil {
		global.Logger.Errorf(ctx, "redisDao.QueryGetKV报错err= %v", err)
		return "", err
	}
	return res, nil
}

func (r RedisToken) ModifySetKV(ctx context.Context, key, token string, expire time.Duration) error {
	rdb := global.RedisEngine
	//覆盖旧值返回0
	_, err := rdb.SetEX(ctx, key, token, expire).Result()
	if err != nil {
		global.Logger.Errorf(ctx, "redisDao.ModifyHsetSingleKV报错err= %v", err)
		return err
	}
	return nil
}

func (r RedisToken) DelKV(ctx context.Context, key string) error {
	rdb := global.RedisEngine

	_, err := rdb.Del(ctx, key).Result()
	if err != nil {
		global.Logger.Errorf(ctx, "redisDao.DelKV报错err= %v", err)
		return err
	}
	return nil
}
