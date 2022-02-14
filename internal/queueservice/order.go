package queueservice

import (
	"github.com/lemon997/lemonMall/internal/dao"
	"github.com/lemon997/lemonMall/internal/model"
)

func (svc MQService) GetOrderStatus(customerID, orderNo int64) (model.Order, error) {
	return dao.NewOrderMethod(svc.ctx).GetStatus(customerID, orderNo)
}

func (svc MQService) SetOrderStatus(customerID, orderNo int64, status uint8) error {
	return dao.NewOrderMethod(svc.ctx).SetStatus(customerID, orderNo, status)
}
