package service

import (
	"errors"
	"fmt"
)

const (
	ErrRedisLoad uint32 = iota
	ErrDBLoad
	ErrRedisGet
)

var (
	ErrLoadRedis = errors.New("写入redis失败")
	ErrLoadDB    = errors.New("写入DB失败")
	ErrGetRedis  = errors.New("查询Redis失败")
)

type ErrService struct {
	//Inner接收未定义错误，或者根据ErrNum自定义错误
	ErrNum uint32
	Msg    string
	Inner  error
}

func (e ErrService) Error() string {
	return fmt.Sprintf("msg=%v, err=%v", e.Msg, e.Inner)
}

func NewErrService(n uint32, msg string, e error) error {
	if e != nil {
		return e
	}
	a := ErrService{}
	switch n {
	case ErrRedisLoad:
		a.Inner = ErrLoadRedis
	case ErrDBLoad:
		a.Inner = ErrLoadDB
	case ErrRedisGet:
		a.Inner = ErrGetRedis
	}
	a.ErrNum = n
	a.Msg = msg
	return a
}
