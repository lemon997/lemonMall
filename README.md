#化妆品商场
##后端工具
gin, Nginx, Redis, MySQL, RabbitMQ<br>
##实现功能
首页, 分类, 购物车, 个人信息, 用户地址管理, 用户收藏, 上传头像, 登录注册, 生产订单, 伪支付, 分类页面排序, 订单详情, 商品详情, 支付失败且过期的订单实现<br>
##实现思路
JWT身份验证, 雪花算法生成订单编号, RabbitMQ实现过期订单,Redis的Zset用作商品排序<br>
##部分代码来源
煎鱼大佬写的Go语言编程之旅第二章案例<br>
go-programming-tour-book/blog-service<br>
##使用
go mod vendor
go run main.go
