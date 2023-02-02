package dao

import (
	"fmt"
	"winter-test/model"
)

// 插入用户注册的数据
func Register(u *model.UserRegister, passwordhash16 string) {
	sqlStr := "insert into user(username,password,nickname,SecretQ,SecretA) values (?,?,?,?,?)"
	r, err := db.Exec(sqlStr, u.UserName, passwordhash16, u.Nickname, u.SecretQ, u.SecretA)
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

// 展示用户资料
func GetUserMessage(username string) (model.UserMessage, error) {
	var userMessage model.UserMessage
	err := db.QueryRow("select human_name, phone_number,nickname, email, gender from user where username = ?", username).Scan(&userMessage.HumanName, &userMessage.PhoneNumber, &userMessage.Nickname, &userMessage.Email, &userMessage.Gender)
	if err != nil {
		return userMessage, err
	}
	return userMessage, nil
}

// 用户修改资料
func UpdateUserMessage(userMessage *model.UserMessage, username string) error {
	query := `update user set`
	fmt.Println(userMessage, "9999999")
	// 通过判断传入的结构体的字段是否为空，来决定是否在sql语句中更新对应的字段
	if userMessage.HumanName != "" {
		query += ` human_name = '%s',`
		query = fmt.Sprintf(query, userMessage.HumanName)
	}
	if userMessage.PhoneNumber != 0 {
		query += ` phone_number = '%d',`
		query = fmt.Sprintf(query, userMessage.PhoneNumber)
	}
	if userMessage.Email != "" {
		query += ` email = '%s',`
		query = fmt.Sprintf(query, userMessage.Email)
	}
	if userMessage.Gender != "" {
		query += ` gender = '%s',`
		query = fmt.Sprintf(query, userMessage.Gender)
	}
	if userMessage.Nickname != "" {
		query += ` nickname = '%s',`
		query = fmt.Sprintf(query, userMessage.Nickname)
	}
	// 将sql语句最后一个逗号删除，并将username查询条件拼接到语句后面
	query = query[:len(query)-1] + " where username = '%s';"
	query = fmt.Sprintf(query, username)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil

}

// 查看余额
func GetMoney(username string) (float64, error) {
	query := "select money from user where username = ?"
	var money float64
	err := db.QueryRow(query, username).Scan(&money)
	if err != nil {
		return 0, err
	}
	return money, nil
}

// 充值余额
func AddMoney(username string, m float64) error {
	query := `update user set money = money + ? where username = ?`
	_, err := db.Exec(query, m, username)
	if err != nil {
		return err
	}
	return nil
}

// 新增收获地址
func CreatAddress(username, place string) (model.Address, error) {
	var address model.Address
	query := "insert into address (user_name, place) values (?, ?)"
	_, err := db.Exec(query, username, place)
	if err != nil {
		return address, err
	}
	err = db.QueryRow("select id,user_name,place from address where user_name = ? and place = ?", username, place).Scan(&address.ID, &address.UserName, &address.Place)
	if err != nil {
		return address, err
	}
	return address, nil
}

// 删除收货地址
func DeleteAddress(username string, addressID int) error {
	// 删除该收货地址
	_, err := db.Exec("delete from address where id = ? and user_name = ?", addressID, username)
	if err != nil {
		return err
	}

	return nil
}
