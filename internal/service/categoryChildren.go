package service

import (
	"github.com/lemon997/lemonMall/internal/dao"
	"github.com/lemon997/lemonMall/internal/model"
)

func (svc Service) GetAllCategoryChildren(categoryID int64) ([]model.CategoryChildren, error) {
	//获取categoryID对应的子分类信息
	return dao.NewCategoryChildrenMethod(svc.ctx).Select(categoryID)
}
