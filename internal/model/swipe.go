package model

//轮播图管理，设计主键ID,
type SwipeImg struct {
	ImgId      int64  `json:"img_id" db:"img_id"`
	ImgUrl     string `json:"img_url" db:"img_url"`
	CreateTime string `json:"create_time" db:"create_time"`
	UpdateTime string `json:"update_time" db:"update_time"`
}

func SwipeTableName() string {
	return "swipe_imgs_url"
}
