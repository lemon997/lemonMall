package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/lemon997/lemonMall/internal/mq"

	_ "github.com/lemon997/lemonMall/internal/redisDao"

	"github.com/lemon997/lemonMall/internal/service"

	"github.com/gin-gonic/gin"

	"github.com/lemon997/lemonMall/common/logger"
	"github.com/lemon997/lemonMall/common/setting"
	"github.com/lemon997/lemonMall/global"
	"github.com/lemon997/lemonMall/internal/model"
	"github.com/lemon997/lemonMall/internal/routers"
	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {

	if err := setupSetting(); err != nil {
		global.Logger.Errorf(context.TODO(), "init.setupSetting err: %v", err)
	}

	if err := setupDBEngine(); err != nil {
		global.Logger.Errorf(context.TODO(), "init.setupDBEngine err: %v", err)
	}

	if err := setupRedis(); err != nil {
		global.Logger.Errorf(context.TODO(), "init.setupRedis err: %v", err)
	}

	if err := setupLogger(); err != nil {
		global.Logger.Errorf(context.TODO(), "init.setupLogger err: %v", err)

	}

	if err := snowflake(context.Background()); err != nil {
		log.Fatalf("init.snowflake err:%v", err)
	}

	if err := SetStock(); err != nil {
		global.Logger.Errorf(context.TODO(), "init.SetStock err: %v", err)
	}

	if err := InitQueue(); err != nil {
		global.Logger.Errorf(context.TODO(), "init.InitQueue err:%v", err)
	}
}

// @title 商城系统
// @version 1.0
// @description 耗费时间较长才写出来的毕业设计
// @termsOfService github.com/lemon997/lemonMall

func main() {
	gin.SetMode(global.ServerSetting.RunMode)
	router := routers.NewRouter()
	s := &http.Server{
		Addr:           "192.168.172.100:" + global.ServerSetting.HttpPort,
		Handler:        router,
		ReadTimeout:    global.ServerSetting.ReadTimeout,
		WriteTimeout:   global.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()

}

func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = setting.ReadSection("Server", &global.ServerSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("App", &global.AppSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("Database", &global.DatabaseSetting)
	if err != nil {
		return err
	}

	global.ServerSetting.ReadTimeout *= time.Second
	global.ServerSetting.WriteTimeout *= time.Second

	//jwt
	err = setting.ReadSection("JWT", &global.JWTSetting)
	if err != nil {
		return err
	}
	global.JWTSetting.Expire *= time.Second

	//redis
	err = setting.ReadSection("Redis", &global.RedisSetting)
	if err != nil {
		return err
	}
	err = setting.ReadSection("RabbitMQ", &global.RabbitMQSetting)
	if err != nil {
		return err
	}
	return nil
}

func setupDBEngine() error {
	//注意不能重新声明Global.DBEngine
	var err error
	global.DBEngine, err = model.NewDBEngine(global.DatabaseSetting)
	if err != nil {
		return err
	}
	global.DBEngineShop, err = model.NewDBEngineShop(global.DatabaseSetting)
	if err != nil {
		return err
	}

	return nil
}

func setupRedis() error {
	ctx := context.Background()
	rdb := global.NewRedis(global.RedisSetting)

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

func setupLogger() error {
	global.Logger = logger.NewLogger(&lumberjack.Logger{
		Filename:  global.AppSetting.LogSavePath + "/" + global.AppSetting.LogFileName + global.AppSetting.LogFileExt,
		MaxSize:   600,
		MaxAge:    10,
		LocalTime: true,
	}, "", log.LstdFlags).WithCaller(2)

	return nil
}

func snowflake(ctx context.Context) error {
	var err error
	global.Node1, err = global.NewNode(1, ctx)
	if err != nil {
		return err
	}
	global.Node2, err = global.NewNode(2, ctx)
	if err != nil {
		return err
	}
	return nil
}

func SetStock() error {
	// redisDao.NewRedisDao(context.Background()).DelAll()
	svc := service.New(context.Background())
	err := svc.SetStockAll()
	return err
}

func InitQueue() error {
	err := mq.MessageQueue()
	err = mq.DelayQueue()
	return err
}
