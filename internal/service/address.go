package service

import (
	"github.com/lemon997/lemonMall/internal/dao"
	"github.com/lemon997/lemonMall/internal/model"
)

func (svc Service) GetAddress(primaryKeyId, id int64) (model.Address, error) {
	//根据自增主键ID查询地址管理信息
	return dao.NewAddressMethod(svc.ctx).Get(primaryKeyId, id)
}

func (svc Service) SelectAddress(customerId int64) ([]model.Address, error) {
	//用户的jwt包含的customerId即可
	return dao.NewAddressMethod(svc.ctx).Select(customerId)
}

func (svc Service) DelAddress(primaryKeyId int64) error {
	return dao.NewAddressMethod(svc.ctx).Del(primaryKeyId)
}

func (svc Service) ModifyAddress(addr model.Address) error {
	//需要有自增主键ID
	return dao.NewAddressMethod(svc.ctx).Modify(addr)
}

func (svc Service) InsertAddress(addr model.Address) error {
	return dao.NewAddressMethod(svc.ctx).Insert(addr)
}

func (svc Service) SetDefaultAddress(pkId, custId int64) error {
	return dao.NewAddressMethod(svc.ctx).SetDefault(pkId, custId)
}
