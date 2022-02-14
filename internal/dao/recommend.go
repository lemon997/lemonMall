package dao

import (
	"context"

	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/model"
)

type RecommendMethod struct {
	ctx context.Context
}

func NewRecommendMethod(c context.Context) RecommendMethod {
	return RecommendMethod{ctx: c}
}

func (r RecommendMethod) Insert(url string) error {
	sqlStr := `insert into recommend_imgs_url (img_url) values(?)`
	res, err := global.DBEngineShop.Exec(sqlStr, url)
	num, _ := res.LastInsertId()
	if num > 0 || err == nil {
		return nil
	}
	return err
}

func (r RecommendMethod) Get(id int64) (model.RecommendImg, error) {
	//查询单个id
	sqlStr := `select img_id, img_url from recommend_imgs_url where img_id = ?`
	img := model.RecommendImg{}
	if err := global.DBEngineShop.Get(&img, sqlStr, id); err != nil {
		return model.RecommendImg{}, err
	}
	return img, nil
}

func (r RecommendMethod) Select(id int64) ([]model.RecommendImg, error) {
	//查询id后所有轮播图的URL
	sqlStr := `select img_id, img_url from recommend_imgs_url where img_id > ?`
	imgs := []model.RecommendImg{}
	if err := global.DBEngineShop.Select(&imgs, sqlStr, id); err != nil {
		return imgs, err
	}
	return imgs, nil
}

func (r RecommendMethod) Del(id int64) (err error) {
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
			global.Logger.Errorf(r.ctx, "dao.recommendTable.Del, panic: %v", p)
			err = tx.Rollback()
			return
		} else if err != nil {
			global.Logger.Infof(r.ctx, "dao.recommendTable.Del, tx, err: %v", err)
			err = tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}()

	sqlStr1 := `SELECT img_id, img_url FROM recommend_imgs_url WHERE img_id = ?`
	sqlStr2 := `DELETE FROM recommend_imgs_url WHERE img_id = ? AND img_url = ?`

	recommend := model.RecommendImg{}
	if err = tx.Get(&recommend, sqlStr1, id); err != nil {
		return
	}

	if _, err = tx.Exec(sqlStr2, id, recommend.ImgUrl); err != nil {
		return
	}
	return
}

func (r RecommendMethod) Modify(recommend model.RecommendImg, newUrl string) (err error) {

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
			global.Logger.Errorf(r.ctx, "dao.recommendTable.Modify, panic: %v", p)
			err = tx.Rollback()
			return
		} else if err != nil {
			global.Logger.Infof(r.ctx, "dao.recommendTable.Modify, err: %v", err)
			err = tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}()

	sqlStr2 := `UPDATA recommend_imgs_url SET img_url = ? WHERE img_id = ? AND img_url = ?`

	if _, err = tx.Exec(sqlStr2, newUrl, recommend.ImgId, recommend.ImgUrl); err != nil {
		return
	}
	return
}
