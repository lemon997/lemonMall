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

func GetAllAddress(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	custId, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GetCustomerIdFaile)
		return
	}

	addrs, err := svc.SelectAddress(custId.(int64))
	if err != nil {
		response.ToErrorResponse(errcode.ServerError)
		global.Logger.Errorf(ctx, "user.address.GetAllAddress.SelectAddress, custId: %d ,err: %v", custId.(int64), err)
		return
	}

	response.ToResponse(gin.H{
		"status":       0,
		"msg":          "获取地址管理信息成功",
		"access_addrs": addrs,
	})
	return
}

func AddAddress(c *gin.Context) {
	//入参：收货人名，手机，收获地址

	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	requestBody := model.Address{} //request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		response.ToErrorResponse(errcode.RequestParamsError) // 解析错误，参数有误
		// global.Logger.Errorf(ctx, "user.address.AddAddress: %v", requestBody)
		return
	}

	//获取jwt的id
	id, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GinGetError)
		return
	}
	requestBody.CustomerId = id.(int64)

	err := svc.InsertAddress(requestBody)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError)
		global.Logger.Errorf(ctx, "user.address.AddAddress.InsertAddress, err: %v", err)
		return
	}
	response.ToResponse(gin.H{
		"status": 0,
		"msg":    "添加地址成功",
	})

}

func DelAddress(c *gin.Context) {
	//入参：customer_addr_id
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	addressId := c.Param("id")
	pkid, _ := strconv.Atoi(addressId)

	err := svc.DelAddress(int64(pkid))
	if err != nil {
		response.ToErrorResponse(errcode.ServerError)
		global.Logger.Errorf(ctx, "user.address.DelAddress.DelAddress, err: %v", err)
		return
	}
	response.ToResponse(gin.H{
		"status": 0,
		"msg":    "删除地址成功",
	})
}

func GetAddressDetail(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	addressId := c.Param("id") //获取path参数，也就是"/getaddressdetail/:id"的id
	pkid, _ := strconv.Atoi(addressId)

	id, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GinGetError)
		return
	}

	adr, err := svc.GetAddress(int64(pkid), id.(int64))
	if err != nil {
		response.ToErrorResponse(errcode.ServerError)
		global.Logger.Errorf(ctx, "user.GetAddressDetail.GetAddress, err= %v", err)
		return
	}

	response.ToResponse(gin.H{
		"access_addrs": adr,
	})
	return
}

func AddressDefault(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	addressId := c.Param("id")
	pkid, _ := strconv.Atoi(addressId)

	custId, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GetCustomerIdFaile)
		return
	}

	err := svc.SetDefaultAddress(int64(pkid), custId.(int64))
	if err != nil {
		response.ToErrorResponse(errcode.ServerError)
		global.Logger.Errorf(ctx, "user.address.AddressDefault err= %v", err)
		return
	}
	response.ToResponse(gin.H{
		"status": 0,
		"msg":    "设置成功",
	})
	return
}

func AddressModify(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	addressId := c.Param("id")
	pkid, _ := strconv.Atoi(addressId)

	custId, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GetCustomerIdFaile)
		return
	}
	requestBody := model.Address{} //request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		response.ToErrorResponse(errcode.RequestParamsError) // 解析错误，参数有误
		// global.Logger.Errorf(ctx, "user.address.AddressModify: %v", requestBody)
		return
	}

	requestBody.CustomerAddrId = int64(pkid)
	requestBody.CustomerId = custId.(int64)

	err := svc.SetDefaultAddress(requestBody.CustomerAddrId, requestBody.CustomerId)
	if err != nil {
		response.ToErrorResponse(errcode.AddressSetDefaultError)
		global.Logger.Errorf(ctx, "user.address.AddressModify.SetDefaultAddress,err= %v, pk= %v, customId= %v",
			err, requestBody.CustomerAddrId, requestBody.CustomerId)
		return
	}

	err = svc.ModifyAddress(requestBody)
	if err != nil {
		response.ToErrorResponse(errcode.AddressModifyRequestError)
		global.Logger.Errorf(ctx, "user.address.AddressModify, err= %v", err)
		return
	}

	response.ToResponse(gin.H{
		"status": 0,
		"msg":    "修改成功",
	})
	return

}
