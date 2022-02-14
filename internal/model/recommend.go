package model

//推荐商品管理，设计主键ID,
type RecommendImg struct {
	ImgId      int64  `json:"img_id" db:"img_id"`
	ImgUrl     string `json:"img_url" db:"img_url"`
	CreateTime string `json:"create_time" db:"create_time"`
	UpdateTime string `json:"update_time" db:"update_time"`
}

func RecommendTableName() string {
	//数据库表名
	return "recommend_imgs_url"
}
