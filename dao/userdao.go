package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"winter-test/model"
)

// 插入用户注册的数据
func Register(u *model.UserRegister, passwordhash16 string) {
	user := model.User{
		UserName: u.UserName,
		Password: u.Password,
		NickName: u.Nickname,
		SecretQ:  u.SecretQ,
		SecretA:  u.SecretA,
	}

	result := DB.Create(&user)
	if result.Error != nil {
		fmt.Printf("err: %v\n", result.Error)
		return
	}
	fmt.Printf("username: %v\n", user.UserName)

}

// 查询用户名是否存在
func QuerryUsername(u string) bool {
	var count int64
	result := DB.Model(&model.User{}).Where("username = ?", u).Count(&count)
	if result.Error != nil {
		fmt.Printf("err: %v\n", result.Error)
		return false
	}
	return count == 1
}

// 用户登录密码验证
func CheckLogin(u, p string) bool {
	var count int64
	result := DB.Model(&model.User{}).Where("username = ? and password = ?", u, p).Count(&count)
	if result.Error != nil {
		fmt.Printf("err: %v\n", result.Error)
		return false
	}
	return count == 1
}

// 根据用户名查询密保问题
func SecretQurryUsername(u string) string {
	var Q string
	result := DB.Model(&model.User{}).Select("secretQ").Where("username = ?", u).First(&Q)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return ""
		}
		fmt.Printf("err: %v\n", result.Error)
		return ""
	}
	return Q
}

// 通过密保查密码
func SecreQurryA(u, Q, A string) (bool, string) {

	var userCountAndPwd model.UserCountAndPwd
	result := DB.Model(&model.User{}).Select("count(*), password").Where("username = ? and secretQ = ? and secretA = ?", u, Q, A).Group("username").Scan(&userCountAndPwd)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, ""
		}
		fmt.Printf("err: %v\n", result.Error)
		return false, ""
	}
	return userCountAndPwd.Count == 1, userCountAndPwd.Password

}

// 修改密码
func ResetPassword(u, np string) error {
	result := DB.Model(&model.User{}).Where("username = ?", u).Update("password", np)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user %s not found", u)
		}
		fmt.Printf("err: %v\n", result.Error)
		return result.Error
	}
	return nil
}

// 展示用户资料
func GetUserMessage(username string) (model.UserMessage, error) {
	var userMessage model.UserMessage
	err := DB.Table("user").Where("username = ?", username).Select("human_name, phone_number, nickname, email, gender").Scan(&userMessage).Error
	if err != nil {
		return userMessage, err
	}
	return userMessage, nil
}

// 用户修改资料
func UpdateUserMessage(userMessage *model.UserMessage, username string) error {
	err := DB.Table("user").Where("username = ?", username).
		Updates(map[string]interface{}{
			"human_name":   userMessage.HumanName,
			"phone_number": userMessage.PhoneNumber,
			"email":        userMessage.Email,
			"gender":       userMessage.Gender,
			"nickname":     userMessage.Nickname,
		}).Error
	if err != nil {
		return err
	}
	return nil
}

// 查看余额
func GetMoney(username string) (float64, error) {
	var money float64
	err := DB.Table("user").Where("username = ?", username).Pluck("money", &money).Error
	if err != nil {
		return 0, err
	}
	return money, nil
}

// 充值余额
func AddMoney(username string, m float64) error {
	err := DB.Table("user").Where("username = ?", username).UpdateColumn("money", gorm.Expr("money + ?", m)).Error
	if err != nil {
		return err
	}
	return nil
}

// 新增收获地址
func CreatAddress(username, place string) (model.Address, error) {
	var address model.Address
	address.UserName = username
	address.Place = place
	result := DB.Create(&address)
	if result.Error != nil {
		return address, result.Error
	}
	return address, nil
}

// 删除收货地址
func DeleteAddress(username string, addressID int) error {
	result := DB.Table("address").Where("id = ? and user_name = ?", addressID, username).Delete(&model.Address{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
