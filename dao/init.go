package dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"winter-test/model"
)

var db *gorm.DB

func Opendata() {
	//连接数据库
	err := InitDB()
	if err != nil {
		fmt.Printf("------err %v\n", err)
	} else {
		fmt.Println("---------连接成功")
	}
	fmt.Printf("db: %v\n", db)
}
func InitDB() (err error) {
	dsn := "sql43_139_195_1:TA5G28rHsB@tcp(localhost:3306)/sql43_139_195_1?charset=utf8mb4"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	} else {
		fmt.Println("---------连接数据库成功---------", db)
	}
	// 迁移表
	err = db.AutoMigrate(&model.User{}, &model.Cart{}, &model.Favorite{}, &model.Order{}, &model.Style{}, &model.Comment{}, &model.Product{}, &model.RotationProduct{}, &model.OrderItem{}, &model.Seller{}, &model.Address{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return
}
