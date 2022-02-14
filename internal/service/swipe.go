package service

import (
	"github.com/lemon997/lemonMall/internal/dao"
	"github.com/lemon997/lemonMall/internal/model"
	"github.com/lemon997/lemonMall/internal/redisDao"
)

func (svc Service) GetSwipeUrlInDB(id int64) ([]model.SwipeImg, error) {
	//从DB获取数据
	return dao.NewSwipeMethod(svc.ctx).Select(id)
}

func (svc Service) InsertSwipeInDB(url string) error {
	//单次插入DB
	return dao.NewSwipeMethod(svc.ctx).Insert(url)
}

func (svc Service) DelSwipeInDB(id int64) error {
	return dao.NewSwipeMethod(svc.ctx).Del(id)
}

func (svc Service) GetSwipeUrlinRedis(key string) ([]string, error) {
	return redisDao.NewMethodsRedisSwipe().Hvals(svc.ctx, key)
}

func (svc Service) InsertSwipeUrlInRedis(key string, value model.SwipeImg) error {
	return redisDao.NewMethodsRedisSwipe().Hset(svc.ctx, key, value)
}

func (svc Service) ModifyMultipleSwipeInRedis(key string, value []model.SwipeImg) error {
	// 描述：修改多个域
	// value中没有的，redis存在，则不会进行删除，
	// value中有的，redis存在则不会覆盖，
	// value中有的，redis没有则创建。
	return redisDao.NewMethodsRedisSwipe().HmsetKV(svc.ctx, key, value)
}

func (svc Service) ModifySwipeInRedis(key string, value model.SwipeImg) error {
	// 修改单个域
	return redisDao.NewMethodsRedisSwipe().Hset(svc.ctx, key, value)
}

func (svc Service) DelSwipeInRedis(key, filed string) (bool, error) {
	//删除单个域
	return redisDao.NewMethodsRedisSwipe().Hdel(svc.ctx, key, filed)
}
