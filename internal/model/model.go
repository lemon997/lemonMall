package model

import (
	"fmt"

	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lemon997/lemonMall/common/setting"
)

type Model struct{}

func NewDBEngine(databaseSetting *setting.DatabaseSettingS) (*sqlx.DB, error) {
	db, err := sqlx.Connect(databaseSetting.DBType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(databaseSetting.MaxOpenConns)
	db.SetMaxIdleConns(databaseSetting.MaxIdleConns)
	db.SetConnMaxLifetime(time.Minute * 5)

	return db, nil
}

func NewDBEngineShop(databaseSetting *setting.DatabaseSettingS) (*sqlx.DB, error) {
	db, err := sqlx.Connect(databaseSetting.DBType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=Local",
		databaseSetting.UserName,
		databaseSetting.Password,
		databaseSetting.Host,
		databaseSetting.DBName2,
		databaseSetting.Charset,
		databaseSetting.ParseTime,
	))
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(databaseSetting.MaxOpenConns)
	db.SetMaxIdleConns(databaseSetting.MaxIdleConns)
	db.SetConnMaxLifetime(time.Minute * 5)

	return db, nil
}
