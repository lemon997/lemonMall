package service

import (
	"sync"

	"github.com/lemon997/lemonMall/internal/dao"
)

//思路，先写map表，防止同时注册同名的错误步骤，写完map表再查询数据库
type checkNameInMap struct {
	m sync.Map //支持并发读写的map
}

var cnOnce sync.Once
var queryMap *checkNameInMap

//避免重复初始化map
func newMap() {
	cnOnce.Do(func() {
		queryMap = new(checkNameInMap)
	})
}

func (c *checkNameInMap) storeMapKey(name string) bool {
	//ok==true表示map中有重复name，返回true
	//再次检查，避免上锁之前有其他请求写入hash表
	if _, ok := c.m.Load(name); ok {
		return true
	} else {
		c.m.Store(name, struct{}{})
	}
	return false
}

func (c *checkNameInMap) deleteMapKey(name string) {
	c.m.Delete(name)
}

type RegisterRequest struct {
	LoginName string `db:"login_name" json:"loginname" binding:"min=1,max=6"`
	Password  string `db:"password" json:"password" binding:"min=1,max=6"`
}

func (svc *Service) CheckUserExistMap(params RegisterRequest) bool {
	newMap()
	//检查map是否存在登录名
	if isExist := queryMap.storeMapKey(params.LoginName); isExist {
		return true
	}
	return false
	//false表示登录名刚加入map,之后需要检查数据库
}

func (svc *Service) DelUserMap(params RegisterRequest) {
	newMap()
	queryMap.deleteMapKey(params.LoginName)
}

func (svc *Service) CheckUserExistDB(name string) bool {
	user := dao.LoginMethod{}
	if _, err := user.GetRow(name); err == nil {
		//查询数据库是否有同名用户，查询到则err==nil
		return true
	}
	return false
}

func (svc *Service) UserStoreDB(params RegisterRequest) error {
	err := dao.LoginMethod{}.StoreDB(params.LoginName, params.Password)
	return err
}
