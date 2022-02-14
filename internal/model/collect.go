package model

type Collect struct {
	CustomerCollectID int64  `json:"customer_collect_id" db:"customer_collect_id"`
	CollectCategory   string `json:"-" db:"collect_category"`
	CustomerID        int64  `json:"customer_id" db:"customer_id"`
	ProductID         int64  `json:"product_id" db:"product_id"`
	ProductName       string `json:"product_name" db:"product_name"`
	ImgURL            string `json:"img_url" db:"img_url"`
	Version           int64  `json:"-" db:"version"`
	CreateAt          string `json:"-" db:"create_at"`
	UpdateAt          string `json:"-" db:"update_at"`
}
