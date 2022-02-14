package model

type Address struct {
	//province省，city市，county区，adr详细地址
	CustomerAddrId int64  `json:"customer_addr_id" db:"customer_addr_id"`
	CustomerId     int64  `json:"customer_id" db:"customer_id"`
	Name           string `json:"name" db:"name"`
	Phone          string `json:"phone" db:"phone"`
	Province       string `json:"province" db:"province"`
	City           string `json:"city" db:"city"`
	County         string `json:"county" db:"county"`
	Adr            string `json:"adr" db:"adr"`
	CreateTime     string `json:"-" db:"create_time"`
	UpdateTime     string `json:"-" db:"update_time"`
	IsDefault      int8   `json:"is_default" db:"is_default"`
}
