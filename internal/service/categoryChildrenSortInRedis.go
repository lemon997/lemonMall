package service

import (
	"encoding/json"
	"strconv"

	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/dao"
	"github.com/lemon997/lemonMall/internal/redisDao"
)

func (svc Service) AddElementsInRedisAndSortByKey(key string, categoryChildrenId int64) error {
	//将分类信息写到Redis中
	res, err := dao.NewProductsMethod(svc.ctx).Select(categoryChildrenId)
	if err != nil {
		global.Logger.Errorf(svc.ctx, "service.categoryChildrenSortInRedis.DB.Select, err= %v", err)
		return err
	}
	goodList := redisDao.NewMethodsRedisGoodList(svc.ctx)
	score := make([]float64, len(res))
	date := make([][]byte, len(res))
	for i := 0; i < len(res); i++ {
		sc, _ := strconv.ParseFloat(res[i].Price, 64)
		score[i] = sc
		b, _ := json.Marshal(res[i])
		date[i] = b
	}
	return goodList.ZAdd(key, score, date...)
}

func (svc Service) GetDescByKeyWithProducts(key string, startIndex, endIndex int64) ([]string, error) {
	// 根据key降序，然后将结果放回
	return redisDao.NewMethodsRedisGoodList(svc.ctx).ZRevRange(key, startIndex, endIndex)
}

func (svc Service) GetAscByKeyWithProducts(key string, startIndex, endIndex int64) ([]string, error) {
	// 根据key升序，然后将结果放回
	return redisDao.NewMethodsRedisGoodList(svc.ctx).ZRange(key, startIndex, endIndex)
}

func (svc Service) SortByKey(key, sortType string, startIndex, endIndex int64) ([]string, error) {
	// 根据sortType选择升序还是降序
	switch sortType {
	case "price":
		//根据价格从小到大排序
		return svc.GetAscByKeyWithProducts(key, startIndex, endIndex)
	case "sales":
		//根据销售量从大到小排序
		return svc.GetDescByKeyWithProducts(key, startIndex, endIndex)
	default:
		return svc.GetDescByKeyWithProducts(key, startIndex, endIndex)
	}

}

func (svc Service) GenerateGoodsDateKeyUseForRedis(key string, categoryChildrenID int64) string {
	// 生成redis的key
	return "goods:" + strconv.Itoa(int(categoryChildrenID)) + ":" + key
}

func (svc Service) CheckCategorySortKeyExists(key string) int64 {
	// 检查分类排序key是否存在
	return redisDao.NewMethodsRedisGoodList(svc.ctx).CheckKeyExists(key)
}
