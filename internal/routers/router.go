package routers

import (
	"net/http"

	"github.com/lemon997/lemonMall/global"

	"github.com/lemon997/lemonMall/internal/routers/api/v1/uploading"

	"github.com/lemon997/lemonMall/internal/routers/api/v1/user"

	"github.com/lemon997/lemonMall/internal/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/lemon997/lemonMall/docs"
	"github.com/lemon997/lemonMall/internal/routers/api/v1"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func NewRouter() *gin.Engine {
	// 初始化，运行，路由、对象、模板管理、调度等
	r := gin.New()
	//中间件使用
	//输出请求日志，异常捕获
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(Cors())

	//初始化docs包，和注册一个针对swagger的路由，swagger.json会默认指向当前应用所启动的域名下的swagger/doc.json路径
	url := ginSwagger.URL("http://127.0.0.1:9090/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	apiv1 := r.Group("/api/v1")
	{
		apiv1.POST("/login", v1.Login)
		apiv1.POST("/register", v1.Register)
		apiv1.GET("/swipe", v1.Swipe)
		apiv1.GET("/recommend", v1.Recommend)
		apiv1.GET("/category", v1.Category)
		apiv1.GET("/getgooddate", v1.GetGoodDate)
		apiv1.GET("/getgoodsdetail/:id", v1.GetGoodsDetail)
		apiv1.GET("/gettabgoodsdata/:id", v1.GetTabGoodsData)
	}

	apiv1User := r.Group("/api/v1/user")
	apiv1User.Use(middleware.JWT()) //该分组下的接口均需要校验token
	{
		apiv1User.GET("/name", user.Name)
		apiv1User.POST("/logout", user.Logout)
		apiv1User.POST("/updatepwd", user.UpdatePwd)
		apiv1User.GET("/getaddressdetail/:id", user.GetAddressDetail)
		apiv1User.GET("/getalladdress", user.GetAllAddress)
		apiv1User.POST("/addaddress", user.AddAddress)
		apiv1User.DELETE("/deladdress/:id", user.DelAddress)
		apiv1User.PATCH("/addressdefault/:id", user.AddressDefault)
		apiv1User.PUT("/addressmodify/:id", user.AddressModify)

		//收藏模块
		apiv1User.PATCH("/collects/cancelgoods/:id", user.CancelCollect)
		apiv1User.PATCH("/collects/setgoods/:id", user.SetCollect)
		apiv1User.GET("/collects", user.GetCollectList)
		apiv1User.GET("/collects/checkgoods/:id", user.CheckCollect)

		//购物车
		apiv1User.GET("/carts", user.GetCartList)
		apiv1User.POST("/carts", user.AddCart)
		apiv1User.PUT("/carts/:id", user.SetCartProductNum)
		apiv1User.POST("/carts/checked", user.ChangeCheckedStatus)
		apiv1User.DELETE("/carts/:id", user.DelCart)

		//订单模块
		apiv1User.GET("/orderlist", user.GetOrderList)
		apiv1User.GET("/orderdetail/:id", user.GetOrderDetail)
		apiv1User.POST("/order/preview", user.GetSettlementData)
		apiv1User.POST("/order/buy", user.SubmitOrder)
		apiv1User.POST("/order/payment", user.Payment)
	}
	// apiv1Adr := r.Group("/addressMotify", user.GetAddressDetail)
	// apiv1Adr.Use(middleware.JWT())

	apiv1Upload := r.Group("/api/v1/uploading")
	apiv1Upload.Use(middleware.JWT())
	{
		upload := uploading.NewUpload()
		apiv1Upload.POST("/file", upload.UploadFile)
	}
	r.StaticFS("/static", http.Dir(global.AppSetting.UploadSavePath)) //返回一个handler允许http请求访问静态资源

	return r
}

func Cors() gin.HandlerFunc {
	//跨域中间件
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin") //请求头部
		method := c.Request.Method
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE, PATCH")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers",
				"X-Requested-With,Authorization, Content-Length, X-CSRF-Token, Token,session,Accept, Origin, Host,Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With,If-Modified-Since, Cache-Control, Content-Type, Pragma")

			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusNoContent, "ok")
		}
		defer func() {
			if err := recover(); err != nil {
				global.Logger.Errorf(c.Request.Context(), "Panic info is: %v", err)
			}
		}()
		// 处理请求
		c.Next()
	}
}
