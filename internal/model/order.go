package model

type Order struct {
	CustomerOrderID int64  `json:"customer_order_id" db:"customer_order_id"`
	CustomerID      int64  `json:"customer_id" db:"customer_id"`
	CustomerAddrId  int64  `json:"customer_addr_id" db:"customer_addr_id"`
	ProductID       int64  `json:"product_id" db:"product_id"`
	Num             int64  `json:"num" db:"num"`
	Version         int64  `json:"-" db:"version"`
	OrderNo         int64  `json:"order_no" db:"order_no"`
	ExpressNo       int64  `json:"express_no" db:"express_no"`
	ExpressType     string `json:"express_type" db:"express_type"`
	PayType         string `json:"pay_type" db:"pay_type"`
	PayTime         string `json:"pay_time" db:"pay_time"`
	Amount          string `json:"amount" db:"amount"`
	CreateAt        string `json:"-" db:"create_at"`
	UpdateAt        string `json:"-" db:"update_at"`
	OrderStatus     uint8  `json:"order_status" db:"order_status"`
}

const orderTableNmae = "customer_order"

type OrderGoods struct {
	ProductID   int64  `json:"product_id" db:"product_id"`
	ProductName string `json:"product_name" db:"product_name"`
	Price       string `json:"price" db:"price"`
	ImgUrl      string `json:"img_url" db:"img_url"`
}

type OrderListResponse struct {
	CustomerOrderID int64      `json:"customer_order_id" db:"customer_order_id"`
	CustomerID      int64      `json:"customer_id" db:"customer_id"`
	ProductID       int64      `json:"product_id" db:"product_id"`
	OrderNo         int64      `json:"order_no" db:"order_no"`
	Num             int64      `json:"num" db:"num"`
	Amount          string     `json:"amount" db:"amount"`
	Goods           OrderGoods `json:"goods"`
	OrderStatus     uint8      `json:"order_status" db:"order_status"`
}

// 其中orderStatus对比：
// 0，全部
// 1，待付款
// 2，待发货
// 3，待收款
// 4，待评价
// 5，已过期

type OrderDetailInfoResponse struct {
	OrderDetail Order        `json:"order_detail"`
	Goods       []OrderGoods `json:"goods"`
}

type SettlementDataRepsonse struct {
	Num   int64      `json:"num"`
	Goods OrderGoods `json:"goods"`
}

type SettlementProductRequest struct {
	//结算时接受的商品数组
	SettlementProductIds []int64 `json:"settlement_product_ids"`
}

//提交订单的请求
type SubmitOrderRequest struct {
	AddrId     int64   `json:"addr_id"`
	ProductIds []int64 `json:"product_ids"`
}

type BlancePayRequest struct {
	OrderNos []int64 `json:"order_nos"`
	Price    []int64 `json:"price"`
	PayType  string  `json:"pay_type"`
}
