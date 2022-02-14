package model

type CategoryChildren struct {
	CategoryChildrenID   int64  `json:"category_children_id" db:"category_children_id"`
	CategoryID           int64  `json:"category_id" db:"category_id"`
	CategoryChildrenName string `json:"category_children_name" db:"category_children_name"`
	CreateAt             string `json:"-" db:"create_at"`
	UpdateAt             string `json:"-" db:"update_at"`
	Version              int64  `json:"-" db:"version"`
}
