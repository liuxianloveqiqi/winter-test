package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"winter-test/dao"
	"winter-test/model"
	"winter-test/service"
)

var username, password string
var conf = model.Conf{
	"b0c68d641806bfde2460",
	"8687d10848b4ca4c4afcf2f3cea351eb14c5f313",
	"http://43.139.195.17:9090/auth/callback/github",
}

// 创建 GitHub OAuth2 配置

func UserRoute(r *gin.Engine) {

	r.GET("/auth/login/github", HandleGithubLogin)
	r.GET("/auth/callback/github", HandleGithubCallback)

	// 用户路业组准备
	us := r.Group("/suning/user")
	{
		us.POST("/register", Register) // 注册

		us.POST("/login", Login)                                                      // 登录
		us.GET("logout", service.JwtAuthMiddleware(), Logout)                         // 退出
		us.POST("/secret", SecretQurry)                                               // 通过密保重置密码
		us.POST("/resetpassword", service.JwtAuthMiddleware(), ResetPassword)         // 修改密码
		us.GET("/favorites/:user_name", service.JwtAuthMiddleware(), ShowFavorites)   // 展示用户收藏
		us.GET("/message", service.JwtAuthMiddleware(), GetUserMessage)               // 展示用户资料
		us.POST("/message", service.JwtAuthMiddleware(), UpdateUserMessage)           // 修改用户资料
		us.GET("/money", service.JwtAuthMiddleware(), GetMoney)                       // 查看余额
		us.POST("/money", service.JwtAuthMiddleware(), AddMoney)                      // 充值余额
		us.POST("/address", service.JwtAuthMiddleware(), CreatAddress)                // 新增收货地址
		us.DELETE("/address/:address_id", service.JwtAuthMiddleware(), DeleteAddress) // 删除收货地址
	}
}

// 进行用户注册
func Register(c *gin.Context) {
	var userregitser model.UserRegister
	err := c.ShouldBind(&userregitser)
	if err != nil {
		c.JSON(400, gin.H{"status": 400,
			"info": "error",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	} else {
		c.JSON(200, gin.H{
			"status": 200,
			"info":   "success",
		})
	}
	// 使用md5加密密码
	password := userregitser.Password
	dao.Register(&userregitser, service.Md5(password))
	// 注册成功重定向到登录界面
	c.Redirect(302, "/user/login")
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
		tokenString, _ := service.GetToken(username, service.GenerateRandomString(10))
		c.JSON(200, gin.H{
			"status":   200,
			"info":     "success",
			"token":    tokenString,
			"username": username,
		})
	}
	c.Redirect(301, "/store")
}

// 处理第三方 Github 登录请求
func HandleGithubLogin(c *gin.Context) {

	//authURL := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", conf.ClientId, conf.RedirectUrl)
	//c.Redirect(307, authURL)
	data := gin.H{
		"ClientId":     "b0c68d641806bfde2460",
		"ClientSecret": "8687d10848b4ca4c4afcf2f3cea351eb14c5f313",
		"RedirectUrl":  "http://43.139.195.17:9090/auth/callback/github",
	}
	c.HTML(http.StatusOK, "index.html", data)
}

// 处理 GitHub 登录回调请求
func HandleGithubCallback(c *gin.Context) {
	// 从查询参数中获取授权码
	code := c.Query("code")

	// 交换授权码获取访问令牌
	tokenAuthUrl := fmt.Sprintf(
		"https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		conf.ClientId, conf.ClientSecret, code)
	// 获取 token
	var token *model.Token
	var err error
	if token, err = service.GetGithubToken(tokenAuthUrl); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("------%v", token)
	// 通过token，获取用户信息
	var userInfo map[string]interface{}
	if userInfo, err = service.GetUserInfo(token); err != nil {
		fmt.Println("获取用户信息失败，错误信息为:", err)
		return
	}
	// 返回用户信息
	c.JSON(200, gin.H{"message": "GitHub 授权登录成功",
		"token":     token.AccessToken,
		"user_info": userInfo})
}

// 进行用户退出
func Logout(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": 200,
		"info":   "success",
	})
	// 重定向到网站首页
	c.Redirect(301, "/suning")
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
			"status": 400,
			"info":   "error",
			"data": gin.H{
				"error": "密码长度应大于等于4小于等于15",
			},
		})
		return
	}
	service.ResetPassword(c, username, newPassword)
}

// 展示用户资料
func GetUserMessage(c *gin.Context) {
	// 获取token里面的username
	username := c.MustGet("claims").(*model.MyClaims).UserName
	// 从数据库中查询用户信息
	userMessage, err := dao.GetUserMessage(username)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"info":   "error",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"info":   "success",
		"data":   userMessage,
	})
}

// 修改用户资料
func UpdateUserMessage(c *gin.Context) {
	// 获取token里面的username
	username := c.MustGet("claims").(*model.MyClaims).UserName

	var userMessage model.UserMessage
	if err := c.ShouldBind(&userMessage); err != nil {
		c.JSON(400, gin.H{"status": 400,
			"info": "error",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	// 调用dao层的UpdateUserMessage函数将修改的信息更新到数据库中
	err := dao.UpdateUserMessage(&userMessage, username)
	if err != nil {
		c.JSON(500, gin.H{"status": 500,
			"info": "error",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"info":   "success",
		"data":   "修改资料成功",
	})
}

// 用户查看余额
func GetMoney(c *gin.Context) {
	// 获取token里面的username
	username := c.MustGet("claims").(*model.MyClaims).UserName
	// 从数据库中查询用户余额
	money, err := dao.GetMoney(username)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"info":   "error",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 200,
		"info":   "success",
		"data": gin.H{
			"money": money,
		},
	})
}

// 充值余额
func AddMoney(c *gin.Context) {
	// 获取token里面的username
	username := c.MustGet("claims").(*model.MyClaims).UserName

	// 获取请求参数

	a := c.PostForm("add_money")
	addmoney, err0 := strconv.ParseFloat(a, 64)
	fmt.Println(addmoney, "7777777")
	if err0 != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"info":   "error",
			"data": gin.H{
				"error": err0.Error(),
			},
		})
	}

	if err := dao.AddMoney(username, addmoney); err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"info":   "error",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"info":   "success",
		"data": gin.H{
			"message": "充值成功",
		},
	})
}

// 新增收获地址
func CreatAddress(c *gin.Context) {
	// 获取 token 中的 username
	username := c.MustGet("claims").(*model.MyClaims).UserName

	// 获取请求参数中的地址信息
	place := c.PostForm("place")

	// 向数据库中新增地址信息
	address, err := dao.CreatAddress(username, place)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"info":   "error",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"info":   "success",
		"data": gin.H{
			"address": address,
		},
	})
}

// 根据address_id删除地址
func DeleteAddress(c *gin.Context) {
	// 获取token里面的username
	username := c.MustGet("claims").(*model.MyClaims).UserName
	// 获取URL中的地址ID
	addressIDStr := c.Param("address_id")
	// 将地址ID字符串转换为int类型
	addressID, _ := strconv.Atoi(addressIDStr)
	// 在数据库中删除该收货地址
	err := dao.DeleteAddress(username, addressID)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"info":   "error",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	c.JSON(200, gin.H{
		"status": 200,
		"info":   "success",
		"data":   "成功删除地址",
	})
}
