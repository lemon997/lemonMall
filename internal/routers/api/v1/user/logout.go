package user

import (
	"github.com/gin-gonic/gin"
	"github.com/lemon997/lemonMall/common/app"

	"github.com/lemon997/lemonMall/common/errcode"

	"github.com/lemon997/lemonMall/internal/redisDao"
	"github.com/lemon997/lemonMall/internal/service"
)

func Logout(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)

	id, exist := c.Get("CustomerId")
	if !exist {
		response.ToErrorResponse(errcode.GinGetError)
		return
	}

	name, exist := c.Get("LoginName")
	if !exist {
		response.ToErrorResponse(errcode.GinGetError)
		return
	}

	svc := service.New(c.Request.Context())
	key := svc.GetRedisTokenKey(id.(int64), name.(string))

	rdb := redisDao.NewMethodsRedisToken()
	err := rdb.DelKV(ctx, key)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError)
	}
	response.ToResponse(gin.H{
		"status": 0,
		"msg":    "退出成功",
	})
}
