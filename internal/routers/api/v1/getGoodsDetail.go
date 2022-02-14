package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lemon997/lemonMall/common/app"
	"github.com/lemon997/lemonMall/common/convert"
	"github.com/lemon997/lemonMall/common/errcode"
	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/service"
)

func GetGoodsDetail(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)
	productId, err := convert.StrTo(c.Param("id")).Int()
	if err != nil {
		global.Logger.Errorf(ctx, "v1.getGoodDetail.StrTo,err= %v, product_id=%v", err, c.Param("id"))
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}

	data, err := svc.GetGoodDatails(productId)
	if err != nil {
		response.ToErrorResponse(errcode.ServerError)
	}

	response.ToResponse(gin.H{
		"status": 0,
		"msg":    "获取详情成功",
		"goods":  data,
	})

}
