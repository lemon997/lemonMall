package dao

import (
	"context"

	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/model"
)

type LoginMethod struct{}

func (l LoginMethod) GetRow(name string) (model.Login, error) {
	//该函数根据用户名获取整行数据
	sqlStr := `SELECT * FROM customer_login WHERE login_name = ?;`

	//声明一个切片获取关于name的信息
	var linfo model.Login
	//分两种情况，1、数据库查询不到但没有报错，返回linfo为空
	if err := global.DBEngine.Get(&linfo, sqlStr, name); err != nil {
		return linfo, err
	}
	return linfo, nil
}

func (l LoginMethod) StoreDB(name, pwd string) (err error) {
	sqlStr := `insert into customer_login (login_name, password) values(?, ?)`
	res, err := global.DBEngine.Exec(sqlStr, name, pwd)
	num, _ := res.LastInsertId()
	if num > 0 && err == nil {
		return nil
	}
	return err
}

func (l LoginMethod) UpdatePwd(ctx context.Context, pwd, name string) error {
	sqlStr := `UPDATE customer_login set password=?, version=version+1 where login_name=?;`
	_, err := global.DBEngine.Exec(sqlStr, pwd, name)
	if err != nil {
		//系统报错
		global.Logger.Errorf(ctx, "dao.loginDBTable.UpdatePwd，err: %v", err)
		return err
	}
	return nil

}

func (l LoginMethod) SetUrl(customerID int64, url string) error {
	sqlStr := `update customer_login set img_url = ? where customer_id = ?`
	_, err := global.DBEngine.Exec(sqlStr, url, customerID)
	return err
}

func (l LoginMethod) GetNameAndUrl(customerID int64) (model.Login, error) {
	sqlStr := `select login_name, img_url, balance from customer_login where customer_id = ?`
	login := model.Login{}
	err := global.DBEngine.Get(&login, sqlStr, customerID)
	return login, err
}
