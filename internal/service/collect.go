package service

import (
	"github.com/lemon997/lemonMall/internal/dao"
	"github.com/lemon997/lemonMall/internal/model"
)

func (svc Service) GetCollectWithCustomerID(productID, customerID int64) (model.Collect, error) {
	return dao.NewCollectMethod(svc.ctx).Get(productID, customerID)
}

func (svc Service) SelectCollectWithCustomerID(customerID, startIndex, endIndex int64) ([]model.Collect, error) {
	return dao.NewCollectMethod(svc.ctx).Select(customerID, startIndex, endIndex)
}

func (svc Service) DelCollectWithCustomeID(productID, customerID int64) error {
	return dao.NewCollectMethod(svc.ctx).Del(productID, customerID)
}

func (svc Service) InsertCollectWithProductID(productID, customerID int64) error {
	product, err := dao.NewProductsMethod(svc.ctx).Get(productID)
	if err != nil {
		return err
	}
	collect := model.Collect{}
	collect.CustomerID = customerID
	collect.ProductID = productID
	collect.ProductName = product.ProductName
	collect.ImgURL = product.ImgUrl
	return dao.NewCollectMethod(svc.ctx).Insert(collect)
}
