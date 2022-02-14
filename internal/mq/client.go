package mq

import (
	"github.com/lemon997/lemonMall/internal/mq/consumer"
)

func DelayQueue() error {
	//程序初始化时就要建立消费者呢
	var err error
	err = consumer.DBDelayOrderReciver()
	return err
}

func MessageQueue() error {
	var err error
	err = consumer.DBReduceInventoryReciver()
	return err
}
