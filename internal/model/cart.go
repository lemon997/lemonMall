package model

// isChecked表示是否选中的状态
// CartIds是传送过来的选中的商品id
// num和isCheck是存储在redis中
//分两次请求，根据cart:customer 获取对应的值，有以下关系：
// product_id:num
// ischeck:1
// 根据product_id获取DB中的name，price,URL
type CartResponse struct {
	ProductID   int64  `json:"product_id"`
	Num         int64  `json:"num"`
	ProductName string `json:"product_name"`
	ImgURL      string `json:"img_url"`
	Price       string `json:"price"`
	IsChecked   bool   `json:"is_checked"`
}

type CartRequest struct {
	ProductID int64   `json:"product_id"`
	Num       int64   `json:"num"`
	CartIds   []int64 `json:"cart_ids"`
}
