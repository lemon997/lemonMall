package service

import (
	"encoding/json"
	"errors"

	"github.com/lemon997/lemonMall/internal/mq/producer"

	"github.com/lemon997/lemonMall/global"

	"github.com/lemon997/lemonMall/internal/dao"
	"github.com/lemon997/lemonMall/internal/model"
	"github.com/lemon997/lemonMall/internal/redisDao"
)

func (svc Service) GetOrder(customerID, orderNo int64) (*model.Order, error) {
	return dao.NewOrderMethod(svc.ctx).Get(customerID, orderNo)
}

func (svc Service) SelectOrder(customerID, startIndex, endIndex int64) (*[]model.Order, error) {
	return dao.NewOrderMethod(svc.ctx).Select(customerID, startIndex, endIndex)
}

func (svc Service) GetOrderList(customerID, startIndex, endIndex, orderStatus int64) ([]model.OrderListResponse, error) {
	l, err := dao.NewOrderMethod(svc.ctx).SelectOrderList(customerID, startIndex, endIndex, orderStatus)
	if len(l) == 0 {
		return l, nil
	}
	if err != nil {
		return l, err
	}
	p := dao.NewProductsMethod(svc.ctx)
	for i := 0; i < len(l); i++ {
		res, err := p.Get(l[i].ProductID)
		if err != nil {
			global.Logger.Errorf(svc.ctx, "GetOrderList err= %v", err)
			continue
		}
		l[i].Goods.ProductName = res.ProductName
		l[i].Goods.ImgUrl = res.ImgUrl
		l[i].Goods.Price = res.Price
	}
	return l, nil
}

func (svc Service) GetOrderDetail(orderNo, customerID int64) (*model.OrderDetailInfoResponse, error) {
	//根据订单编号，也就是order_no，获取信息
	order, err := dao.NewOrderMethod(svc.ctx).Get(customerID, orderNo)

	if err != nil {

		return nil, err
	}

	good, err := dao.NewProductsMethod(svc.ctx).Get(order.ProductID)
	if err != nil {

		return nil, err
	}
	goodDetail := model.OrderDetailInfoResponse{}
	goodDetail.OrderDetail = *order

	orderGood := model.OrderGoods{}
	orderGood.ImgUrl = good.ImgUrl
	orderGood.Price = good.Price
	orderGood.ProductName = good.ProductName

	goodDetail.Goods = append(goodDetail.Goods, orderGood)
	return &goodDetail, nil
}

func (svc Service) ReduceInventoryWithOrder(customerID int64, productID []int64) ([]int64, error) {
	//返回扣减的数量以及错误，还进行删除购物车对应键值对操作
	//因为提交订单的时候，购物数量和选择的东西在cart:customer上，前端数据不可信，做校验，
	// 前端的商品和cart:customerID 和 isCheck:productID校验
	//检查cart的选择项与购买的商品是否符合
	cartKey := model.CartKey(customerID)
	fields := make([]string, len(productID))
	//检查选择项
	for i := 0; i < len(productID); i++ {
		fields[i] = model.CheckedStatusField(productID[i])
	}

	exist, err := svc.ExistFieldsWithCart(cartKey, fields)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("下订单失败")
	}

	//获取购物车的商品数量
	num, err := redisDao.NewMethodsCartProductNum(svc.ctx).GetProductNumsAtCart(cartKey, productID)
	if num == nil || len(num) != len(productID) {
		return nil, errors.New("获取cart数量错误")
	}
	if err != nil {
		return nil, err
	}

	//进行减库存操作,keys的键形式是"product:productID"
	keys := make([]string, len(productID))
	for i := 0; i < len(productID); i++ {
		keys[i] = model.ProductKey(productID[i])
	}

	err = redisDao.NewPlaceOrderMethod(svc.ctx).ReduceInventory(keys, num)
	if err != nil {
		return nil, err
	}

	//删除对应键值对
	for i := 0; i < len(fields); i++ {
		_ = svc.DelCartField(cartKey, fields[i])
	}

	// 删除cart:customerID 下的field "productID"
	for i := 0; i < len(productID); i++ {
		field := model.ProductIDField(productID[i])
		err = svc.DelCartField(cartKey, field)
	}

	return num, nil
}

func (svc Service) PlaceOrder(customerID, addrID int64, productID []int64) ([]int64, error) {
	//返回订单编号
	//获取购物车商品数量，订单ID,地址id,用户id,将订单信息写入DB
	nums, err := svc.ReduceInventoryWithOrder(customerID, productID)
	if err != nil {
		return nil, err
	}

	//获取订单ID
	ordersID := make([]int64, len(productID))
	insert := dao.NewOrderMethod(svc.ctx)

	for i := 0; i < len(productID); i++ {
		orderNO := global.Node1.Generate().Int64()
		insert.Insert(customerID, addrID, productID[i], orderNO, nums[i])
		ordersID[i] = orderNO
	}

	//上面已经在mysql生成订单了，下面就发送消息给MySQL扣减库存以及延时队列
	delay := model.DelayOrder{}
	delay.CustomerId = customerID
	delay.OrderNo = ordersID
	delay.ProductId = productID
	delay.Num = nums
	data, _ := json.Marshal(delay)

	err = producer.DBDelayOrderSend(data)

	if err != nil {
		global.Logger.Errorf(svc.ctx, "send delay queue fail, err= %v", err)
		return nil, err
	}

	// 扣减库存
	reduce := model.ReduceInventory{}
	reduce.ProductId = productID
	reduce.Num = nums
	body, _ := json.Marshal(reduce)

	err = producer.DBReduceInventorySend(body)

	if err != nil {
		global.Logger.Errorf(svc.ctx, "send message queue fail, err= %v", err)
		return nil, err
	}

	return ordersID, err
}
