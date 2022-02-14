package dao

import (
	"context"

	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/model"
)

type SwipeMethod struct {
	ctx context.Context
}

func NewSwipeMethod(c context.Context) SwipeMethod {
	return SwipeMethod{ctx: c}
}

func (s SwipeMethod) Insert(url string) error {
	sqlStr := `insert into swipe_imgs_url (img_url) values(?)`
	res, err := global.DBEngineShop.Exec(sqlStr, url)
	num, _ := res.LastInsertId()
	if num > 0 || err == nil {
		return nil
	}
	return err
}

func (s SwipeMethod) Get(id int64) (model.SwipeImg, error) {
	//查询单个id
	sqlStr := `select img_id , img_url from swipe_imgs_url where img_id = ?`
	img := model.SwipeImg{}
	if err := global.DBEngineShop.Get(&img, sqlStr, id); err != nil {
		return model.SwipeImg{}, err
	}
	return img, nil
}

func (s SwipeMethod) Select(id int64) ([]model.SwipeImg, error) {
	//查询id后所有轮播图的URL
	sqlStr := `select img_id , img_url from swipe_imgs_url where img_id > ?`
	imgs := []model.SwipeImg{}
	if err := global.DBEngineShop.Select(&imgs, sqlStr, id); err != nil {
		return imgs, err
	}
	return imgs, nil
}

func (s SwipeMethod) Del(id int64) (err error) {
	//返回值：error,如果发生panic, 会捕捉这个错误，同时回滚，如果是事务错误，则会返回回滚的err或者提交的err
	//入参：id是swipe_imgs_url的自增主键id值,根据id进行查询和删除
	//描述：先查询目标id是否存在，存在则删除，需要开启事务
	tx, err := global.DBEngineShop.Beginx()
	if err != nil {
		return
	}
	//defer的err针对两句DB执行语句
	defer func() {
		if p := recover(); p != nil {
			global.Logger.Errorf(s.ctx, "dao.swipeTable.Del, panic: %v", p)
			err = tx.Rollback()
			return
		} else if err != nil {
			global.Logger.Infof(s.ctx, "dao.swipeTable.Del, tx, err: %v", err)
			err = tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}()

	sqlStr1 := `SELECT img_id, img_url FROM swipe_imgs_url WHERE img_id = ?`
	sqlStr2 := `DELETE FROM swipe_imgs_url WHERE img_id = ? AND img_url = ?`

	swipe := model.SwipeImg{}
	if err = tx.Get(&swipe, sqlStr1, id); err != nil {
		return
	}

	if _, err = tx.Exec(sqlStr2, id, swipe.ImgUrl); err != nil {
		return
	}
	return
}

func (s SwipeMethod) Modify(swipe model.SwipeImg, newUrl string) (err error) {
	//返回值：error,如果发生panic, 会捕捉这个错误，同时回滚，如果是事务错误，则会返回回滚的err或者提交的err
	//入参：id是swipe_imgs_url的自增主键id值, 根据id进行查询和修改， newUrl是准备替换后的值
	//描述：先查询目标id是否存在，存在则使用乐观锁，对该id的img_url进行修改，需要开启事务
	tx, err := global.DBEngineShop.Beginx()
	if err != nil {
		return
	}
	//defer的err针对两句DB执行语句
	defer func() {
		if p := recover(); p != nil {
			global.Logger.Errorf(s.ctx, "dao.swipeTable.Modify, panic: %v", p)
			err = tx.Rollback()
			return
		} else if err != nil {
			global.Logger.Infof(s.ctx, "dao.swipeTable.Modify, err: %v", err)
			err = tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}()

	// sqlStr1 := `SELECT img_id FROM swipe_imgs_url WHERE img_id = ?`
	sqlStr2 := `UPDATA swipe_imgs_url SET img_url = ? WHERE img_id = ? AND img_url = ?`

	// swipe := model.SwipeImg{}
	// if err = tx.Get(&swipe, sqlStr1, id); err != nil {
	// 	return
	// }

	if _, err = tx.Exec(sqlStr2, newUrl, swipe.ImgId, swipe.ImgUrl); err != nil {
		return
	}
	return
}
