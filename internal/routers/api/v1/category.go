package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lemon997/lemonMall/common/app"
	"github.com/lemon997/lemonMall/common/errcode"
	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/service"
)

func Category(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)
	var pkId int64 = 0
	resCategory, err := svc.GetAllCategory(pkId)
	if err != nil {
		response.ToErrorResponse(errcode.GetCategoryError)
		global.Logger.Errorf(ctx, "v1.category.GetAllCategory, err= %v", err)
		return
	}

	res := []service.CategoryRequest{}

	for i := 0; i < len(resCategory); i++ {
		categ := service.CategoryRequest{}
		categ.CategoryID = resCategory[i].CategoryID
		categ.CategoryName = resCategory[i].CategoryName
		childrenRes, err := svc.GetAllCategoryChildren(resCategory[i].CategoryID)

		if err != nil {
			response.ToErrorResponse(errcode.GetCategoryError)
			global.Logger.Errorf(ctx, "v1.category.GetAllCategoryChildren, err= %v", err)
			return
		}
		categ.Children = childrenRes
		res = append(res, categ)
	}

	response.ToResponse(gin.H{
		"status":   0,
		"msg":      "获取分类成功",
		"category": res,
	})

}
