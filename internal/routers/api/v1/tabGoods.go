package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lemon997/lemonMall/common/app"
	"github.com/lemon997/lemonMall/common/convert"
	"github.com/lemon997/lemonMall/common/errcode"
	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/service"
)

func GetTabGoodsData(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	id, _ := convert.StrTo(c.Param("id")).Int()
	if id <= 0 {
		id = 1
	}
	startIndex, endIndex := app.GetPageInterval(id)

	res, err := svc.SelectProducts(startIndex, endIndex)
	if err != nil {
		global.Logger.Errorf(ctx, "v1.tabGoods.SelectProducts,err= %v", err)
		response.ToErrorResponse(errcode.ErrorGetGoodDate)
		return
	}
	response.ToResponse(gin.H{
		"status": 0,
		"msg":    "获取成功",
		"tab":    res,
	})
	return
}
