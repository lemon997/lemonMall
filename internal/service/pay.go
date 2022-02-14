package service

import (
	"fmt"

	"github.com/lemon997/lemonMall/internal/dao"
)

func (svc Service) Payment(customerID int64, orderNos, prices []int64, payType string) error {
	//先查余额，余额够了，再查订单状态，如果为2,则报错返回，如果为1,则修改该订单，然后付款，为了保证原子性，整个操作使用事务以及悲观锁
	// 不能使用Redis分布式锁,因为无法避免丢失锁
	if len(orderNos) != len(prices) {
		return fmt.Errorf("pay fail")
	}
	return dao.NewOrderMethod(svc.ctx).SetPay(customerID, orderNos, prices, payType)

}
