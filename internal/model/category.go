package model

type Category struct {
	CategoryID   int64  `json:"category_id" db:"category_id"`
	Version      int64  `json:"-" db:"version"`
	CategoryName string `json:"category_name" db:"category_name"`
	CreateAt     string `json:"-" db:"create_at"`
	UpdateAt     string `json:"-" db:"update_at"`
}

var CategoryArray = []int64{1, 2, 3}
