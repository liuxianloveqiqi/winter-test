package model

import "github.com/dgrijalva/jwt-go"

// 用户注册
type UserRegister struct {
	UserName string `form:"username" binding:"required,min=6,max=15"`
	Nickname string `form:"nickname" binding:"required,min=4,max=15"`
	Password string `form:"password" binding:"required,min=6,max=15"`
	SecretQ  string `form:"secretQ" binding:"required"`
	SecretA  string `form:"secretA" binding:"required"`
}

type MyClaims struct {
	UserName string `json:"username" form:"username"`
	jwt.StandardClaims
}

// 用户资料
type UserMessage struct {
	HumanName   string `json:"human_name" form:"human_name"`
	PhoneNumber int    `json:"phone_number" form:"phone_number"`
	Email       string `json:"email" form:"email"`
	Gender      string `json:"gender" form:"gender"`
	Nickname    string `json:"lock_name" form:"lock_name"`
}

// 用户收货地址
type Address struct {
	ID       string `json:"id" form:"id"`
	UserName string `json:"username" form:"username"`
	Place    string `json:"place" form:"place"`
}
