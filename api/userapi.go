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
	us := r.Group("/user")
	{
		us.POST("/register", Register)                                             //注册
		us.POST("/login", Login)                                                   //登录
		us.GET("logout", Logout)                                                   //退出
		us.POST("/secret", SecretQurry)                                            //通过密保找回密码
		us.POST("/resetpassword", service.AuthMiddleWare(username), ResetPassword) //修改密码
	}
}

// 进行用户注册
func Register(c *gin.Context) {
	var userregitser model.UserRegister
	err := c.ShouldBind(&userregitser)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "注册成功",
			"status":  200,
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

			"error":  "密码错误",
			"status": 401})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "登录成功",
			"hello":   username,
			"status":  200,
		})

	}
	// 验证成功，设置 cookie 并跳转到首页

	c.SetCookie("username", username, 3600*3, "/", "localhost", false, true)
	c.Redirect(301, "/store")
}

// 进行用户退出
func Logout(c *gin.Context) {
	// 删除用户的cookie
	c.SetCookie("username", username, -1, "/", "localhost", false, true)
	// 重定向到网站首页
	c.Redirect(301, "/store")
}

// 密保查询
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
	service.ResetPassword(username, newPassword)
}
