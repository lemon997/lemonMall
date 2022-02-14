package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lemon997/lemonMall/common/app"
	"github.com/lemon997/lemonMall/common/authJWT"
	"github.com/lemon997/lemonMall/common/errcode"
	"github.com/lemon997/lemonMall/common/util"
	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/redisDao"
	"github.com/lemon997/lemonMall/internal/service"
)

//1、向数据库查询信息
// 2、数据库查询失败
// 3、数据库查询成功
// 4、JWT验证
// 5、JWT过期处理
// 6、JWT生成
// 7、JWT解析

// @Summary 实现登录功能
// @Produce json
// @Param loginname body string true "登录名"
// @Param password body string true "密码" char(20)
// @Param jwt body string false "JWT"
// @Success 200 {object} model.Login
// @Failure 400 {object} errcode.Error "请求错误"
// @Failure 401 {object} errcode.Error"JWT验证失败"
// @Router /api/v1/Login [post]
func Login(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)

	var requestBody service.LoginRequest //request的body中的数据，自动按照json格式解析到结构体
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		response.ToErrorResponse(errcode.RequestParamsError) // 解析错误，参数有误
		return
	}

	// svc := service.New(c.Request.Context())
	svc := service.New(ctx)
	customerInfos, err := svc.GetCustomerInfosInDB(requestBody.LoginName) //检查DB是否存在该用户,不存在则返回错误
	if err != nil {
		response.ToErrorResponse(errcode.CustomerNotFound)
		return
	}

	//校验用户密码是否正确
	if same := util.IsSame(requestBody.Password, customerInfos.Password); !same {
		response.ToErrorResponse(errcode.PwdError)
		// c.JSON(http.StatusBadRequest, gin.H{
		// 	"status": 20101,
		// 	"msg":    "用户输入密码错误",
		// })
		// return
	}

	tokenString, err := authJWT.GenerateToken(customerInfos.LoginName, customerInfos.CustomerId, customerInfos.Version)
	if err != nil { //生成token失败，记录日志
		response.ToErrorResponse(errcode.UnauthorizedTokenGenerate)
		global.Logger.Infof(ctx, "api.v1.login生成token失败，name=%s, id=%d, version=%d, err: ",
			customerInfos.LoginName, customerInfos.CustomerId, customerInfos.Version, err)
		return
	}

	key := svc.GetRedisTokenKey(customerInfos.CustomerId, customerInfos.LoginName) //生成key

	rdbMethod := redisDao.NewMethodsRedisToken() //以MySQL数据库生成的token为准，无论有无，使用setex覆盖或生成
	err = rdbMethod.ModifySetKV(ctx, key, tokenString, global.JWTSetting.Expire)
	if err != nil { //写入redis失败，日志记录
		response.ToErrorResponse(errcode.ServerError)
		return
	}
	// data := gin.H{
	// 	"status": 0,
	// 	"token":  tokenString,
	// 	"msg":    "登录成功",
	// }
	// c.JSON(http.StatusOK, data)

	response.ToResponse(gin.H{
		"status":       0,
		"access_token": tokenString,
		"msg":          "登录成功",
	}) //登录成功，返回token
	return
}
