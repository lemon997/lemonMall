package dao

import (
	"context"

	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/model"
)

func CollectTableName() string {
	return "customer_collect"
}

type CollectMethod struct {
	ctx context.Context
}

func NewCollectMethod(c context.Context) CollectMethod {
	return CollectMethod{ctx: c}
}

func (c CollectMethod) Insert(collect model.Collect) error {
	sqlStr := `insert into customer_collect (customer_id,product_id,product_name,img_url) values(?,?,?,?)`
	_, err := global.DBEngine.Exec(sqlStr, collect.CustomerID, collect.ProductID, collect.ProductName, collect.ImgURL)
	return err
}

func (c CollectMethod) Get(productID, customerID int64) (model.Collect, error) {

	sqlStr := `select customer_collect_id from customer_collect where customer_id = ? AND product_id = ?`
	categ := model.Collect{}
	if err := global.DBEngine.Get(&categ, sqlStr, customerID, productID); err != nil {
		return model.Collect{}, err
	}
	return categ, nil
}

func (c CollectMethod) Select(customerID, startIndex, endIndex int64) ([]model.Collect, error) {
	sqlStr := `select customer_collect_id, customer_id, product_id, product_name, img_url from customer_collect where customer_id = ? limit ? offset ?`
	adrs := []model.Collect{}
	//加入endIndex是9,但是limit是输出多少个，输出10个需要endIndex+1
	if err := global.DBEngine.Select(&adrs, sqlStr, customerID, endIndex-startIndex+1, startIndex); err != nil {
		return adrs, err
	}
	return adrs, nil
}

func (c CollectMethod) Del(productID, customerID int64) (err error) {
	tx, err := global.DBEngine.Beginx()
	if err != nil {
		return
	}
	//defer的err针对两句DB执行语句
	defer func() {
		if p := recover(); p != nil {
			global.Logger.Errorf(c.ctx, "dao.collectTable.Del, panic: %v", p)
			err = tx.Rollback()
			return
		} else if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}()

	sqlStr1 := `SELECT version FROM customer_collect WHERE customer_id = ? AND product_id = ?`
	sqlStr2 := `DELETE FROM customer_collect WHERE customer_id = ? AND product_id = ? AND version = ?`

	adr := model.Collect{}
	if err = tx.Get(&adr, sqlStr1, customerID, productID); err != nil {
		return
	}

	if _, err = tx.Exec(sqlStr2, customerID, productID, adr.Version); err != nil {
		return
	}
	return
}
