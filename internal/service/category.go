package service

import (
	"github.com/lemon997/lemonMall/internal/dao"
	"github.com/lemon997/lemonMall/internal/model"
)

type CategoryRequest struct {
	CategoryID   int64                    `json:"category_id"`
	CategoryName string                   `json:"category_name"`
	Children     []model.CategoryChildren `json:"children"`
}

func (svc Service) GetAllCategory(primaryKeyId int64) ([]model.Category, error) {
	//获取总分类信息
	return dao.NewCategoryMethod(svc.ctx).Select(primaryKeyId)
}
