package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/lemon997/lemonMall/common/convert"

	"github.com/lemon997/lemonMall/common/util"

	"github.com/lemon997/lemonMall/global"

	"github.com/lemon997/lemonMall/internal/model"
)

type OrderMethod struct {
	ctx context.Context
}

func NewOrderMethod(c context.Context) OrderMethod {
	return OrderMethod{ctx: c}
}

func (o OrderMethod) Get(customerID, orderNo int64) (*model.Order, error) {
	//入参：customer_order_id的自增id
	//返回值：订单信息和error
	//描述：根据customer_order_id，获取信息
	sqlStr := `select customer_order_id, customer_id, customer_addr_id, product_id, num, order_no, express_no, express_type, pay_type, pay_time, amount, order_status from customer_order where customer_id = ? AND order_no = ?`
	order := model.Order{}
	if err := global.DBEngine.Get(&order, sqlStr, customerID, orderNo); err != nil {
		return nil, err
	}
	return &order, nil
}

func (o OrderMethod) Select(customerID, startIndex, endIndex int64) (*[]model.Order, error) {
	//描述根据customerID获取所有信息
	sqlStr := `select customer_order_id, customer_id, order_no, express_no, express_type, pay_type, pay_time, amount, order_status from customer_order where customer_id = ? limit ? offset ?`
	order := []model.Order{}
	if err := global.DBEngine.Select(&order, sqlStr, customerID, endIndex-startIndex+1, startIndex); err != nil {
		return nil, err
	}
	return &order, nil
}

func (o OrderMethod) SelectOrderList(customerID, startIndex, endIndex, orderStatus int64) ([]model.OrderListResponse, error) {
	if orderStatus == 0 {
		sqlStr1 := `select customer_order_id, customer_id, product_id, order_no, num, amount from customer_order where customer_id = ? limit ? offset ?`
		l := []model.OrderListResponse{}
		err := global.DBEngine.Select(&l, sqlStr1, customerID, endIndex-startIndex+1, startIndex)
		return l, err
	}
	sqlStr := `select customer_order_id, customer_id, product_id, order_no, num, amount from customer_order where customer_id = ? and order_status = ?  limit ? offset ?`
	l := []model.OrderListResponse{}
	err := global.DBEngine.Select(&l, sqlStr, customerID, orderStatus, endIndex-startIndex+1, startIndex)
	return l, err
}

func (o OrderMethod) Insert(customerID, addrID, productID, orderNO, num int64) (err error) {
	tx, err := global.DBEngine.Beginx()
	if err != nil {
		return
	}
	//defer的err针对两句DB执行语句
	defer func() {
		if p := recover(); p != nil {
			global.Logger.Errorf(o.ctx, "dao.order.Insert, panic: %v", p)
			err = tx.Rollback()
			return
		} else if err != nil {
			_ = tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}()

	sqlStr1 := `SELECT order_no FROM customer_order WHERE order_no = ? AND customer_id = ?`
	sqlStr2 := `insert into customer_order (customer_id, customer_addr_id, product_id, num, order_no, order_status) values(?,?,?,?,?,1)`

	order := model.Order{}
	err = tx.Get(&order, sqlStr1, orderNO, customerID)
	if err == nil || order.OrderNo > 0 {
		//默认是0,所以不能判断等于0,判断0则无法写入DB
		return errors.New("订单已经存在")
	}

	if _, err = tx.Exec(sqlStr2, customerID, addrID, productID, num, orderNO); err != nil {
		return
	}
	return
}

func (o OrderMethod) GetStatus(customerID, orderNo int64) (model.Order, error) {
	//获取订单状态
	sqlStr := `select order_status from customer_order where customer_id = ? AND order_no = ?`
	order := model.Order{}
	err := global.DBEngine.Get(&order, sqlStr, customerID, orderNo)
	return order, err
}

func (o OrderMethod) SetStatus(customerID, orderNo int64, status uint8) error {
	sqlStr := `update customer_order set order_status = ? where customer_id = ? AND order_no = ?`
	_, err := global.DBEngine.Exec(sqlStr, status, customerID, orderNo)
	return err
}

//帐号余额写在customer_login表中，这个金钱应该放在单独的金融表
func (o OrderMethod) SetPay(customerID int64, orderNo, prices []int64, payType string) (err error) {
	tx, err := global.DBEngine.Beginx()
	defer func() {
		if r := recover(); r != nil {
			global.Logger.Errorf(o.ctx, "panic in dao.order.SetPay panic= %v", r)
			err = tx.Rollback()
			return
		} else if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}()

	var money int64
	for i := 0; i < len(prices); i++ {
		if prices[i] <= 0 {
			return fmt.Errorf("money must > 0")
		}
		money += prices[i]
	}

	//查询多个订单状态,如果订单状态不为1,则退出,如果为1,则查询余额,
	sqlStr3 := `select version, order_status from customer_order where customer_id = ? AND order_no = ? for update`
	arr := make([]model.Order, len(orderNo))

	for i := 0; i < len(orderNo); i++ {
		version := model.Order{}
		err = tx.Get(&version, sqlStr3, customerID, orderNo[i])
		if err != nil {
			global.Logger.Errorf(o.ctx, "model.Order,err=%v", err)
			return
		}
		if version.OrderStatus != uint8(1) {
			return fmt.Errorf("already purchased")
		}
		arr[i] = version
	}

	//查询余额,如果不够则退出,够则扣款
	sqlStr1 := `select version, balance from customer_login where customer_id = ? for update`
	login := model.Login{}
	err = tx.Get(&login, sqlStr1, customerID)
	if err != nil {
		global.Logger.Errorf(o.ctx, "model.login,err=%v", err)
		return
	}

	//余额
	balance, err := convert.StrTo(login.Balance).MustInt64()
	if err != nil {
		global.Logger.Errorf(o.ctx, "StrTo,err=%v", err)
		return
	}
	//余额减总价格
	if balance-money < 0 {
		global.Logger.Errorf(o.ctx, "money=%d, balance=%d", money, balance)
		return fmt.Errorf("insufficient balance")
	}

	sqlStr2 := `update customer_login set version = version + 1, balance = balance - ? where customer_id = ? AND version = ?`
	_, err = tx.Exec(sqlStr2, money, customerID, login.Version)
	if err != nil {
		global.Logger.Errorf(o.ctx, "exec money,err=%v", err)
		return
	}

	sqlStr4 := `update customer_order set version = version + 1, amount = ?, pay_type = ?, order_status = 2, pay_time = ? where customer_id = ? AND order_no = ? AND version = ?`
	for i := 0; i < len(arr); i++ {
		payTime := util.NowTimeWithCTS()
		_, err = tx.Exec(sqlStr4, prices[i], payType, payTime, customerID, orderNo[i], arr[i].Version)
		if err != nil {
			global.Logger.Errorf(o.ctx, "model.exec order,err=%v", err)
			return
		}
	}

	return
}
