package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"winter-test/dao"
	"winter-test/model"
	"winter-test/service"
)

func OderRounte(r *gin.Engine) {
	// 订单功能
	order := r.Group("suning/order")
	{
		order.GET("/add", service.JwtAuthMiddleware(), ShowMessage)            //展示收货地址
		order.GET("/add/:address_id", service.JwtAuthMiddleware(), Creatorder) //创造订单

	}

}

// 展示收货信息
func ShowMessage(c *gin.Context) {
	un := c.MustGet("claims").(*model.MyClaims).UserName
	addresses, err := dao.ShowMessage(un)
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
		"data":   addresses,
	})
}

// 建立一个订单
func Creatorder(c *gin.Context) {
	un := c.MustGet("claims").(*model.MyClaims).UserName
	id := c.Param("address_id")
	address_id, _ := strconv.Atoi(id)
	service.CreatOrder(un, address_id, c)
}
