package consumer

import (
	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/mq/delayqueue"
	"github.com/lemon997/lemonMall/internal/mq/messagequeue"
)

func DBDelayOrderReciver() error {
	uri := global.NewMQURL()
	exchange := "mysql_order"
	exchangeType := "fanout"
	queueName := "delayafter_order"
	delaryName := "delay_order"

	_, err := delayqueue.NewConsumerDelay(uri, exchange, exchangeType, delaryName, queueName)
	if err != nil {
		return err
	}
	return nil
}

func DBReduceInventoryReciver() error {
	uri := global.NewMQURL()
	exchange := "mysql_inventory"
	exchangeType := "topic"
	queueName := "reduce_inventory"
	routerKey := "amq.mysql.reduceInventory"
	_, err := messagequeue.NewConsumer(uri, exchange, exchangeType, queueName, routerKey)
	if err != nil {
		return err
	}
	return nil
}
