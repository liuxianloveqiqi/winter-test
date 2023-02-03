package service

import (
	"github.com/gin-gonic/gin"
	"winter-test/dao"
)

// 创建订单
func CreateOrder(username string, address_id int, c *gin.Context) {
	order, err := dao.CreateOrder(settlecarts, totalprice, username, address_id, ids)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"info":   "error",
			"data": gin.H{
				"error": err,
			},
		})
		return
	} else {
		c.JSON(200, gin.H{
			"status": 200,
			"info":   "success",
			"data": gin.H{
				"order": order,
			},
		})
	}
}
