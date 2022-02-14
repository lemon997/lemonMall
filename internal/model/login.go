package model

//管理数据库的模型
type Login struct {
	UserStats    int8   `db:"user_stats" json:"userstats",omitempty`
	Version      int8   `db:"version" json:"version",omitempty`
	CustomerId   int64  `db:"customer_id" json:"customerid",omitempty`
	LoginName    string `db:"login_name" json:"loginname",omitempty`
	Password     string `db:"password" json:"password",omitempty`
	ModifiedTime string `db:"modified_time" json:"modifiedtime",omitempty`
	Balance      string `db:"balance" json:"balance",omitempty`
	ImgUrl       string `db:"img_url" json:"img_url",omitempty`
}
