package util

import (
	"github.com/shopspring/decimal"
)

func CentsToYuan(price string) string {
	cents, _ := decimal.NewFromString(price)
	nums, _ := decimal.NewFromString("100")
	return cents.Mul(nums).StringFixedBank(2)
}

func YuanToCents(price string) string {
	decimal.DivisionPrecision = 2
	yuan, _ := decimal.NewFromString(price)
	nums, _ := decimal.NewFromString("100")
	return yuan.Div(nums).StringFixedBank(2)
}
