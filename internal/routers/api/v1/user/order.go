package user

import (
	"github.com/gin-gonic/gin"
	"github.com/lemon997/lemonMall/common/app"
	"github.com/lemon997/lemonMall/common/convert"
	"github.com/lemon997/lemonMall/common/errcode"
	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/model"
	"github.com/lemon997/lemonMall/internal/service"
)

func GetOrderList(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	custId, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GetCustomerIdFaile)
		return
	}

	orderStatus, _ := convert.StrTo(c.Query("status")).MustInt64()
	if orderStatus < 0 || orderStatus > 5 {
		response.ToErrorResponse(errcode.ErrorSelectFail)
		return
	}

	page := app.GetPage(c)
	if page <= 0 {
		page = 1
	}
	startIndex, endIndex := app.GetPageInterval(page)

	data, err := svc.GetOrderList(custId.(int64), startIndex, endIndex, orderStatus)
	if err != nil {
		global.Logger.Errorf(ctx, "SelectOrder err= %q, customerID= %d, startIndex= %d, endIndex= %d",
			err, custId.(int64), startIndex, endIndex)
		response.ToErrorResponse(errcode.ErrorSelectFail)
		return
	}
	response.ToResponse(gin.H{
		"status": 0,
		"msg":    "获取成功",
		"data":   data,
	})
	return
}

func GetOrderDetail(c *gin.Context) {
	//返回参数：商品信息，订单信息，订单状态
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	orderID, _ := convert.StrTo(c.Param("id")).MustInt64()

	custId, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GetCustomerIdFaile)
		return
	}

	data, err := svc.GetOrderDetail(orderID, custId.(int64))
	if err != nil {
		global.Logger.Errorf(ctx, "user.order.GetOrderDetail err= %q, customer_no= %d",
			err, orderID)
		response.ToErrorResponse(errcode.ErrorSelectFail)
		return
	}
	response.ToResponse(gin.H{
		"status": 0,
		"msg":    "获取成功",
		"data":   *data,
	})
	return
}

func GetSettlementData(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	custId, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GetCustomerIdFaile)
		return
	}

	//settlementProductIDs是需要购买的商品ID
	var settlementProductIDs model.SettlementProductRequest
	if err := c.ShouldBindJSON(&settlementProductIDs); err != nil {
		response.ToErrorResponse(errcode.InvalidParams)
		global.Logger.Errorf(ctx, "user.order.GetSettlementData.c.ShouldBindJSON, err= %v", err)
		return
	}

	addrs, err := svc.SelectAddress(custId.(int64))
	if err != nil {
		response.ToErrorResponse(errcode.ServerError)
		global.Logger.Errorf(ctx, "user.address.GetAllAddress.SelectAddress, custId: %d ,err: %v",
			custId.(int64), err)
		return
	}

	// 获取商品信息和数量
	nums, err := svc.SelectProductNums(custId.(int64), settlementProductIDs.SettlementProductIds)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorSelectFail)
		global.Logger.Errorf(ctx, "user.GetSettementData.SelectProductNums, err= %v", err)
		return
	}

	//根据商品ID获取商品信息
	err = svc.SelectProductsWithSettlementDataResponse(nums)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorGetGoodDate)
		global.Logger.Errorf(ctx, "user.GetSettementData.SelectProductsWithSettlementDataResponse, err= %v", err)
		return
	}

	response.ToResponse(gin.H{
		"status":       0,
		"msg":          "获取地址管理信息成功",
		"access_addrs": addrs,
		"goods":        *nums,
	})
	return
}

func SubmitOrder(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	custId, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GetCustomerIdFaile)
		return
	}

	var order model.SubmitOrderRequest
	if err := c.ShouldBindJSON(&order); err != nil {
		response.ToErrorResponse(errcode.InvalidParams)
		global.Logger.Errorf(ctx, "user.order.SubmitOrder.c.ShouldBindJSON, err= %v", err)
		return
	}

	orderNOs, err := svc.PlaceOrder(custId.(int64), order.AddrId, order.ProductIds)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorOrder)
		global.Logger.Errorf(ctx, "user.order.SubmitOrder.svc.PlaceOrder, err= %v", err)
		return
	}

	response.ToResponse(gin.H{
		"status":    0,
		"msg":       "下单成功",
		"order_nos": orderNOs,
	})
	return

}

func Payment(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	custId, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GetCustomerIdFaile)
		return
	}

	var pay model.BlancePayRequest
	if err := c.ShouldBindJSON(&pay); err != nil {
		response.ToErrorResponse(errcode.InvalidParams)
		global.Logger.Errorf(ctx, "user.order.BlancePay.c.ShouldBindJSON, err= %v", err)
		return
	}

	//接受到实际付款的金额数组，订单号数组，支付类型，
	// 就要修改MySQL的pay_type, express_type, pay_time, amount, order_status, version, express_no
	// 还要删除延时队列的订单
	if len(pay.OrderNos) != len(pay.Price) {
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	err := svc.Payment(custId.(int64), pay.OrderNos, pay.Price, pay.PayType)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorPay)
		global.Logger.Errorf(ctx, "user.order.payment.payment, err= %v", err)
		return
	}
	response.ToResponse(gin.H{
		"status": 0,
		"msg":    "购买成功",
	})
	return
}
