package mysql

import (
	"context"
	"fmt"
	"govote/app/config"
	"govote/app/db/model"
	"govote/app/tools/log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var _db *gorm.DB

func NewMysql() {
	my := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Conf.MySQLConfig.User,
		config.Conf.MySQLConfig.Password,
		config.Conf.MySQLConfig.Host,
		config.Conf.MySQLConfig.Port,
		config.Conf.MySQLConfig.DbName)
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
	err = conn.AutoMigrate(&model.Vote{}, &model.User{}, &model.VoteOpt{}, &model.VoteOptUser{})
	if err != nil {
		log.L.Panicf("数据表AutoMigrate失败, err:%s\n", err)
	}

	// 设置连接池配置
	sqlDB, err := conn.DB()
	if err != nil {
		log.L.Panicf("获取数据库实例失败, err:%s\n", err)
	}
	sqlDB.SetMaxOpenConns(config.Conf.MySQLConfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.Conf.MySQLConfig.MaxIdleConns)

	_db = conn
}

func NewDBClient(ctx context.Context) *gorm.DB {
	db := _db
	return db.WithContext(ctx)
}

func Close() {
	db, _ := _db.DB()
	_ = db.Close()
}
