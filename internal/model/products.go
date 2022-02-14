package model

type Products struct {
	ProductID          int64  `json:"product_id" db:"product_id",omitempty`
	CategoryChildrenID int64  `json:"category_children_id" db:"category_children_id",omitempty`
	RemainingAmount    int64  `json:"remaining_amount" db:"remaining_amount",omitempty`
	Version            int64  `json:"version" db:"version",omitempty`
	ProductName        string `json:"product_name" db:"product_name",omitempty`
	Price              string `json:"price" db:"price",omitempty`
	ImgUrl             string `json:"img_url" db:"img_url",omitempty`
	Description        string `json:"description" db:"description",omitempty`
	CreateAt           string `json:"-" db:"create_at",omitempty`
	UpdateAt           string `json:"-" db:"update_at",omitempty`
}
