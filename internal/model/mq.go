package model

//通过customer和order_no查看订单表的状态，如果已经购买，则不做操作，没有购买，则根据productId返回库存
type DelayOrder struct {
	CustomerId int64   `json:"customer_id" db:"customer_id"`
	OrderNo    []int64 `json:"order_no" db:"order_no"`
	ProductId  []int64 `json:"product_id" db:"product_id"`
	Num        []int64 `json:"num" db:"num"`
}

//减库存
type ReduceInventory struct {
	ProductId []int64 `json:"product_id" db:"product_id"`
	Num       []int64 `json:"num" db:"num"`
}
