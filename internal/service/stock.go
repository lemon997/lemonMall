package service

import (
	"fmt"

	"github.com/lemon997/lemonMall/internal/dao"
	"github.com/lemon997/lemonMall/internal/model"
	"github.com/lemon997/lemonMall/internal/redisDao"
)

func (svc Service) SetStock(productID int64) error {
	res, err := dao.NewProductsMethod(svc.ctx).GetStock(productID)
	if err != nil {
		return err
	}
	key := model.ProductKey(productID)
	ok, err := redisDao.NewStockMethod(svc.ctx).SetNX(key, res.RemainingAmount)
	if err != nil {
		return NewErrService(0, "", err)
	}

	s := fmt.Sprintf("SetSrock,key=%v", key)
	if !ok {
		return NewErrService(ErrRedisLoad, s, nil)
	}
	return nil
}

func (svc Service) DelStock(productID int64) error {
	key := model.ProductKey(productID)
	return redisDao.NewStockMethod(svc.ctx).Del(key)
}

func (svc Service) SetStockAll() error {
	array := model.CategoryArray
	childrenID := make([]int64, 0, 8)

	for i := 0; i < len(array); i++ {
		res, err := dao.NewCategoryChildrenMethod(svc.ctx).Select(array[i])
		if err != nil {
			return err
		}
		for j := 0; j < len(res); j++ {
			childrenID = append(childrenID, res[j].CategoryChildrenID)
		}
	}

	productIDs := make([]int64, 0, 8)
	nums := make([]int64, 0, 8)

	d := dao.NewProductsMethod(svc.ctx)
	for i := 0; i < len(childrenID); i++ {
		id, err := d.Select(childrenID[i])
		if err != nil {
			return err
		}

		for j := 0; j < len(id); j++ {
			productIDs = append(productIDs, id[j].ProductID)
			n, err := d.GetStock(id[j].ProductID)
			if err != nil {
				return err
			}
			nums = append(nums, n.RemainingAmount)
		}
	}

	for i := 0; i < len(productIDs); i++ {
		key := model.ProductKey(productIDs[i])
		_, err := redisDao.NewStockMethod(svc.ctx).SetNX(key, nums[i])
		if err != nil {
			return err
		}
	}
	return nil
}
