package dao

import (
	"context"

	"github.com/lemon997/lemonMall/global"

	"github.com/lemon997/lemonMall/internal/model"
)

type GoodsResponseMethod struct {
	ctx context.Context
}

func NewGoodsResponseMethodMethod(c context.Context) GoodsResponseMethod {
	return GoodsResponseMethod{ctx: c}
}

func (g GoodsResponseMethod) Get(productId int64) (model.GoodDetailResponse, error) {
	sqlStr := `select product_id, product_name, price, img_url, description, details from products where product_id = ? `
	categ := model.GoodDetailResponse{}
	err := global.DBEngineShop.Get(&categ, sqlStr, productId)
	return categ, err
}
