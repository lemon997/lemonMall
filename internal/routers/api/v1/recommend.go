package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/lemon997/lemonMall/common/app"
	"github.com/lemon997/lemonMall/common/errcode"
	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/service"
)

//读缓存，缓存不设置过期，因为轮播图是热点数据，时常访问，
// 更新轮播图操作：先修改数据库,然后删除缓存数据，修改数据库操作应该在管理系统操作
// 读操作：从缓存读取数据，读取到就返回，缓存读取不到，就从DB读取数据返回，把数据放到缓存
// 避免删除缓存后有大量请求涌向DB,采用singlefly套路，也就是合并请求，统一回复,只不过这个是用在多个协程上的。
func Recommend(c *gin.Context) {
	key := "recommend"
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)
	imgs, err := svc.GetRecommendUrlinRedis(key)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorGetRecommend)
		global.Logger.Errorf(ctx, "GetRecommendURLInRedis, err: %v", err)
	}

	if len(imgs) > 0 {
		response.ToResponse(gin.H{
			"access_imgs_url": imgs,
		})
		return
	}

	//redis获取不到则从DB获取，然后更新redis
	recommend, err := svc.GetRecommendUrlInDB(0)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorGetRecommend)
		global.Logger.Errorf(ctx, "GetRecommendUrlInDB, err: %v", err)
		return
	}

	//将值写入redis
	err = svc.ModifyMultipleRecommendInRedis(key, recommend)
	if err != nil {
		global.Logger.Errorf(ctx, "v1.recommend.ModifyMultipleRecommendInRedis, err: %v", err)
	}

	//返回urls给客户端
	urls := make([]string, len(recommend))
	for i := 0; i < len(recommend); i++ {
		urls[i] = recommend[i].ImgUrl

	}

	response.ToResponse(gin.H{
		"access_imgs_url": urls,
	})
	return

}
