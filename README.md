Go语言实现化妆品商场后端
====
前端使用vue3+vue-cli4+vant3实现，地址如下
```
https://github.com/lemon997/lemonMall-front
```

项目已经上线，地址如下
```
http://120.79.132.3
```

后端工具
----
gin, Nginx, Redis, MySQL, RabbitMQ<br>

实现功能
---
首页, 分类, 购物车, 个人信息, 用户地址管理, 用户收藏, 上传头像, 登录注册, 生产订单, 伪支付, 分类页面排序, 订单详情, 商品详情, 支付失败且过期的订单实现<br>

实现思路
---
JWT身份验证, 雪花算法生成订单编号, RabbitMQ实现过期订单,Redis的Zset用作商品排序<br>


部分代码来源
---

煎鱼大佬写的Go语言编程之旅第二章案例，地址如下
```
https://github.com/go-programming-tour-book/blog-service
```

使用
---

请配置好configs目录下的配置再启动，根据实际情况更改

```
go mod vendor
```

```
go run main.go
```

目录结构
---
```
├── build
│   └── lemonMall
├── common
│   ├── app
│   │   ├── app.go
│   │   └── pagination.go
│   ├── authJWT
│   │   └── authJWT.go
│   ├── convert
│   │   └── convert.go
│   ├── errcode
│   │   └── errcode.go
│   ├── logger
│   │   └── logger.go
│   ├── setting
│   │   ├── section.go
│   │   └── setting.go
│   ├── upload
│   │   └── file.go
│   └── util
│       ├── byteSlipeAndString.go
│       ├── isSame.go
│       ├── md5.go
│       ├── priceFormatConversion.go
│       └── timeFromat.go
├── configs
│   └── config.yaml
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── global
│   ├── db.go
│   ├── mq.go
│   ├── redis.go
│   ├── setting.go
│   ├── snowflake.go
│   └── snowflake_test.go
├── go.mod
├── go.sum
├── internal
│   ├── dao
│   │   ├── address.go
│   │   ├── categoryChildren.go
│   │   ├── categoryChildren_test.go
│   │   ├── category.go
│   │   ├── category_test.go
│   │   ├── collect.go
│   │   ├── dao.go
│   │   ├── goodsResponse.go
│   │   ├── loginDB.go
│   │   ├── order.go
│   │   ├── products.go
│   │   ├── products_test.go
│   │   ├── recommend.go
│   │   └── swipe.go
│   ├── middleware
│   │   └── jwt.go
│   ├── model
│   │   ├── address.go
│   │   ├── cart.go
│   │   ├── categoryChildren.go
│   │   ├── category.go
│   │   ├── collect.go
│   │   ├── goodsResponse.go
│   │   ├── login.go
│   │   ├── model.go
│   │   ├── mq.go
│   │   ├── order.go
│   │   ├── products.go
│   │   ├── recommend.go
│   │   ├── redisKey.go
│   │   └── swipe.go
│   ├── mq
│   │   ├── client.go
│   │   ├── consumer
│   │   │   └── consumer.go
│   │   ├── delayqueue
│   │   │   └── mysqlDelay.go
│   │   ├── messagequeue
│   │   │   └── mysqlReduce.go
│   │   └── producer
│   │       └── producer.go
│   ├── queueservice
│   │   ├── mqservice.go
│   │   ├── order.go
│   │   ├── products.go
│   │   └── queueReceiver.go
│   ├── redisDao
│   │   ├── cart.go
│   │   ├── delayQueue.go
│   │   ├── errors.go
│   │   ├── order.go
│   │   ├── redisDao.go
│   │   ├── redisGoodDetail.go
│   │   ├── redisGoodList.go
│   │   ├── redisGoodList_test.go
│   │   ├── redisRecommend.go
│   │   ├── redisSwipe.go
│   │   ├── redisToken.go
│   │   └── stock.go
│   ├── routers
│   │   ├── api
│   │   │   └── v1
│   │   │       ├── category.go
│   │   │       ├── getGoodDate.go
│   │   │       ├── getGoodsDetail.go
│   │   │       ├── login.go
│   │   │       ├── recommend.go
│   │   │       ├── register.go
│   │   │       ├── swipe.go
│   │   │       ├── tabGoods.go
│   │   │       ├── uploading
│   │   │       │   └── file.go
│   │   │       └── user
│   │   │           ├── address.go
│   │   │           ├── cart.go
│   │   │           ├── collect.go
│   │   │           ├── getUserName.go
│   │   │           ├── logout.go
│   │   │           ├── order.go
│   │   │           └── updatePwd.go
│   │   └── router.go
│   └── service
│       ├── address.go
│       ├── auth.go
│       ├── cart.go
│       ├── categoryChildren.go
│       ├── categoryChildrenSortInRedis.go
│       ├── category.go
│       ├── collect.go
│       ├── errors.go
│       ├── goods.go
│       ├── login.go
│       ├── logout.go
│       ├── modifyInfos.go
│       ├── order.go
│       ├── pay.go
│       ├── products.go
│       ├── recommend.go
│       ├── register.go
│       ├── service.go
│       ├── stock.go
│       ├── swipe.go
│       └── upload.go
├── LICENSE
├── main.go
├── README.md
└── storage
    └── logs
        └── app.log
```

