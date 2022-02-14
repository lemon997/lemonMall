package redisDao

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/lemon997/lemonMall/global"
)

//分类的右边商品排序
type RedisGoodList struct {
	ctx context.Context
}

func NewMethodsRedisGoodList(c context.Context) RedisGoodList {
	return RedisGoodList{ctx: c}
}

func (r RedisGoodList) ZAdd(key string, score []float64, date ...[]byte) error {
	//data应该为结构体序列化后的数据
	rdb := global.RedisEngine
	var z []*redis.Z
	for i := range date {
		tmp := redis.Z{Score: score[i], Member: date[i]}
		z = append(z, &tmp)
	}

	_, err := rdb.TxPipelined(r.ctx, func(pipe redis.Pipeliner) error {
		err := pipe.ZAddNX(r.ctx, key, z...).Err()
		return err
	})
	return err
}

func (r RedisGoodList) ZRange(key string, startIndex, endIndex int64) ([]string, error) {
	//根据score从小到大排序
	rdb := global.RedisEngine
	res, err := rdb.ZRange(r.ctx, key, startIndex, endIndex).Result()
	return res, err
}

func (r RedisGoodList) ZRevRange(key string, startIndex, endIndex int64) ([]string, error) {
	//根据score从大到小排序
	rdb := global.RedisEngine
	res, err := rdb.ZRevRange(r.ctx, key, startIndex, endIndex).Result()
	return res, err
}

func (r RedisGoodList) CheckKeyExists(key string) int64 {
	//如果num是0则不存在
	rdb := global.RedisEngine
	num, _ := rdb.Exists(r.ctx, key).Result()
	return num
}
