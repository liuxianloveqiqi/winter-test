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
