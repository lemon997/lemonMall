package dao

import (
	"context"

	"github.com/lemon997/lemonMall/common/util"

	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/model"
)

func CategoryTableName() string {
	return "category"
}

type CategoryMethod struct {
	ctx context.Context
}

func NewCategoryMethod(c context.Context) CategoryMethod {
	return CategoryMethod{ctx: c}
}

func (c CategoryMethod) Insert(name string) error {
	sqlStr := `insert into category (category_name) values(?)`
	_, err := global.DBEngineShop.Exec(sqlStr, name)
	return err
}

func (c CategoryMethod) Get(primaryKeyId int64) (model.Category, error) {

	sqlStr := `select * from category where category_id = ? `
	categ := model.Category{}
	if err := global.DBEngineShop.Get(&categ, sqlStr, primaryKeyId); err != nil {
		return model.Category{}, err
	}
	return categ, nil
}

func (c CategoryMethod) Select(primaryKeyId int64) ([]model.Category, error) {
	sqlStr := `select category_id, category_name from category where category_id > ? `
	adrs := []model.Category{}

	if err := global.DBEngineShop.Select(&adrs, sqlStr, primaryKeyId); err != nil {
		return adrs, err
	}
	return adrs, nil
}

func (c CategoryMethod) Del(primaryKeyId int64) (err error) {
	tx, err := global.DBEngineShop.Beginx()
	if err != nil {
		return
	}
	//defer的err针对两句DB执行语句
	defer func() {
		if p := recover(); p != nil {
			global.Logger.Errorf(c.ctx, "dao.categoryTable.Del, panic: %v", p)
			err = tx.Rollback()
			return
		} else if err != nil {
			// global.Logger.Infof(c.ctx, "dao.categoryTable.Del, tx, err: %v", err)
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}()

	sqlStr1 := `SELECT update_at FROM category WHERE category_id = ?`
	sqlStr2 := `DELETE FROM category WHERE category_id = ? AND update_at = ?`

	adr := model.Category{}
	if err = tx.Get(&adr, sqlStr1, primaryKeyId); err != nil {
		return
	}

	//时间格式化
	t, err := util.MySQLTimeToGoTime(adr.UpdateAt)
	if err != nil {
		return
	}

	if _, err = tx.Exec(sqlStr2, primaryKeyId, t); err != nil {
		return
	}
	return
}
