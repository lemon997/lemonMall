package v1

import (
	"encoding/json"

	"github.com/lemon997/lemonMall/common/util"

	"github.com/lemon997/lemonMall/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/lemon997/lemonMall/common/app"
	"github.com/lemon997/lemonMall/common/convert"
	"github.com/lemon997/lemonMall/common/errcode"
	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/service"
)

func GetGoodDate(c *gin.Context) {
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	sortType := c.Query("type")
	// tmp3 := c.Query("title")
	categoryChildrenID, err := convert.StrTo(c.Query("category_children_id")).MustInt64()
	if err != nil {
		global.Logger.Errorf(ctx, "v1.getGoodDate.StrTo,err= %v, category_children_id=%v", err, c.Query("category_children_id"))
		response.ToErrorResponse(errcode.InvalidParams)
		return
	}
	if categoryChildrenID == 0 {
		categoryChildrenID = 1
	}

	//生成用于Redis的key
	key := svc.GenerateGoodsDateKeyUseForRedis(sortType, categoryChildrenID)
	if num := svc.CheckCategorySortKeyExists(key); num == 0 {
		err := svc.AddElementsInRedisAndSortByKey(key, categoryChildrenID)
		if err != nil {
			global.Logger.Errorf(ctx, "v1.getGoodDate.AddElements, err= %v", err)
			response.ToErrorResponse(errcode.ServerError)
			return
		}
		n := svc.CheckCategorySortKeyExists(key)
		if n == 0 {
			response.ToErrorResponse(errcode.ErrorGetGoodDate)
			return
		}
	}

	//根据页码转换成区间下标，页码1对应0-9,页码2对应10-19
	page := app.GetPage(c)
	startIndex, endIndex := app.GetPageInterval(page)

	//获取排序结果
	res, err := svc.SortByKey(key, sortType, startIndex, endIndex)
	if err != nil {
		global.Logger.Errorf(ctx, "v1.getGoodDate.SortByKey, err= %v", err)
		response.ToErrorResponse(errcode.ErrorGetGoodDate)
		return
	}

	length := len(res)

	//Zset的member是已经序列化的，懒得前端再反序列化，将res字符串切片转成byte二维切片，使用指针强转类型方式
	b := make([][]byte, length)
	for i := 0; i < length; i++ {
		b[i] = util.S2B(res[i])
	}

	//反序列化结果放到Products切片中
	var goods []model.Products
	for i := 0; i < length; i++ {
		tmp := model.Products{}
		err := json.Unmarshal(b[i], &tmp)
		if err != nil {
			global.Logger.Errorf(ctx, "v1.getGoodDate.json.Unmarshal, err= %v", err)
			response.ToErrorResponse(errcode.ServerError)
			return
		}
		goods = append(goods, tmp)
	}

	response.ToResponse(gin.H{
		"status": 0,
		"msg":    "商品获取成功",
		"goods":  goods,
	})

}
