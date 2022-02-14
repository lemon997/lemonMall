package model

type GoodDetailResponse struct {
	ProductID   int64  `json:"product_id" db:"product_id"`
	ProductName string `json:"product_name" db:"product_name"`
	Price       string `json:"price" db:"price"`
	ImgUrl      string `json:"img_url" db:"img_url"`
	Description string `json:"description" db:"description"`
	Details     string `json:"details" db:"details"`
	IsRecommend int8   `json:"is_recommend" db:"is_recommend"`
}
