package dao

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"winter-test/model"
)

var DB *gorm.DB
var db *sql.DB

func Opendata() {
	//连接数据库
	err := InitDB()
	if err != nil {
		fmt.Printf("------err %v\n", err)
	} else {
		fmt.Println("---------连接成功")
	}
	fmt.Printf("db: %v\n", db)
	err = InitGormDB()
	if err != nil {
		fmt.Printf("------err %v\n", err)
	} else {
		fmt.Println("---------连接成功")
	}
	fmt.Printf("db: %v\n", db)
}
func InitGormDB() (err error) {
	dsn := "sql43_139_195_1:TA5G28rHsB@tcp(localhost:3306)/sql43_139_195_1?charset=utf8mb4"

	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	} else {
		fmt.Println("---------连接数据库成功---------", DB)
	}
	// 迁移表
	err = DB.AutoMigrate(&model.User{}, &model.Cart{}, &model.Favorite{}, &model.Order{}, &model.Style{}, &model.Comment{}, &model.Product{}, &model.RotationProduct{}, &model.OrderItem{}, &model.Seller{}, &model.Address{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return
}
func InitDB() (err error) {
	dsn := "sql43_139_195_1:TA5G28rHsB@tcp(localhost:3306)/sql43_139_195_1?charset=utf8mb4"

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 与数据库建立连接
	err2 := db.Ping()
	if err2 != nil {
		return err2
	}
	return nil

}
