package user

import (
	"github.com/gin-gonic/gin"
	"github.com/lemon997/lemonMall/common/app"
	"github.com/lemon997/lemonMall/common/errcode"
	"github.com/lemon997/lemonMall/internal/redisDao"
	"github.com/lemon997/lemonMall/internal/service"
)

func UpdatePwd(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	name, exist := c.Get("LoginName")
	if !exist {
		response.ToErrorResponse(errcode.GinGetError)
		// global.Logger.Infof("updatePwd.exist:",)
		return
	}
	var requestBody service.LoginRequest //request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		response.ToErrorResponse(errcode.RequestParamsError) // 解析错误，参数有误
		return
	}
	svc := service.New(c.Request.Context())
	err := svc.UpdatePwd(requestBody.Password, name.(string))
	if err != nil {
		response.ToErrorResponse(errcode.ServerError)
		return
	}

	id, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GinGetError)
		return
	}
	key := svc.GetRedisTokenKey(id.(int64), name.(string))

	rdb := redisDao.NewMethodsRedisToken()
	err = rdb.DelKV(ctx, key)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError)
	}
	// response.ToResponse(gin.H{
	// 	"status": 20303,
	// 	"msg":    "退出成功",
	// })

	response.ToErrorResponse(errcode.LogoutRequest)

}
