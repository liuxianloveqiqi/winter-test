package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"winter-test/dao"
	"winter-test/model"
	"winter-test/service"
)

func OderRoute(r *gin.Engine) {
	// 订单功能
	order := r.Group("suning/order")
	{
		order.GET("/add", service.JwtAuthMiddleware(), ShowMessage)             //展示收货地址
		order.GET("/add/:address_id", service.JwtAuthMiddleware(), CreateOrder) //创造订单
		order.GET("/show", service.JwtAuthMiddleware(), ShowOrdersByStep)       //分类展示订单
		order.POST("/update", service.JwtAuthMiddleware(), UpdateOrderStep)     //修改订单状态
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
func CreateOrder(c *gin.Context) {
	un := c.MustGet("claims").(*model.MyClaims).UserName
	id := c.Param("address_id")
	address_id, _ := strconv.Atoi(id)
	service.CreateOrder(un, address_id, c)
}

// 根据订单状态分类展示订单
func ShowOrdersByStep(c *gin.Context) {
	stepStr := c.Query("step")
	step, err := strconv.Atoi(stepStr)
	if err != nil {
		c.JSON(400, gin.H{
			"status": 400,
			"info":   "error",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	un := c.MustGet("claims").(*model.MyClaims).UserName
	orders, err := dao.ShowOrdersByStep(un, step)
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
		"data":   orders,
	})
}

// 改变订单状态
func UpdateOrderStep(c *gin.Context) {
	idStr := c.PostForm("order_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(400, gin.H{
			"status": 400,
			"info":   "error",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	stepStr := c.PostForm("step")
	step, err := strconv.Atoi(stepStr)
	if err != nil {
		c.JSON(400, gin.H{
			"status": 400,
			"info":   "error",
			"data": gin.H{
				"error": err.Error(),
			},
		})
		return
	}

	un := c.MustGet("claims").(*model.MyClaims).UserName
	err = dao.UpdateOrderStep(id, un, step)
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
		"data":   fmt.Sprintf("修改订单id:%v状态为%v", id, step),
	})
}
