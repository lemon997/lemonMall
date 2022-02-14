package dao_test

import (
	"fmt"
	"testing"

	"github.com/lemon997/lemonMall/internal/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// func TestProductsInsert(t *testing.T) {
// 	dsn := "root:root@tcp(192.168.172.10:3306)/shop?charset=utf8mb4&parseTime=True&loc=Local"
// 	db, err := sqlx.Connect("mysql", dsn)

// 	if err != nil {
// 		t.Error(err)
// 	}
// 	pro := model.Products{}
// 	pro.CategoryChildrenID = 6
// 	pro.RemainingAmount = 10000
// 	pro.ProductName = `LAB SERIES保湿修护爽肤水（蓝宝瓶清爽补水平衡舒缓）`
// 	pro.Price = `260.00`
// 	pro.ImgUrl = `https://ssl3.sephorastatic.cn/products/V/S/1/0/2/6/7/9/6/1_n_12645_350x350.jpg`

// 	sqlStr := `insert into products(category_children_id, remaining_amount, product_name, price, img_url) values(?,?,?,?,?)`
// 	_, err = db.Exec(sqlStr, pro.CategoryChildrenID, pro.RemainingAmount, pro.ProductName, pro.Price, pro.ImgUrl)

// 	if err != nil {
// 		t.Error(err)
// 	}

// }

func TestProductsSelect(t *testing.T) {
	dsn := "root:root@tcp(192.168.172.10:3306)/shop?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sqlx.Connect("mysql", dsn)

	if err != nil {
		t.Error(err)
	}

	sqlStr := `select product_id, category_children_id, product_name, price, img_url from products where category_children_id = ? `
	adrs := []model.Products{}
	categoryChildrenId := 6
	err = db.Select(&adrs, sqlStr, categoryChildrenId)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(adrs)
}
