package middleware

import (
	"strings"

	"github.com/lemon997/lemonMall/internal/redisDao"
	"github.com/lemon997/lemonMall/internal/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lemon997/lemonMall/common/app"
	"github.com/lemon997/lemonMall/common/authJWT"
	"github.com/lemon997/lemonMall/common/errcode"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		response := app.NewResponse(c) //传入gin.context
		var tokenString string

		auth := c.GetHeader("Authorization") //获取头域Authorization的token
		if len(auth) == 0 {
			c.Abort()
			response.ToErrorResponse(errcode.UnauthorizedTokenError)
			return
		}

		tokenString = strings.Fields(auth)[1]
		if tokenString == "" {
			c.Abort()
			response.ToErrorResponse(errcode.UnauthorizedTokenError)
			return
		}

		claimsRequest, err := authJWT.ParseToken(tokenString) //以上三点均校验token是否有效

		if err != nil { //根据jwt规定的错误放回错误信息
			ecode := errcode.Success

			switch err.(*jwt.ValidationError).Errors {

			case jwt.ValidationErrorExpired:
				ecode = errcode.UnauthorizedTokenTimeout

			default:
				ecode = errcode.UnauthorizedTokenError
			}

			response.ToErrorResponse(ecode)
			c.Abort()
			return
		}

		//查询Redis是否存在该token
		rdb := redisDao.NewMethodsRedisToken()
		svc := service.New(c.Request.Context())

		key := svc.GetRedisTokenKey(claimsRequest.CustomerId, claimsRequest.LoginName)
		redisToken, err := rdb.QueryGetKV(ctx, key)

		if err != nil { //Redis查询不到，说明该token错误
			response.ToErrorResponse(errcode.UnauthorizedTokenError)
			c.Abort()
			return
		}

		claimsRedis, err := authJWT.ParseToken(redisToken)
		if err != nil { //解析redis的token失败
			ecode := errcode.Success

			switch err.(*jwt.ValidationError).Errors {

			case jwt.ValidationErrorExpired:
				ecode = errcode.UnauthorizedTokenTimeout

			default:
				ecode = errcode.UnauthorizedTokenError
			}

			response.ToErrorResponse(ecode)
			c.Abort()
			return
		}

		if claimsRequest.Version != claimsRedis.Version { //过期token
			response.ToErrorResponse(errcode.UnauthorizedTokenTimeout)
			c.Abort()
			return
		}
		//到这一步则说明Redis的token有效
		// response.ToResponse(gin.H{
		// 	"status": 0,
		// 	"msg":    "token有效",
		// })

		c.Set("CustomerId", claimsRequest.CustomerId)
		c.Set("LoginName", claimsRequest.LoginName)
		c.Next()
	}
}
