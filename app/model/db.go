package model

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// 数据库操作都放在这里

var Conn *gorm.DB

func NewMysql() {
	my := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "123456", "localhost:3306", "vote")
	var ormLogger logger.Interface
	ormLogger = logger.Default.LogMode(logger.Info)
	conn, err := gorm.Open(mysql.Open(my), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		fmt.Printf("err:%s\n", err)
		panic(err)
	}
	err = conn.AutoMigrate(&Vote{}, &User{}, &VoteOpt{}, &VoteOptUser{})
	if err != nil {
		return
	}
	Conn = conn
}

func Close() {
	db, _ := Conn.DB()
	_ = db.Close()
}
