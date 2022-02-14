package service

import (
	"strconv"
	"strings"

	"github.com/lemon997/lemonMall/global"

	"github.com/lemon997/lemonMall/internal/dao"
	"github.com/lemon997/lemonMall/internal/model"
	"github.com/lemon997/lemonMall/internal/redisDao"
)

//添加到购物车操作只记录商品id以及商品数量
func (svc Service) AddCartProduct(customerID, productID, num int64) error {
	key := model.CartKey(customerID)
	return redisDao.NewMethodsCartProductNum(svc.ctx).AddProductAtCart(key, productID, num)
}

func (svc Service) DelCartProduct(customerID, productID int64) error {
	key := model.CartKey(customerID)
	field := model.ProductIDField(productID)
	return redisDao.NewMethodsCartProductNum(svc.ctx).DelProductAtCart(key, field)
}

func (svc Service) DelCartField(key, field string) error {
	return redisDao.NewMethodsCartProductNum(svc.ctx).DelProductAtCart(key, field)
}

func (svc Service) GetCartProductList(customerID int64) ([]model.CartResponse, error) {
	key := model.CartKey(customerID)
	nums, err := redisDao.NewMethodsCartProductNum(svc.ctx).GetProductListAtCart(key)
	if err != nil {
		return []model.CartResponse{}, err
	}
	// tmpHash是保存选择状态的商品
	tmpHash := make(map[int64]struct{}, 8)
	// resp是返回的结果切片
	resp := make([]model.CartResponse, 0, len(nums))
	//遍历map,边查询DB
	for i, v := range nums {
		// 判断前缀是否是选择状态
		isPrefix := strings.HasPrefix(i, model.CheckedStatusFieldPrefix)
		if isPrefix {
			pID := strings.TrimPrefix(i, model.CheckedStatusFieldPrefix)
			n, _ := strconv.ParseInt(pID, 10, 64)
			tmpHash[n] = struct{}{}
			continue
		}

		productID, _ := strconv.ParseInt(i, 10, 64)
		num, _ := strconv.ParseInt(v, 10, 64)

		productInfo, err := dao.NewProductsMethod(svc.ctx).Get(productID)
		if err != nil {
			global.Logger.Infof(svc.ctx, "dao.NewProductsMethod(),err= %v, productID= %v", err, productID)
			continue
		}
		tmp := model.CartResponse{}
		tmp.ProductName = productInfo.ProductName
		tmp.ProductID = productID
		tmp.ImgURL = productInfo.ImgUrl
		tmp.Price = productInfo.Price
		tmp.Num = num
		resp = append(resp, tmp)
	}

	for i := 0; i < len(resp); i++ {
		if _, ok := tmpHash[resp[i].ProductID]; ok {
			resp[i].IsChecked = true
		}
	}
	return resp, nil
}

func (svc Service) SetCartProductNum(customerID, productID, num int64) error {
	key := model.CartKey(customerID)
	return redisDao.NewMethodsCartProductNum(svc.ctx).SetProductNumAtCart(key, productID, num)
}

func (svc Service) ChangeCheckedStatus(customerID int64, cartIds []int64) error {
	//redis不存在则默认是没有选择，false,有选择则true
	// 选择状态，cartIds的值是选择后的商品ID,前端传来的是一个数组，如果没有选择，则一定是空数组，因此采用删除再添加策略
	key := model.CartKey(customerID)
	methods := redisDao.NewMethodsCartProductNum(svc.ctx)
	var value bool = true
	res, _ := methods.GetProductListAtCart(key)
	for i := range res {
		if exist := strings.HasPrefix(i, "isCheck:"); exist {
			methods.DelProductAtCart(key, i)
		}
	}

	for i := 0; i < len(cartIds); i++ {
		field := model.CheckedStatusField(cartIds[i])
		err := methods.ChangeCheckedStatus(key, field, value)
		if err != nil {
			global.Logger.Infof(svc.ctx, "methods.ChangeCheckStatus,err= %v, key=%v, field= %v",
				err, key, field)
		}
	}
	return nil
}

func (svc Service) GetProductNum(customerID int64, productID int64) (int64, error) {
	key := model.CartKey(customerID)
	return redisDao.NewMethodsCartProductNum(svc.ctx).GetProductNumAtCart(key, productID)
}

func (svc Service) SelectProductNums(customerID int64, productIDs []int64) (*[]model.SettlementDataRepsonse, error) {
	//因为map和结构体切片的创建会逃逸到堆上,还是声明结构体切片较好，使用map还需要构建切片，因此直接使用切片
	res := make([]model.SettlementDataRepsonse, len(productIDs))
	for i := 0; i < len(productIDs); i++ {
		num, err := svc.GetProductNum(customerID, productIDs[i])
		if err != nil {
			return nil, err
		}
		if num < 0 {
			return nil, NewErrService(ErrRedisGet, "svc.SelectProductNums.GetProductNum.", err)
		}
		res[i].Num = num
		res[i].Goods.ProductID = productIDs[i]
	}
	return &res, nil
}

func (svc Service) ExistFieldWithCart(key, field string) (bool, error) {
	red := redisDao.NewMethodsCartProductNum(svc.ctx)
	return red.ExistFieldAtCart(key, field)
}

func (svc Service) ExistFieldsWithCart(key string, fields []string) (bool, error) {
	return redisDao.NewMethodsCartProductNum(svc.ctx).ExistFieldsAtCart(key, fields)
}
