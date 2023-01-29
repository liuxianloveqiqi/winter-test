package model

import "github.com/dgrijalva/jwt-go"

type UserRegister struct {
	UserName string `form:"username" binding:"required,min=6,max=15"`
	LockName string `form:"lockname" binding:"required,min=4,max=15"`
	Password string `form:"password" binding:"required,min=6,max=15"`
	SecretQ  string `form:"secretQ" binding:"required"`
	SecretA  string `form:"secretA" binding:"required"`
}
type MyClaims struct {
	UserName string `json:"username"`
	jwt.StandardClaims
}
type UserMessage struct {
	HumanName   string `json:"human_name" form:"human_name"`
	PhoneNumber int    `json:"phone_number" form:"phone_number"`
	Email       string `json:"email" form:"email"`
	Gender      string `json:"gender" form:"gender"`
	LockName    string `json:"lock_name" form:"lock_name"`
}
