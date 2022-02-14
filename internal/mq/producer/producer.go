package producer

import (
	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/mq/delayqueue"
	"github.com/lemon997/lemonMall/internal/mq/messagequeue"
)

func DBDelayOrderSend(body []byte) error {
	//DB订单过期延时队列
	queueName := "delay_order" //延时队列
	expiration := "1800000"    //30分钟
	// expiration := "5000"
	uri := global.NewMQURL()
	err := delayqueue.PublishDelary(uri, queueName, expiration, body, true)
	return err
}

func DBReduceInventorySend(body []byte) error {
	uri := global.NewMQURL()
	exchange := "mysql_inventory"
	exchangeType := "topic"
	routerKey := "amq.mysql.reduceInventory"
	err := messagequeue.Publish(uri, exchange, exchangeType, routerKey, body, true)
	return err
}
