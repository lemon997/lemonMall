package redisDao

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/lemon997/lemonMall/global"
)

const (
	expiredOrderKey = "order:expired"
	min             = "0"
	timeOut         = time.Second * 60
)

type DelayQueue struct {
	ctx context.Context
}

func NewDelayQueueMethod(c context.Context) DelayQueue {
	return DelayQueue{ctx: c}
}

//定时扫描zset列表，扫描过期订单，key为expired order
func (d DelayQueue) DelayTask() {
	var times int64
	var s string
	rdb := global.RedisEngine

	t := time.NewTicker(timeOut)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			//过期订单
			times = time.Now().Unix()
			s = strconv.FormatInt(times, 10)

			z := redis.ZRangeBy{Max: s, Min: min}
			_, err := rdb.ZRevRangeByScore(d.ctx, expiredOrderKey, &z).Result()

			//获取到字符串切片，获取对应的值并进行回滚数量
			if err != nil {
			}

			//删除到时间戳之前的过期订单
			err = d.ZRemByScore(s)
			if err != nil {

			}
		}
	}
}

func (d DelayQueue) ZRemByScore(max string) error {
	rdb := global.RedisEngine
	err := rdb.ZRemRangeByScore(d.ctx, expiredOrderKey, "0", max).Err()
	return err
}

func (d DelayQueue) A() {

}
