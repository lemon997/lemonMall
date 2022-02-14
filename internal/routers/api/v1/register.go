package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lemon997/lemonMall/common/app"
	"github.com/lemon997/lemonMall/common/errcode"
	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/service"
)

func Register(c *gin.Context) {
	ctx := c.Request.Context()
	params := service.RegisterRequest{}
	response := app.NewResponse(c) //声明返回处理变量
	svc := service.New(ctx)        //声明顶级上下文

	if err := c.ShouldBindJSON(&params); err != nil {
		//根据自定义错误码，使用ToErrorResponse返回gin.H{}信息
		response.ToErrorResponse(errcode.RequestParamsError)
		global.Logger.Infof(ctx, "注册请求参数解析出错err: %v", err)
		return
	}
	//检查Map
	if isExist := svc.CheckUserExistMap(params); isExist {
		//已经存在
		response.ToErrorResponse(errcode.RequestParmsIsExistError)
		return
	}
	//检查数据库
	if isExist := svc.CheckUserExistDB(params.LoginName); isExist {
		//存在则返回错误
		svc.DelUserMap(params)
		response.ToErrorResponse(errcode.RequestParmsIsExistError)
		return
	}
	//数据库不存在用户名，插入数据库
	if err := svc.UserStoreDB(params); err != nil {
		global.Logger.Infof(ctx, "注册请求参数解析出错err: %v", err)
		response.ToErrorResponse(errcode.ServerError)
		svc.DelUserMap(params)

		return
	}
	svc.DelUserMap(params)
	response.ToResponse(gin.H{
		"status": 200,
		"msg":    "注册成功",
	})

}
