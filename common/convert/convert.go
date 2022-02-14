package convert

import (
	"strconv"
)

//类型转换
type StrTo string

func (s StrTo) String() string {
	return string(s)
}

func (s StrTo) Int() (int, error) {
	v, err := strconv.Atoi(s.String())
	return v, err
}

func (s StrTo) MustInt() int {
	v, _ := s.Int()
	return v
}

func (s StrTo) UInt32() (uint32, error) {
	v, err := strconv.Atoi(s.String())
	return uint32(v), err
}

func (s StrTo) MustUInt32() uint32 {
	v, _ := s.UInt32()
	return v
}

func (s StrTo) MustInt64() (int64, error) {
	v, err := strconv.Atoi(s.String())
	return int64(v), err
}

type IntTo int64

func (i IntTo) Int64() int64 {
	return int64(i)
}

func (i IntTo) String() string {
	return strconv.FormatInt(i.Int64(), 10)
}
