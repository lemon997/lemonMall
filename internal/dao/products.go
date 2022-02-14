package dao

import (
	"context"

	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/model"
)

func ProductsTableName() string {
	return "products"
}

type ProductsMethod struct {
	ctx context.Context
}

func NewProductsMethod(c context.Context) ProductsMethod {
	return ProductsMethod{ctx: c}
}

func (p ProductsMethod) Insert(pro model.Products) error {

	sqlStr := `insert into products(category_children_id, remaining_amount, product_name, price, img_url) values(?,?,?,?,?)`
	_, err := global.DBEngineShop.Exec(sqlStr, pro.CategoryChildrenID, pro.RemainingAmount, pro.ProductName, pro.Price, pro.ImgUrl)
	return err
}

func (p ProductsMethod) Get(productId int64) (model.Products, error) {

	sqlStr := `select product_name, price, img_url, description from products where product_id = ? `
	categ := model.Products{}
	err := global.DBEngineShop.Get(&categ, sqlStr, productId)
	return categ, err
}

func (p ProductsMethod) Select(categoryChildrenId int64, page ...int64) ([]model.Products, error) {
	//page[0]=startIndex,page[1]=endIndex
	if len(page) == 2 {
		sqlStr1 := `select product_id, product_name, price, img_url from products limit ? offset ?`
		res := []model.Products{}
		err := global.DBEngineShop.Select(&res, sqlStr1, page[1]-page[0]+1, page[0])
		return res, err
	}
	sqlStr := `select product_id, category_children_id, product_name, price, img_url from products where category_children_id = ? `
	adrs := []model.Products{}
	err := global.DBEngineShop.Select(&adrs, sqlStr, categoryChildrenId)
	return adrs, err
}

func (p ProductsMethod) Del(primaryKeyId int64) (err error) {
	tx, err := global.DBEngineShop.Beginx()
	if err != nil {
		return
	}
	//defer的err针对两句DB执行语句
	defer func() {
		if pan := recover(); pan != nil {
			global.Logger.Errorf(p.ctx, "dao.productsTable.Del, panic: %v", p)
			err = tx.Rollback()
			return
		} else if err != nil {
			// global.Logger.Infof(p.ctx, "dao.productsTable.Del, tx, err: %v", err)
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}()

	sqlStr1 := `SELECT version FROM products WHERE products_id = ?`
	sqlStr2 := `DELETE FROM products WHERE product_id = ? AND version = ?`

	adr := model.Products{}
	if err = tx.Get(&adr, sqlStr1, primaryKeyId); err != nil {
		return
	}

	if _, err = tx.Exec(sqlStr2, primaryKeyId, adr.Version); err != nil {
		return
	}
	return
}

func (p ProductsMethod) GetStock(productID int64) (model.Products, error) {
	sqlStr := `select remaining_amount from products where product_id = ?`
	stock := model.Products{}
	err := global.DBEngineShop.Get(&stock, sqlStr, productID)
	return stock, err
}

func (p ProductsMethod) ReduceInventory(productID, nums int64) (err error) {
	tx, err := global.DBEngineShop.Beginx()
	if err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			global.Logger.Errorf(context.TODO(), "panic in dao.products.ReduceInventory panic= %v", r)
			err = tx.Rollback()
			return
		}
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}()

	sqlStr1 := `select version from products where product_id = ? for update`
	sqlStr2 := `update products set version = version + 1, remaining_amount = remaining_amount - ? where product_id = ? AND version = ?`

	product := model.Products{}
	err = tx.Get(&product, sqlStr1, productID)

	if err != nil {
		return
	}

	_, err = tx.Exec(sqlStr2, nums, productID, product.Version)
	if err != nil {
		return
	}
	return
}

func (p ProductsMethod) AddInventory(productID, nums int64) (err error) {
	tx, err := global.DBEngineShop.Beginx()
	if err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			global.Logger.Errorf(context.TODO(), "panic in dao.products.AddInventory panic= %v", r)
			err = tx.Rollback()
			return
		}
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}()

	sqlStr1 := `select version from products where product_id = ? for update`
	sqlStr2 := `update products set version = version + 1, remaining_amount = remaining_amount + ? where product_id = ? AND version = ?`

	product := model.Products{}
	err = tx.Get(&product, sqlStr1, productID)

	if err != nil {
		return
	}

	_, err = tx.Exec(sqlStr2, nums, productID, product.Version)
	if err != nil {
		return
	}
	return
}
