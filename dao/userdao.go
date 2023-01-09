package dao

import (
	"fmt"
	"winter-test/model"
)

// 插入用户注册的数据
func Register(u *model.UserRegister, passwordhash16 string) {
	sqlStr := "insert into user(username,password,lockname,SecretQ,SecretA) values (?,?,?,?,?)"
	r, err := db.Exec(sqlStr, u.UserName, passwordhash16, u.LockName, u.SecretQ, u.SecretA)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	i2, err2 := r.LastInsertId()
	if err2 != nil {
		fmt.Printf("err2: %v\n", err2)
		return
	}
	fmt.Printf("--------i2: %v\n", i2)

}

// 查询用户名是否存在
func QuerryUsername(u string) bool {
	var count int
	err := db.QueryRow("select count(*) from user where username = ?", u).Scan(&count)

	return err == nil && count == 1
}

// 用户登录密码验证
func CheckLogin(u, p string) bool {
	var count int
	//查询验证用户名和密码是否匹配
	err := db.QueryRow("select count(*) from user where username = ? and password = ?", u, p).Scan(&count)
	return err == nil && count == 1
}

// 根据用户名查询密保问题
func SecretQurryUsername(u string) string {
	var Q string
	err := db.QueryRow("select secretQ from user where username = ?", u).Scan(&Q)
	if err != nil {
		panic(err)
	}
	return Q
}

// 通过密保查密码
func SecreQurryA(u, Q, A string) (bool, string) {
	var count int
	var p string
	//查询验证用户密保
	err := db.QueryRow("select count(*),password from user where username = ? and secretQ = ? and secretA = ? group by username", u, Q, A).Scan(&count, &p)
	fmt.Println(err, "----------------", count)
	return err == nil && count == 1, p
}

// 修改密码
func ResetPassword(u, np string) error {
	strSql := "update user set password=? where username=?"
	_, err := db.Exec(strSql, np, u)
	return err
}
