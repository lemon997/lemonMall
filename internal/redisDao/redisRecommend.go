package redisDao

import (
	"context"

	"strconv"

	"github.com/go-redis/redis/v8"

	"github.com/lemon997/lemonMall/internal/model"

	"github.com/lemon997/lemonMall/global"
)

type RedisRecommend struct {
	ctx context.Context
}

func NewMethodsRedisRecommend(c context.Context) RedisRecommend {
	return RedisRecommend{ctx: c}
}

func (r RedisRecommend) Hmset(key string, value []model.RecommendImg) error {

	//返回值，error，服务器能连上则返回nil
	//入参，r.ctx,key,model.SwipeImg切片
	//描述：根据切片插入到redis中，域为id, 值为url,目标id存在则跳过覆盖操作
	index := 0
	rdb := global.RedisEngine

	txf := func(tx *redis.Tx) error {

		exist, err := tx.HExists(r.ctx, key, strconv.Itoa(int(value[index].ImgId))).Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if exist {
			global.Logger.Infof(r.ctx, "filed exist, insert new value fail, filed is: %d", value[index].ImgId)
			return err
		}

		_, err = tx.TxPipelined(r.ctx, func(pipe redis.Pipeliner) error {
			result := pipe.HMSet(r.ctx, key, value[index].ImgId, value[index].ImgUrl)
			return result.Err()
		})

		return err
	}

	for i := 0; i < len(value); i++ {
		err := rdb.Watch(r.ctx, txf, key)
		index++
		if err == redis.TxFailedErr {
			continue
		} else if err != nil {
			global.Logger.Errorf(r.ctx, "Watch fail, err= %v", err)
			return err
		}
	}
	return nil
}

func (r RedisRecommend) Hkeys(key string) ([]string, error) {
	rdb := global.RedisEngine
	strings, err := rdb.HKeys(r.ctx, key).Result()
	return strings, err
}

func (r RedisRecommend) Hvals(key string) ([]string, error) {
	rdb := global.RedisEngine
	res, err := rdb.HVals(r.ctx, key).Result()
	return res, err
}

func (r RedisRecommend) Hdel(key, filed string) (bool, error) {
	//返回值，删除成功返回true，nil，失败则返回false,err
	//入参数，r.ctx,key，和域
	//删除单个域，先查询，再删除，使用watch和piple，事务

	rdb := global.RedisEngine

	txf := func(tx *redis.Tx) error {
		exist, err := tx.HExists(r.ctx, key, filed).Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if !exist {
			return err
		}

		_, err = tx.TxPipelined(r.ctx, func(pipe redis.Pipeliner) error {
			res := pipe.HDel(r.ctx, key, filed)
			return res.Err()
		})
		return err
	}

	err := rdb.Watch(r.ctx, txf, key)
	if err != redis.TxFailedErr || err != nil {
		return false, err
	}

	return true, nil
}

func (r RedisRecommend) Hset(key string, value model.RecommendImg) error {
	//返回值: 服务器错误返回err
	//入参: r.ctx, key, model.SwipeImg类型的value
	//描述: 插入1个结构体，根据key,结构体的id和url赋值
	rdb := global.RedisEngine
	res := rdb.HSet(r.ctx, key, value.ImgId, value.ImgUrl)
	return res.Err()
}
