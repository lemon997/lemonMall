package user

import (
	"github.com/lemon997/lemonMall/common/app"
	"github.com/lemon997/lemonMall/common/convert"
	"github.com/lemon997/lemonMall/common/errcode"
	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/service"

	"github.com/gin-gonic/gin"
)

func GetCollectList(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	page := app.GetPage(c)
	startIndex, endIndex := app.GetPageInterval(page)

	customerID, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GetCustomerIdFaile)
		return
	}

	res, err := svc.SelectCollectWithCustomerID(customerID.(int64), startIndex, endIndex)
	if err != nil {
		global.Logger.Errorf(ctx, "v1.user.collect.svc.GetCollectWithCustomerID, err= %v, customer_id= %v",
			err, customerID.(int64))
		response.ToErrorResponse(errcode.ServerError)
		return
	}
	response.ToResponse(gin.H{
		"status": 0,
		"msg":    "获取成功",
		"data":   res,
	})
}

func CancelCollect(c *gin.Context) {
	//取消收藏
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	customerID, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GetCustomerIdFaile)
		return
	}

	productID, err := convert.StrTo(c.Param("id")).MustInt64()
	if err != nil {
		global.Logger.Errorf(ctx, "v1.user.CancelCollect.StrTo,err= %v, product_id=%v", err, c.Param("id"))
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	err = svc.DelCollectWithCustomeID(productID, customerID.(int64))
	if err != nil {
		global.Logger.Infof(ctx, "v1.user.CancelCollect.svc.Del,err= %v, productID= %v, customerID= %v",
			err, productID, customerID.(int64))
		response.ToErrorResponse(errcode.ErrorDelFail)
		return
	}

	response.ToResponse(gin.H{
		"status": 0,
		"msg":    "删除成功",
	})
}

func SetCollect(c *gin.Context) {
	// 加入收藏
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	productID, err := convert.StrTo(c.Param("id")).MustInt64()
	if err != nil {
		global.Logger.Infof(ctx, "v1.user.CheckCollect.StrTo,err= %v, product_id=%v", err, c.Param("id"))
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	customerID, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GetCustomerIdFaile)
		return
	}

	err = svc.InsertCollectWithProductID(productID, customerID.(int64))
	if err != nil {
		global.Logger.Infof(ctx, "v1.user.SetCollect.svc.InsertCollectWithProductID,err= %v, productID= %v, customerID= %v",
			err, productID, customerID.(int64))
		response.ToErrorResponse(errcode.ServerError)
		return
	}

	response.ToResponse(gin.H{
		"status": 0,
		"msg":    "添加成功",
	})
}

func CheckCollect(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	//获取该productID是否已经收藏
	customerID, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GetCustomerIdFaile)
		return
	}
	productId, err := convert.StrTo(c.Param("id")).Int()
	if err != nil {
		global.Logger.Errorf(ctx, "v1.user.collect.CheckCollect.StrTo,err= %v, product_id=%v", err, c.Param("id"))
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	_, err = svc.GetCollectWithCustomerID(int64(productId), customerID.(int64))
	var num int
	if err == nil {
		num = 1
	}

	response.ToResponse(gin.H{
		"is_collect": num,
	})
}
