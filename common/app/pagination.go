package app

import (
	"github.com/lemon997/lemonMall/global"

	"github.com/gin-gonic/gin"
	"github.com/lemon997/lemonMall/common/convert"
)

//分页处理
func GetPageInterval(page int) (int64, int64) {
	var base int = 10
	var startIndex int64 = int64(base*page - 10)
	var endIndex int64 = int64(base*page - 1)
	return startIndex, endIndex
}
func GetPage(c *gin.Context) int {
	page := convert.StrTo(c.Query("page")).MustInt()
	if page <= 0 {
		return 1
	}

	return page
}

func GetPageSize(c *gin.Context) int {
	pageSize := convert.StrTo(c.Query("page_size")).MustInt()
	if pageSize <= 0 {
		return global.AppSetting.DefaultPageSize
	}
	if pageSize > global.AppSetting.MaxPageSize {
		return global.AppSetting.MaxPageSize
	}

	return pageSize
}

func GetPageOffset(page, pageSize int) int {
	result := 0
	if page > 0 {
		result = (page - 1) * pageSize
	}

	return result
}
