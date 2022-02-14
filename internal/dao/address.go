package dao

import (
	"context"

	"github.com/lemon997/lemonMall/common/util"

	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/model"
)

func AddressTableName() string {
	return "customer_addr"
}

//该数据表中，customer_id是不可能改动，只能根据customer_id建立索引

// 增删改查
type AddressMethod struct {
	ctx context.Context
}

func NewAddressMethod(c context.Context) AddressMethod {
	return AddressMethod{ctx: c}
}

func (a AddressMethod) Insert(addr model.Address) error {
	//入参：传入地址管理结构体
	//返回值：error
	//描述：插入一个地址管理行
	sqlStr := `insert into customer_addr (is_default,customer_id,name,phone,province,city,county,adr) values(?,?,?,?,?,?,?,?)`
	_, err := global.DBEngine.Exec(sqlStr, addr.IsDefault, addr.CustomerId, addr.Name, addr.Phone, addr.Province, addr.City, addr.County, addr.Adr)

	// num, _ := res.LastInsertId()
	// if num > 0 || err == nil {
	// 	return nil
	// }

	return err
}

func (a AddressMethod) Get(primaryKeyId, id int64) (model.Address, error) {
	//入参：customer_addr_id的自增id
	//返回值：地址管理信息和error
	//描述：根据customer_addr_id，获取地址管理信息
	sqlStr := `select customer_addr_id, is_default, name, phone, province, city, county, adr from customer_addr where customer_addr_id = ? AND customer_id = ?`
	adrr := model.Address{}

	if err := global.DBEngine.Get(&adrr, sqlStr, primaryKeyId, id); err != nil {
		return model.Address{}, err
	}
	return adrr, nil
}

func (a AddressMethod) Select(jwtId int64) ([]model.Address, error) {
	//入参：jwt的customer_id
	//返回值：地址管理信息切片和error
	//描述：根据customer_id，获取地址管理信息
	sqlStr := `select customer_addr_id, is_default, name, phone, province, city, county, adr from customer_addr where customer_id = ?`
	adrs := []model.Address{}
	if err := global.DBEngine.Select(&adrs, sqlStr, jwtId); err != nil {
		return adrs, err
	}
	return adrs, nil
}

func (a AddressMethod) Del(primaryKeyId int64) (err error) {
	//返回值：error,如果发生panic, 会捕捉这个错误，同时回滚，如果是事务错误，则会返回回滚的err或者提交的err
	//入参：id是swipe_adrs_url的自增主键id值,根据id进行查询和删除
	//描述：先查询目标id是否存在，存在则删除，需要开启事务
	tx, err := global.DBEngine.Beginx()
	if err != nil {
		return
	}
	//defer的err针对两句DB执行语句
	defer func() {
		if p := recover(); p != nil {
			global.Logger.Errorf(a.ctx, "dao.addressTable.Del, panic: %v", p)
			err = tx.Rollback()
			return
		} else if err != nil {
			// global.Logger.Infof(a.ctx, "dao.addressTable.Del, tx, err: %v", err)
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}()

	sqlStr1 := `SELECT update_time FROM customer_addr WHERE customer_addr_id = ?`
	sqlStr2 := `DELETE FROM customer_addr WHERE customer_addr_id = ? AND update_time = ?`

	adr := model.Address{}
	if err = tx.Get(&adr, sqlStr1, primaryKeyId); err != nil {
		return
	}

	//时间格式化
	t, err := util.MySQLTimeToGoTime(adr.UpdateTime)
	if err != nil {
		return
	}

	if _, err = tx.Exec(sqlStr2, primaryKeyId, t); err != nil {
		return
	}
	return
}

func (a AddressMethod) Modify(addr model.Address) (err error) {
	//PUT接口，要求幂等,采用乐观锁
	//返回值：error,如果发生panic, 会捕捉这个错误，同时回滚，如果是事务错误，则会返回回滚的err或者提交的err
	//入参：id是swipe_adrs_url的自增主键id值, 根据id进行查询和修改， newUrl是准备替换后的值
	//描述：先查询目标id是否存在，存在则使用乐观锁，对该id的img_url进行修改，需要开启事务
	tx, err := global.DBEngine.Beginx()
	if err != nil {
		return
	}
	//defer的err针对两句DB执行语句
	defer func() {
		if p := recover(); p != nil {
			global.Logger.Errorf(a.ctx, "dao.addressTable.Modify, panic: %v", p)
			_ = tx.Rollback()
			return
		} else if err != nil {
			// global.Logger.Infof(a.ctx, "dao.addressTable.Modify, err: %v", err)
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}()

	sqlStr1 := `SELECT update_time FROM customer_addr WHERE customer_addr_id = ?`
	sqlStr2 := `UPDATE customer_addr SET is_default=?,name=?,phone=?,province=?,city=?,county=?,adr=? WHERE customer_addr_id=? AND update_time=?`

	updateTime := model.Address{}
	if err = tx.Get(&updateTime, sqlStr1, addr.CustomerAddrId); err != nil {
		return
	}

	t, err := util.MySQLTimeToGoTime(updateTime.UpdateTime)
	if err != nil {
		return
	}

	if _, err = tx.Exec(sqlStr2, addr.IsDefault, addr.Name, addr.Phone, addr.Province,
		addr.City, addr.County, addr.Adr, addr.CustomerAddrId, t); err != nil {
		return
	}
	return
}

func (a AddressMethod) SetDefault(primaryKeyId int64, customerId int64) (err error) {
	tx, err := global.DBEngine.Beginx()
	if err != nil {
		return
	}
	//defer的err针对两句DB执行语句
	defer func() {
		if p := recover(); p != nil {
			global.Logger.Errorf(a.ctx, "dao.addressTable.SetDefault, panic: %v", p)
			_ = tx.Rollback()
			return
		} else if err != nil {
			// global.Logger.Infof(a.ctx, "dao.addressTable.SetDefault, err: %v", err)
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}()

	sqlStr1 := `SELECT customer_addr_id FROM customer_addr WHERE customer_id = ? AND is_default = 1`
	sqlStr2 := `UPDATE customer_addr SET is_default = 1 where customer_addr_id = ? `
	sqlStr3 := `UPDATE customer_addr SET is_default = 0 where customer_addr_id = ? `

	oldDefaultId := model.Address{}
	if err = tx.Get(&oldDefaultId, sqlStr1, customerId); err != nil {
		return
	}

	if _, err = tx.Exec(sqlStr2, primaryKeyId); err != nil {
		return
	}
	if _, err = tx.Exec(sqlStr3, oldDefaultId.CustomerAddrId); err != nil {
		return
	}
	return
}
