package service

import "C"
import (
	"github.com/gin-gonic/gin"
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

// 根据用户名及密保问题查询密码
func SecretQurry(c *gin.Context, u string) {
	Q := dao.SecretQurryUsername(u)
	c.JSON(200, gin.H{
		"status":         200,
		"你的密保问题为": Q,
	})
	A := c.PostForm("secretA")
	is, p := dao.SecreQurryA(u, Q, A)
	if A == "" {

	} else if !is {
		c.JSON(401, gin.H{
			"status": 401,
			"error":  "输入的密保答案错误,请重新输入",
		})
	} else {
		c.JSON(200, gin.H{
			"status":     200,
			"你的密码是": p,
		})
	}
}
