package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"winter-test/dao"
)

// 商品详情
func Productdata(c *gin.Context, id string) {
	products, err := dao.Productdata(id)
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
			"data":   products,
		})
	}
}

// 商品款式
func GetStyles(c *gin.Context, id string) {
	styles, err := dao.GetStyles(id)
	fmt.Println("########", err)
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
			"data":   styles,
		})
	}
}
