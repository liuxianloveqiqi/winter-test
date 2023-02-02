package service

import (
	"github.com/gin-gonic/gin"
	"winter-test/dao"
)

// 添加评论
func CreateComment(c *gin.Context, username, comment string, productID int, parentID int) {
	// 向数据库中存储评论信息
	commentAll, err := dao.CreateComment(username, productID, comment, parentID)
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
			"comment": commentAll,
		},
	})
}
