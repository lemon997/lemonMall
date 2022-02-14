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
func Swipe(c *gin.Context) {
	key := "swipe"
	ctx := c.Request.Context()
	response := app.NewResponse(c)
	svc := service.New(ctx)

	imgs, err := svc.GetSwipeUrlinRedis(key)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorGetSwipe)
		global.Logger.Errorf(ctx, "GetSwipeURLInRedis, err: %v", err)
	}

	if len(imgs) > 0 {
		response.ToResponse(gin.H{
			"access_imgs_url": imgs,
		})
		return
	}

	//redis获取不到则从DB获取，然后更新redis
	swipe, err := svc.GetSwipeUrlInDB(0)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorGetSwipe)
		global.Logger.Errorf(ctx, "GetSwipeURLInDB, err: %v", err)
		return
	}

	//将值写入redis
	err = svc.ModifyMultipleSwipeInRedis(key, swipe)
	if err != nil {
		global.Logger.Errorf(ctx, "v1.swipe.ModifyMultipleSwipeInRedis, err: %v", err)
	}

	//返回urls给客户端
	urls := make([]string, len(swipe))
	for i := 0; i < len(swipe); i++ {
		urls[i] = swipe[i].ImgUrl

	}

	response.ToResponse(gin.H{
		"access_imgs_url": urls,
	})
	return

}
