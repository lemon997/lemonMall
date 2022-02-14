package dao_test

import (
	"fmt"
	"testing"

	"github.com/lemon997/lemonMall/internal/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func TestInsert(t *testing.T) {
	// dsn := "user:password@tcp(192.168.172.10:3306)/shop?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := sqlx.Connect("mysql", dsn)
	// if err != nil {
	// 	t.Error(err)
	// }
	// name := "男士"
	// sqlStr := `insert into category (category_name) values(?)`
	// _, err = db.Exec(sqlStr, name)
	// if err != nil {
	// 	t.Error(err)
	// }
}
func TestSelect(t *testing.T) {
	dsn := "user:password@tcp(192.168.172.10:3306)/shop?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		t.Error(err)
	}
	sqlStr := `select category_id, category_name from category where category_id > ? `
	adrs := []model.Category{}
	category_id := 0

	if err := db.Select(&adrs, sqlStr, category_id); err != nil {
		t.Error(err)
	}
	fmt.Println(adrs)
}
