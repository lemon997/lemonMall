package user

import (
	"github.com/gin-gonic/gin"
	"github.com/lemon997/lemonMall/common/app"
	"github.com/lemon997/lemonMall/common/errcode"
	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/service"
)

func Name(c *gin.Context) {
	response := app.NewResponse(c)
	ctx := c.Request.Context()
	svc := service.New(ctx)

	customerId, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GetCustomerIdFaile)
		return
	}
	data, err := svc.GetNameAndUrl(customerId.(int64))
	if err != nil {
		global.Logger.Errorf(ctx, "user.getUserName,err= %v", err)
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	response.ToResponse(gin.H{
		"data": data,
	})
}
