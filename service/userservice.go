package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
	"winter-test/dao"
	"winter-test/model"
)

// 检查用户名是否存在
func CheckUsernamelive(c *gin.Context, u string) bool {

	if !dao.QuerryUsername(u) {
		c.JSON(401, gin.H{
			"status": 401,
			"info":   "error",
			"data": gin.H{
				"error": "用户名不存在",
			},
		})
	}
	return dao.QuerryUsername(u)
}

// md5加密
func Md5(pasaword string) string {
	hash := md5.New()
	hash.Write([]byte(pasaword))
	passwordHash := hash.Sum(nil)
	// 将哈希密码转换为16进制储存
	passwordHash16 := hex.EncodeToString(passwordHash)
	return passwordHash16
}

// 根据用户名及密保问题重置密码
func SecretQurry(c *gin.Context, u string, pa string) {
	Q := dao.SecretQurryUsername(u)
	fmt.Println("###########", pa)
	c.JSON(200, gin.H{
		"status": 200,
		"info":   "success",
		"data": gin.H{
			"你的密保答案": Q,
		},
	})
	A := c.PostForm("secretA")
	is, phash := dao.SecreQurryA(u, Q, A)
	fmt.Println("$$$$$$$$$$$$$$$$$$$", phash)
	if A == "" {

	} else if !is {
		c.JSON(401, gin.H{
			"status": 401,
			"info":   "error",
			"data": gin.H{
				"error": "输入的密保答案错误,请重新输入",
			},
		})
	} else {
		newPassword := c.PostForm("newpassword")
		if len(newPassword) < 6 || len(newPassword) > 15 {
			c.JSON(400, gin.H{
				"status": 400,
				"info":   "error",
				"data": gin.H{
					"error": "密码长度应大于等于6小于等于15",
				},
			})
			return
		}
		ResetPassword(c, u, newPassword)
	}
}

// // 中间件cookie凭证
//
//	func AuthMiddleWare(username string) gin.HandlerFunc {
//		return func(c *gin.Context) {
//			// 获取客户端cookie并校验
//			cookie, err := c.Cookie("username")
//			fmt.Println("oooooooooooooo", cookie)
//			fmt.Println(username)
//			if cookie == username && err == nil {
//				c.Next()
//			} else {
//				// 返回错误
//				c.JSON(http.StatusUnauthorized, gin.H{"error": "没有登录"})
//				// 若验证不通过，不再调用后续的函数处理
//				c.Abort()
//			}
//		}
//	}
//

// JWT中间件，验证 token

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(401, gin.H{
				"status": 401,
				"info":   "error",
				"error":  "请求未携带token，无权限访问",
			})
			c.Abort()
			return
		}
		claims, err := ParseToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{
				"status": 401,
				"info":   "error",
				"data": gin.H{
					"error": "无效的token",
				},
			})
			c.Abort()
			return
		}

		username := claims.UserName
		if username == "" {
			c.JSON(401, gin.H{
				"status": 401,
				"info":   "error",
				"data": gin.H{
					"error": "无效的token",
				},
			})
			c.Abort()
			return
		}
		// 验证通过
		c.Set("claims", claims)
		c.Next()
	}
}

const TokenExpireTime = time.Hour * 24 //设置过期时间

var Secret = []byte("liuxian123") //设置密钥

func GetToken(username string) (string, error) {
	// 创建一个Claims
	c := model.MyClaims{
		username, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireTime).Unix(), // 过期时间
			Issuer:    "liuxian",                              // 签发人
		},
	}
	// 使用HS256算法签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用secret签名并获得字符串token
	return token.SignedString(Secret)
}

// 解析token
func ParseToken(tokenString string) (*model.MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.MyClaims{},
		func(token *jwt.Token) (i interface{}, err error) {
			return Secret, nil
		})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*model.MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token验证错误")
}

// 修改密码
func ResetPassword(c *gin.Context, u, np string) {
	err := dao.ResetPassword(u, Md5(np))
	if len(np) < 6 || len(np) > 15 {
		c.JSON(400, gin.H{
			"status": 400,
			"info":   "error",
			"data": gin.H{
				"error": "密码长度应大于等于6小于等于15",
			},
		})
		return
	}

	if err != nil {
		fmt.Println("修改密码错误：", err)
	} else {
		c.JSON(200, gin.H{
			"status": 200,
			"info":   "success",
			"data": gin.H{
				"你的新密码": np,
			},
		})
	}
}
