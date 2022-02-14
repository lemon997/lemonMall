package global

import (
	"fmt"
)

func NewMQURL() string {
	return fmt.Sprintf("%s://%s:%s@%s/", RabbitMQSetting.Head, RabbitMQSetting.Username,
		RabbitMQSetting.Password, RabbitMQSetting.Addr)
}
