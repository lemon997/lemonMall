package user

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lemon997/lemonMall/common/app"
	"github.com/lemon997/lemonMall/common/errcode"
	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/model"
	"github.com/lemon997/lemonMall/internal/service"
)

func AddCart(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	var requestBody model.CartRequest //request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		response.ToErrorResponse(errcode.RequestParamsError) // 解析错误，参数有误
		return
	}

	customerId, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GetCustomerIdFaile)
		return
	}

	err := svc.AddCartProduct(customerId.(int64), requestBody.ProductID, requestBody.Num)
	if err != nil {
		global.Logger.Infof(ctx, "svc.AddCartProduct(),err= %v,customerID= %v, productID= %v",
			err, customerId.(int64), requestBody.ProductID)
		response.ToErrorResponse(errcode.WriteRedisError)
		return
	}

	response.ToResponse(gin.H{
		"status": 0,
		"msg":    "加入购物车成功",
	})
	return
}

func SetCartProductNum(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	productId, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	var requestBody model.CartRequest //request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		response.ToErrorResponse(errcode.RequestParamsError) // 解析错误，参数有误
		return
	}

	customerId, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GetCustomerIdFaile)
		return
	}

	err := svc.SetCartProductNum(customerId.(int64), productId, requestBody.Num)
	if err != nil {
		global.Logger.Infof(ctx, "svc.AddCartProduct(),err= %v,customerID= %v, productID= %v",
			err, customerId.(int64), requestBody.ProductID)
		response.ToErrorResponse(errcode.WriteRedisError)
		return
	}

	response.ToResponse(gin.H{
		"status": 0,
		"msg":    "修改成功",
	})
	return
}

func ChangeCheckedStatus(c *gin.Context) {
	//修改选择状态
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	var requestBody model.CartRequest //request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		response.ToErrorResponse(errcode.RequestParamsError) // 解析错误，参数有误
		return
	}

	customerId, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GetCustomerIdFaile)
		return
	}

	svc.ChangeCheckedStatus(customerId.(int64), requestBody.CartIds)
	response.ToResponse(gin.H{
		"status": 0,
	})
	return
}

func GetCartList(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	customerId, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GetCustomerIdFaile)
		return
	}

	res, _ := svc.GetCartProductList(customerId.(int64))
	response.ToResponse(gin.H{
		"data": res,
	})
	return
}

func DelCart(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	customerId, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GetCustomerIdFaile)
		return
	}

	productId, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := svc.DelCartProduct(customerId.(int64), productId)
	if err != nil {
		global.Logger.Errorf(ctx, "svc.DelCart,err= %v, customerID= %v, productID= %v",
			err, customerId.(int64), productId)
		response.ToErrorResponse(errcode.ErrorDelFail)
		return
	}
	response.ToResponse(gin.H{
		"status": 0,
	})
	return
}
