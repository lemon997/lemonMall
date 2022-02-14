package service

import (
	"strconv"
	// "github.com/lemon997/lemonMall/internal/redisDao"

	"github.com/lemon997/lemonMall/common/authJWT"
)

//声明auth接口入参校验
type AuthRequest struct {
	// //playload存储ID,username,version,version用于区别是否有更改密码
	authJWT.Claims
}

func (svc *Service) GetRedisTokenKey(id int64, name string) string {
	return strconv.Itoa(int(id)) + ":" + name
}

func (svc *Service) CheckJWT(tokenString string) (bool, error) {
	//JWT解析错误原因：
	// 密钥错误
	// 有效时间失效
	claims, err := authJWT.ParseToken(tokenString)
	if err != nil || claims == nil {
		//密钥错误或者过期
		return false, err
	}
	// //jwt有效，查询redis中是否存在该token
	// key := svc.GetRedisTokenKey(claims.CustomerId, claims.LoginName)
	// rdb := redisDao.NewMethodsRedisToken()
	// redisValue, err := rdb.QueryGetKV(key)
	// if err != nil {
	// 	//Redis查询不到则报错
	// 	return false, err
	// }
	return true, nil
}
