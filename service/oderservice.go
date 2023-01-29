package service

import (
	"github.com/gin-gonic/gin"
	"winter-test/dao"
)

func CreatOrder(username string, c *gin.Context) {
	orderID, err := dao.CreateOrder(settlecarts, totalprice, username)
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
				"orderID": orderID,
			},
		})
	}
}
