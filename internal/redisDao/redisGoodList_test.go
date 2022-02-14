package redisDao_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/go-redis/redis/v8"
)

type SwipeImg struct {
	ImgUrl string `json:"img_url" db:"img_url"`
	ImgId  int    `json:"img_id" db:"img_id"`
}

func Add(args ...[]byte) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.172.10:6379",
		Password: "123456", // no password set
		DB:       0,        // use default DB
	})
	fmt.Println(args)
	var z []*redis.Z
	var a float64 = 1
	for i := range args {
		tmp := redis.Z{Score: a, Member: args[i]}
		z = append(z, &tmp)
		a++
	}
	key := "priceSort"
	ctx := context.Background()
	res := rdb.ZAddNX(ctx, key, z...)
	if res.Err() != nil {
		fmt.Println(res.Err())
	}
}

func TestRedisGoodList_ZAdd(t *testing.T) {
	tmp := SwipeImg{ImgUrl: "123", ImgId: 200}
	slipe1, _ := json.Marshal(tmp)
	tmp2 := SwipeImg{ImgUrl: "199", ImgId: 300}
	slipe12, _ := json.Marshal(tmp2)
	Add(slipe1, slipe12)
}
