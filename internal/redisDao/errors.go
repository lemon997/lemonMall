package redisDao

import (
	"errors"
)

var (
	DestNotExists = errors.New("目标不存在")
	TxFail        = errors.New("事务失败")
)
