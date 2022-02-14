package redisDao

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/lemon997/lemonMall/common/convert"

	"github.com/go-redis/redis/v8"
	"github.com/lemon997/lemonMall/global"
)

type CartProductNum struct {
	ctx context.Context
}

func NewMethodsCartProductNum(c context.Context) CartProductNum {
	return CartProductNum{ctx: c}
}

func (c CartProductNum) AddProductAtCart(key string, productID, num int64) (err error) {
	//isCheck是选择的状态
	//不存在则添加域值
	field1 := strconv.FormatInt(productID, 10)
	// field2 := "isCheck:" + field1
	rdb := global.RedisEngine

	txf := func(tx *redis.Tx) error {
		_, err = tx.TxPipelined(c.ctx, func(pipe redis.Pipeliner) error {
			pipe.HSetNX(c.ctx, key, field1, num)
			// pipe.HSetNX(c.ctx, key, field2, isCheck)
			return nil
		})
		return err
	}

	err = rdb.Watch(c.ctx, txf, key)
	if err == redis.TxFailedErr {
		return nil
	}
	return err

}

func (c CartProductNum) SetProductNumAtCart(key string, productID, num int64) error {
	//修改域的值
	rdb := global.RedisEngine
	field := strconv.FormatInt(productID, 10)

	txf := func(tx *redis.Tx) error {
		exist, err := tx.HExists(c.ctx, key, field).Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if !exist {
			return errors.New(fmt.Sprintf("field not exist,field= %v", field))
		}

		_, err = tx.TxPipelined(c.ctx, func(pipe redis.Pipeliner) error {
			return pipe.HSet(c.ctx, key, field, num).Err()
		})
		return err
	}
	err := rdb.Watch(c.ctx, txf, key)
	if err == redis.TxFailedErr {
		return nil
	}
	return err

}

func (c CartProductNum) DelProductAtCart(key, field string) error {
	//乐观锁，事务实现删除
	rdb := global.RedisEngine

	txf := func(tx *redis.Tx) error {
		exist, err := tx.HExists(c.ctx, key, field).Result()
		if err != nil && err != redis.Nil {
			return err
		}
		if !exist {
			return err
		}

		_, err = tx.TxPipelined(c.ctx, func(pipe redis.Pipeliner) error {
			return pipe.HDel(c.ctx, key, field).Err()
		})
		return err
	}
	err := rdb.Watch(c.ctx, txf, key)
	if err == redis.TxFailedErr {
		return nil
	}
	return err
}

func (c CartProductNum) GetProductListAtCart(key string) (map[string]string, error) {
	rdb := global.RedisEngine
	return rdb.HGetAll(c.ctx, key).Result()
}

func (c CartProductNum) ChangeCheckedStatus(key, field string, value bool) error {
	rdb := global.RedisEngine

	txf := func(tx *redis.Tx) error {
		_, err := tx.TxPipelined(c.ctx, func(pipe redis.Pipeliner) error {
			return pipe.HSet(c.ctx, key, field, value).Err()
		})
		return err
	}
	err := rdb.Watch(c.ctx, txf, key)
	if err == redis.TxFailedErr {
		return nil
	}
	return err
}

func (c CartProductNum) GetProductNumAtCart(key string, productID int64) (int64, error) {
	//获取cart:productID中，field为productID的数量
	// 建议给购物车上锁

	rdb := global.RedisEngine
	fields := strconv.FormatInt(productID, 10)

	res, err := rdb.HGet(c.ctx, key, fields).Result()
	if err == redis.Nil {
		return -1, errors.New("get product num is nil")
	} else if err != nil {
		return -1, err
	}
	num, _ := convert.StrTo(res).MustInt64()
	return num, nil
}

func (c CartProductNum) GetProductNumsAtCart(key string, productID []int64) ([]int64, error) {
	//获取cart:productID中，field为productID的数量
	// 建议给购物车上锁

	rdb := global.RedisEngine
	fields := make([]string, len(productID))
	for i := 0; i < len(productID); i++ {
		fields[i] = strconv.FormatInt(productID[i], 10)
	}

	res, err := rdb.HMGet(c.ctx, key, fields...).Result()
	if err == redis.Nil {
		return nil, errors.New("get product num is nil")
	} else if err != nil {
		return nil, err
	}

	nums := make([]int64, len(res))
	for i := 0; i < len(res); i++ {
		nums[i], _ = convert.StrTo(res[i].(string)).MustInt64()
	}
	return nums, nil
}

func (c CartProductNum) ExistFieldAtCart(key string, field string) (bool, error) {
	rdb := global.RedisEngine
	return rdb.HExists(c.ctx, key, field).Result()
}

func (c CartProductNum) ExistFieldsAtCart(key string, fields []string) (bool, error) {
	rdb := global.RedisEngine
	var exist bool
	var err error
	for i := 0; i < len(fields); i++ {
		exist, err = rdb.HExists(c.ctx, key, fields[i]).Result()
		if err != nil {
			return false, err
		}
		if !exist {
			return false, nil
		}
	}
	return true, nil

}
