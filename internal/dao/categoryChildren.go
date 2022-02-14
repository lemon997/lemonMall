package dao

import (
	"context"

	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/model"
)

func CategoryChildrenTableName() string {
	return "category_children_children"
}

type CategoryChildrenMethod struct {
	ctx context.Context
}

func NewCategoryChildrenMethod(c context.Context) CategoryChildrenMethod {
	return CategoryChildrenMethod{ctx: c}
}

func (c CategoryChildrenMethod) Insert(name string, categoryID int64) error {

	sqlStr := `insert into category_children (category_id, category_children_name) values(?,?)`
	_, err := global.DBEngineShop.Exec(sqlStr, categoryID, name)
	return err
}

func (c CategoryChildrenMethod) Get(categoryChildrenID int64) (model.CategoryChildren, error) {

	sqlStr := `select category_children_id, category_id, category_children_name from category_children where category_children_id = ? `
	categ := model.CategoryChildren{}
	err := global.DBEngineShop.Get(&categ, sqlStr, categoryChildrenID)
	return categ, err
}

func (c CategoryChildrenMethod) Select(categoryID int64) ([]model.CategoryChildren, error) {
	sqlStr := `select category_children_id, category_id, category_children_name from category_children where category_id = ?`
	adrs := []model.CategoryChildren{}
	err := global.DBEngineShop.Select(&adrs, sqlStr, categoryID)
	return adrs, err
}

func (c CategoryChildrenMethod) Del(categoryChildrenID int64) (err error) {
	tx, err := global.DBEngineShop.Beginx()
	if err != nil {
		return
	}
	//defer的err针对两句DB执行语句
	defer func() {
		if p := recover(); p != nil {
			global.Logger.Errorf(c.ctx, "dao.categoryChildrenTable.Del, panic: %v", p)
			err = tx.Rollback()
			return
		} else if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}()

	sqlStr1 := `SELECT version FROM category_children WHERE category_children_id = ?`
	sqlStr2 := `DELETE FROM category_children WHERE category_children_id = ? AND version = ?`

	adr := model.CategoryChildren{}
	if err = tx.Get(&adr, sqlStr1, categoryChildrenID); err != nil {
		return
	}

	if _, err = tx.Exec(sqlStr2, categoryChildrenID, adr.Version); err != nil {
		return
	}
	return
}
