package service

import (
	"encoding/json"
	"strconv"

	"github.com/lemon997/lemonMall/internal/model"

	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/dao"
	"github.com/lemon997/lemonMall/internal/redisDao"
)

func (svc Service) AddElementsInRedisWithGoodDetail(productId int) error {
	//将分类信息写到Redis中
	goodDetail, err := dao.NewGoodsResponseMethodMethod(svc.ctx).Get(int64(productId))
	if err != nil {
		global.Logger.Errorf(svc.ctx, "service.goods.AddE.DB.get, err= %v", err)
		return err
	}
	date, err := json.Marshal(goodDetail)
	if err != nil {
		global.Logger.Errorf(svc.ctx, "service.goods.AddE.json.Mar, err= %v", err)
		return err
	}
	key := svc.GenerateGoodsDatailsKeyUseForRedis()
	field := strconv.Itoa(productId)
	return redisDao.NewMethodsRedisGoodDetail(svc.ctx).HSetNX(key, field, date)
}

func (svc Service) GetGoodDatails(productId int) (model.GoodDetailResponse, error) {
	//获取商品详情数据
	key := svc.GenerateGoodsDatailsKeyUseForRedis()
	field := strconv.Itoa(productId)

	//检查Redis是否有该商品，没有则添加
	rdb := global.RedisEngine
	exist, err := rdb.HExists(svc.ctx, key, field).Result()
	if err != nil {
		global.Logger.Errorf(svc.ctx, "service.goods.GetGoodDatails.Hexists, err= %v", err)
		return model.GoodDetailResponse{}, err
	}
	if !exist {
		err = svc.AddElementsInRedisWithGoodDetail(productId)
		if err != nil {
			return model.GoodDetailResponse{}, err
		}
	}
	//从Redis获取
	data, err := rdb.HGet(svc.ctx, key, field).Bytes()
	if err != nil {
		global.Logger.Errorf(svc.ctx, "service.goods.GetDatails.HGet,err= %v", err)
		return model.GoodDetailResponse{}, err
	}
	goodDetails := model.GoodDetailResponse{}
	json.Unmarshal(data, &goodDetails)
	return goodDetails, nil
}

func (svc Service) GenerateGoodsDatailsKeyUseForRedis() string {
	// 生成redis的key
	return "goodsDetail"
}
