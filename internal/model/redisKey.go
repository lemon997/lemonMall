package model

import (
	"strconv"
)

//数量和选择后的状态存储在redis
func CartKey(customerID int64) string {
	// cart:customerID, type:hash
	// value=productID:num, isCheck:1
	return "cart:" + strconv.FormatInt(customerID, 10)
}

func CheckedStatusField(productID int64) string {
	//以isCheck:productID形式做field
	return "isCheck:" + strconv.FormatInt(productID, 10)
}

// 从isCheck:productID解析出来的商品ID
var CheckedStatusFieldPrefix string = "isCheck:"

func ProductIDField(productID int64) string {
	return strconv.FormatInt(productID, 10)
}

func ProductKey(productID int64) string {
	//获取redis商品对应的key，值是数量
	// key=product:priductID
	// value=num
	return "product:" + strconv.FormatInt(productID, 10)
}
