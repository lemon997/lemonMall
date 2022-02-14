package errcode

import (
	"fmt"
	"net/http"
)

//错误码处理
// 1是服务器级别，2是用户级别
// 00是模块级别，01是用户模块，03是请求表单模块，04是JWT模块,05是中间件模块，06上传模块，07获取静态资源模块,08是redis模块
// 09是mysql模块,10是订单模块

var (
	Success                   = NewError(0, "成功")         //200
	ServerError               = NewError(10000, "服务内部错误") //500
	TooManyRequests           = NewError(10300, "请求过多")   //500
	GinGetError               = NewError(10500, "gin.get找不到值")
	WriteRedisError           = NewError(10801, "redis写入失败")
	InvalidParams             = NewError(20000, "入参错误") //400
	CustomerNotFound          = NewError(20100, "找不到该用户")
	PwdError                  = NewError(20101, "密码错误")
	RequestParamsError        = NewError(20301, "用户注册表单参数错误") //400
	RequestParmsIsExistError  = NewError(20302, "注册用户名已经存在")  //400
	LogoutRequest             = NewError(20303, "退出成功")       //告诉前端删除token
	AddressModifyRequestError = NewError(20304, "地址管理信息修改失败")
	AddressSetDefaultError    = NewError(20305, "地址管理信息设置默认地址失败")
	GetCategoryError          = NewError(20306, "获取商品分类失败")
	UnauthorizedAuthNotExist  = NewError(20401, "鉴权失败") //以下401
	UnauthorizedTokenError    = NewError(20402, "鉴权失败，Token 错误")
	UnauthorizedTokenTimeout  = NewError(20403, "鉴权失败，Token 超时")
	UnauthorizedTokenGenerate = NewError(20404, "鉴权失败，Token 生成失败")
	GetCustomerIdFaile        = NewError(20501, "获取customerId失败")
	ErrorUploadFileFail       = NewError(20601, "上传文件失败")
	ErrorGetSwipe             = NewError(20701, "轮播图请求失败")
	ErrorGetRecommend         = NewError(20702, "推荐商品请求失败")
	ErrorGetGoodDate          = NewError(20703, "获取商品数据失败")
	ErrorDelFail              = NewError(20804, "删除失败")
	ErrorSelectFail           = NewError(20901, "select获取失败")
	ErrorOrder                = NewError(21001, "下单失败，请稍后尝试")
	ErrorPay                  = NewError(21002, "购买失败")
)

type Error struct {
	code    int      `json:"code"`
	msg     string   `json:"msg"`
	details []string `json:"details"`
}

var codes = map[int]string{}

//对状态码统一管理
func NewError(code int, msg string) *Error {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("错误码%d已经存在，请换一个", code))
	}
	codes[code] = msg
	return &Error{code: code, msg: msg}
}

func (e *Error) Error() string {
	return fmt.Sprintf("错误码:%d,错误信息:%s", e.Code(), e.Msg())
}

func (e *Error) Code() int {
	return e.code
}

func (e *Error) Msg() string {
	return e.msg
}

func (e *Error) Msgf(args []interface{}) string {
	return fmt.Sprintf(e.msg, args...)
}

func (e *Error) Details() []string {
	return e.details
}

func (e *Error) WithDetails(details ...string) *Error {
	newError := *e
	newError.details = []string{}
	for _, d := range details {
		newError.details = append(newError.details, d)
	}

	return &newError
}

func (e *Error) StatusCode() int {
	switch e.Code() {

	case GinGetError.Code():
		return http.StatusInternalServerError

	case LogoutRequest.Code():
		return http.StatusOK

	case Success.Code():
		return http.StatusOK

	case ServerError.Code():
		//系统错误返回500
		return http.StatusInternalServerError

	case PwdError.Code():
		//密码错误
		return http.StatusBadRequest

	case CustomerNotFound.Code():
		//用户找不到
		return http.StatusNotFound

	case WriteRedisError.Code():
		return http.StatusBadRequest

	case ErrorOrder.Code():
		return http.StatusBadRequest

	case RequestParamsError.Code():
		//注册表单参数错误返回400
		// fallthrough
		return http.StatusBadRequest

	case ErrorPay.Code():
		return http.StatusBadRequest

	case ErrorDelFail.Code():
		return http.StatusBadRequest

	case RequestParmsIsExistError.Code():
		//用户信息重复
		// fallthrough
		return http.StatusBadRequest

	case InvalidParams.Code():
		//参数问题返回400
		return http.StatusBadRequest

	case ErrorSelectFail.Code():
		return http.StatusBadRequest

	case UnauthorizedAuthNotExist.Code():
		// fallthrough
		return http.StatusUnauthorized

	case UnauthorizedTokenError.Code():
		// fallthrough
		return http.StatusUnauthorized

	case UnauthorizedTokenGenerate.Code():
		// fallthrough
		return http.StatusUnauthorized

	case GetCustomerIdFaile.Code():
		return http.StatusUnauthorized

	case UnauthorizedTokenTimeout.Code():
		//jwt验证失败统一401
		return http.StatusUnauthorized

	case TooManyRequests.Code():
		return http.StatusInternalServerError
	}
	//都不是则回复500
	return http.StatusInternalServerError
}
