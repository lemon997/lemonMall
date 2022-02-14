package redisDao

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/lemon997/lemonMall/global"
)

//下单操作，前端传来地址,商品id,商品数量信息，使用hset存储，同时生成分布式ID,key为order:customer_id:product_id,
// 值为字符串地址，商品数量,商品ID,订单编号

var reduceInventoryScript string = `
local array1 = KEYS
local array2 = ARGV
for k, v in ipairs(array1) do
	local value = redis.call("GET", v)
	if ( value - array2[k] < 0 ) 
	then 
		return -1
	end
end
for k, v in ipairs(array1) do
	redis.call("decrby", v, array2[k])
end
return 0
`

func generalNewDelayedOrdelKey(customerID, productID int64) string {
	return "order:" + strconv.FormatInt(customerID, 10) + ":" + strconv.FormatInt(productID, 10)
}

type PlaceOrder struct {
	ctx context.Context
}

func NewPlaceOrderMethod(c context.Context) PlaceOrder {
	return PlaceOrder{ctx: c}
}

func (p PlaceOrder) ReduceInventory(key []string, num []int64) error {
	//key的下标与num下标一一对应
	//减库存操作，使用LUA脚本
	rdb := global.RedisEngine
	if num == nil || len(num) == 0 {
		return errors.New("num must be greater than 0")
	}
	nums := make([]interface{}, len(num))
	for i := 0; i < len(num); i++ {
		nums[i] = num[i]
	}

	n, err := redis.NewScript(reduceInventoryScript).Run(p.ctx, rdb, key, nums).Result()
	if err != nil {
		return err
	}
	if n.(int64) < 0 {
		return errors.New("Not enough stock")
	}
	return nil
}
