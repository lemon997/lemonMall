package service

import (
	"github.com/lemon997/lemonMall/internal/dao"
	"github.com/lemon997/lemonMall/internal/model"
)

func (svc *Service) UpdatePwd(pwd, name string) error {
	//修改密码应该使用手机或者邮箱认证，能力有限，目前没有写这个功能
	return dao.LoginMethod{}.UpdatePwd(svc.ctx, pwd, name)
}

func (svc Service) GetNameAndUrl(customerID int64) (model.Login, error) {
	return dao.LoginMethod{}.GetNameAndUrl(customerID)
}
