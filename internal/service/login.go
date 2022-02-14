package service

import (
	"github.com/lemon997/lemonMall/internal/dao"
	"github.com/lemon997/lemonMall/internal/model"
)

//几个逻辑点：
// 1.用户名存在但是密码不正确，返回403
// 2.用户名不存在，请注册,返回404
// 3.用户名存在密码正确，返回200,statusok

type LoginRequest struct {
	RegisterRequest
}

func (svc *Service) GetCustomerInfosInDB(name string) (model.Login, error) {
	return dao.LoginMethod{}.GetRow(name)
}
