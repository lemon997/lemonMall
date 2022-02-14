package queueservice

import (
	"context"
	"encoding/json"

	"github.com/lemon997/lemonMall/global"

	"github.com/lemon997/lemonMall/internal/model"
)

//以下是收到消息的逻辑
func MySQLReduceInventory(body []byte) error {
	reduce := model.ReduceInventory{}
	_ = json.Unmarshal(body, &reduce)
	ctx := context.TODO()
	svc := MQNew(ctx)
	for i := 0; i < len(reduce.ProductId); i++ {
		err := svc.ReduceInventory(reduce.ProductId[i], reduce.Num[i])
		if err != nil {
			return err
		}
	}
	return nil
}

// 其中orderStatus对比：
// 0，全部
// 1，待付款
// 2，待发货
// 3，待收款
// 4，待评价
// 5，已过期
func MySQLAddInventory(body []byte) {
	delay := model.DelayOrder{}
	json.Unmarshal(body, &delay)
	ctx := context.TODO()
	svc := MQNew(ctx)

	for i := 0; i < len(delay.OrderNo); i++ {
		data, err := svc.GetOrderStatus(delay.CustomerId, delay.OrderNo[i])
		if err != nil {
			global.Logger.Errorf(ctx, "service.queueReceiver.MySQLAddInventory.GetOrderStatus, err= %v", err)
			continue
		}
		if data.OrderStatus != uint8(1) {
			continue
		}

		//未付款，且是过期订单，修改该订单信息，设置status为5
		err = svc.SetOrderStatus(delay.CustomerId, delay.OrderNo[i], uint8(5))
		if err != nil {
			global.Logger.Errorf(ctx, "service.queueReceiver.MySQLAddInventory.setOrderStatus, err= %v", err)
		}
		//增加库存
		err = svc.AddInventory(delay.ProductId[i], delay.Num[i])
		if err != nil {
			global.Logger.Errorf(ctx, "service.queueReceiver.MySQLAddInventory.Addinventory, err= %v")
		}

	}
}
