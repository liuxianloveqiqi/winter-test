package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

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
}
func InitDB() (err error) {
	dsn := "root:xian712525@tcp(127.0.0.1:3306)/store?charset=utf8mb4"

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
