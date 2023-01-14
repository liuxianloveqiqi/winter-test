package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"winter-test/dao"
	"winter-test/model"
	"winter-test/service"
)

var username, password string

func UserRoute(r *gin.Engine) {
	// 用户路业组准备
	us := r.Group("/suning/user")
	{
		us.POST("/register", Register)                                              //注册
		us.POST("/login", Login)                                                    //登录
		us.GET("logout", service.JwtAuthMiddleware(), Logout)                       //退出
		us.POST("/secret", SecretQurry)                                             //通过密保重置密码
		us.POST("/resetpassword", service.JwtAuthMiddleware(), ResetPassword)       //修改密码
		us.GET("/favorites/:user_name", service.JwtAuthMiddleware(), ShowFavorites) //展示用户收藏
	}
}

// 进行用户注册
func Register(c *gin.Context) {
	var userregitser model.UserRegister
	err := c.ShouldBind(&userregitser)
	if err != nil {
		c.JSON(400, gin.H{"status": 401,
			"info": "error",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": 200,
			"info":   "success",
		})
	}
	// 使用md5加密密码
	password := userregitser.Password
	dao.Register(&userregitser, service.Md5(password))
	// 注册成功重定向到登录界面
	c.Redirect(http.StatusFound, "/user/login")
}

// 进行用户登录
func Login(c *gin.Context) {
	username = c.PostForm("username")
	password = c.PostForm("password")
	//检查用户名是否存在
	if !service.CheckUsernamelive(c, username) {
		fmt.Println("用户名不存在")
		return
	}
	// 验证密码是否符合

	if !dao.CheckLogin(username, service.Md5(password)) {
		// 验证失败，返回错误信息
		c.JSON(401, gin.H{
			"status": 401,
			"info":   "error",
			"data": gin.H{
				"error": "密码错误",
			},
		})
		return
	} else {
		//验证成功生成token
		tokenString, _ := service.GetToken(username)
		c.JSON(http.StatusOK, gin.H{
			"status":   200,
			"info":     "success",
			"token":    tokenString,
			"username": username,
		})
	}
	c.Redirect(301, "/store")
}

// 进行用户退出
func Logout(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": 200,
		"info":   "success",
	})
	// 重定向到网站首页
	c.Redirect(301, "/store")
}

// 密保重置密码
func SecretQurry(c *gin.Context) {
	username1 := c.PostForm("username")
	// 先验证用户名是否存在
	if !service.CheckUsernamelive(c, username1) {
		fmt.Println("用户名不存在")
		c.Redirect(301, "/store/login")
		return
	}
	service.SecretQurry(c, username1, password)
}

// 修改密码
func ResetPassword(c *gin.Context) {
	newPassword := c.PostForm("newpassword")
	if len(newPassword) < 4 || len(newPassword) > 15 {
		c.JSON(400, gin.H{
			"status": 401,
			"info":   "error",
			"data": gin.H{
				"error": "密码长度应大于等于4小于等于15",
			},
		})
		return
	}
	service.ResetPassword(c, username, newPassword)
}
