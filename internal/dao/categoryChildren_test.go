package dao_test

import (
	"fmt"
	"testing"

	"github.com/lemon997/lemonMall/internal/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func TestChildrenInsert(t *testing.T) {
	// dsn := "user:password@tcp(192.168.172.10:3306)/shop?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := sqlx.Connect("mysql", dsn)

	// if err != nil {
	// 	t.Error(err)
	// }
	// name := "清洁"
	// var categoryID int64 = 3

	// sqlStr := `insert into category_children (category_id, category_children_name) values(?,?)`

	// _, err = db.Exec(sqlStr, categoryID, name)
	// if err != nil {
	// 	t.Error(err)
	// }
	// name = "面部护理"
	// categoryID = 3

	// sqlStr = `insert into category_children (category_id, category_children_name) values(?,?)`

	// _, err = db.Exec(sqlStr, categoryID, name)
	// if err != nil {
	// 	t.Error(err)
	// }
}
func TestChildrenSelect(t *testing.T) {
	dsn := "user:password@tcp(192.168.172.10:3306)/shop?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sqlx.Connect("mysql", dsn)

	if err != nil {
		t.Error(err)
	}

	sqlStr := `select category_children_id, category_id, category_children_name from category_children where category_id = ?`
	adrs := []model.CategoryChildren{}
	categoryID := 3
	err = db.Select(&adrs, sqlStr, categoryID)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(adrs)
}
