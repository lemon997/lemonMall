package service

import (
	"github.com/lemon997/lemonMall/internal/dao"
	"github.com/lemon997/lemonMall/internal/model"
)

func (svc Service) GetProducts(ProductId int64) (model.Products, error) {
	//获取Product_id对应的商品信息
	return dao.NewProductsMethod(svc.ctx).Get(ProductId)
}

func (svc Service) SelectProductsWithSettlementDataResponse(nums *[]model.SettlementDataRepsonse) error {
	var err error = nil
	for i := 0; i < len(*nums); i++ {
		data, err := svc.GetProducts((*nums)[i].Goods.ProductID)
		if err != nil {
			break
		}
		(*nums)[i].Goods.ImgUrl = data.ImgUrl
		(*nums)[i].Goods.Price = data.Price
		(*nums)[i].Goods.ProductName = data.ProductName
	}
	return err
}

func (svc Service) SelectProducts(startIndex, endIndex int64) ([]model.Products, error) {
	return dao.NewProductsMethod(svc.ctx).Select(0, startIndex, endIndex)
}
