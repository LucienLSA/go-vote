package model

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func Init() error {
	my := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root", "123456", "127.0.0.1:3306", "vote")
	conn, err := gorm.Open(mysql.Open(my), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		fmt.Printf("数据库连接失败: %s\n", err)
		return err
	}

	// 先赋值DB变量，再进行自动迁移
	DB = conn

	err = DB.AutoMigrate(&User{})
	if err != nil {
		fmt.Printf("数据库迁移失败: %s\n", err)
		return err
	}
	return nil
}

func Close() {
	db, _ := DB.DB()
	_ = db.Close()
}
