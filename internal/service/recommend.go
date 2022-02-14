package service

import (
	"github.com/lemon997/lemonMall/internal/dao"
	"github.com/lemon997/lemonMall/internal/model"
	"github.com/lemon997/lemonMall/internal/redisDao"
)

func (svc Service) GetRecommendUrlInDB(id int64) ([]model.RecommendImg, error) {
	//从DB获取数据

	return dao.NewRecommendMethod(svc.ctx).Select(id)
}

func (svc Service) InsertRecommendInDB(url string) error {
	//单次插入DB
	return dao.NewRecommendMethod(svc.ctx).Insert(url)
}

func (svc Service) DelRecommendInDB(id int64) error {
	return dao.NewRecommendMethod(svc.ctx).Del(id)
}

func (svc Service) GetRecommendUrlinRedis(key string) ([]string, error) {
	return redisDao.NewMethodsRedisRecommend(svc.ctx).Hvals(key)
}

func (svc Service) InsertRecommendUrlInRedis(key string, value model.RecommendImg) error {
	return redisDao.NewMethodsRedisRecommend(svc.ctx).Hset(key, value)
}

func (svc Service) ModifyMultipleRecommendInRedis(key string, value []model.RecommendImg) error {
	// 描述：修改多个域
	// value中没有的，redis存在，则不会进行删除，
	// value中有的，redis存在则不会覆盖，
	// value中有的，redis没有则创建。
	return redisDao.NewMethodsRedisRecommend(svc.ctx).Hmset(key, value)
}

func (svc Service) ModifyRecommendInRedis(key string, value model.RecommendImg) error {
	// 修改单个域
	return redisDao.NewMethodsRedisRecommend(svc.ctx).Hset(key, value)
}

func (svc Service) DelRecommendInRedis(key, filed string) (bool, error) {
	//删除单个域

	return redisDao.NewMethodsRedisRecommend(svc.ctx).Hdel(key, filed)
}
