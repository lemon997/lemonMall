package queueservice

import (
	"github.com/lemon997/lemonMall/internal/dao"
)

func (svc MQService) ReduceInventory(productID, num int64) error {
	return dao.NewProductsMethod(svc.ctx).ReduceInventory(productID, num)
}

func (svc MQService) AddInventory(productID, num int64) error {
	return dao.NewProductsMethod(svc.ctx).AddInventory(productID, num)
}
