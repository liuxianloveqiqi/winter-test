package service

import "C"
import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"winter-test/dao"
)

// 检查用户名是否存在
func CheckUsernamelive(c *gin.Context, u string) bool {

	if !dao.QuerryUsername(u) {
		c.JSON(401, gin.H{
			"error":  "没有此用户名",
			"status": 401})
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
		"status":  200,
		"你的密保问题为": Q,
	})
	A := c.PostForm("secretA")
	is, phash := dao.SecreQurryA(u, Q, A)
	fmt.Println("$$$$$$$$$$$$$$$$$$$", phash)
	if A == "" {

	} else if !is {
		c.JSON(401, gin.H{
			"status": 401,
			"error":  "输入的密保答案错误,请重新输入",
		})
	} else {
		newPassword := c.PostForm("newpassword")
		if len(newPassword) < 4 || len(newPassword) > 15 {
			c.JSON(400, gin.H{
				"error": "密码长度应大于等于4小于等于15",
			})
			return
		}
		ResetPassword(c, u, newPassword)
	}
}

// 中间件cookie凭证
func AuthMiddleWare(username string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取客户端cookie并校验
		if cookie, err := c.Cookie("username"); err == nil {
			if cookie == username {
				c.Next()
			}
		} else {
			// 返回错误
			c.JSON(http.StatusUnauthorized, gin.H{"error": "没有登录"})
			// 若验证不通过，不再调用后续的函数处理
			c.Abort()
		}
	}
}

// 修改密码
func ResetPassword(c *gin.Context, u, np string) {
	err := dao.ResetPassword(u, Md5(np))
	if err != nil {
		fmt.Println("修改密码错误：", err)
	} else {
		c.JSON(200, gin.H{
			"status": 200,
			"你的新密码为": np,
		})
	}
}
