package model

import (
	"context"
	"fmt"

	"govote/app/tools/log"
	"govote/app/tools/session"

	"github.com/rbcervilla/redisstore/v9"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 数据库操作都放在这里

var Conn *gorm.DB

func NewMysql() {
	my := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "123456", "localhost:3306", "vote")
	ormLogger := logger.Default.LogMode(logger.Info)
	conn, err := gorm.Open(mysql.Open(my), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		log.L.Panicf("数据库变量初始化失败, err:%s\n", err)
	}
	err = conn.AutoMigrate(&Vote{}, &User{}, &VoteOpt{}, &VoteOptUser{})
	if err != nil {
		log.L.Panicf("数据表AutoMigrate失败, err:%s\n", err)
	}
	Conn = conn
}

var Rdb *redis.Client

func NewRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       8,  // use default DB
	})
	Rdb = rdb
	// 初始化session
	var err error
	session.SessionStore, err = redisstore.NewRedisStore(context.TODO(), Rdb)
	if err != nil {
		log.L.Panicf("初始化redisStore失败, err:%s\n", err)
	}
}

func Close() {
	db, _ := Conn.DB()
	_ = db.Close()
}
