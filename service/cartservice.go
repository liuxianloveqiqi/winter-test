package service

import (
	"github.com/gin-gonic/gin"
	"winter-test/dao"
	"winter-test/model"
)

var settlecarts []*model.Cart
var totalprice float64
var ids []int

func SettleCart(c *gin.Context, request model.SettleCartRequest) {
	// 调用dao层函数进行结算操作
	var err error
	ids = request.IDs
	settlecarts, totalprice, err = dao.SettleCart(request.IDs)
	if err != nil {
		c.JSON(400, gin.H{
			"status": 400,
			"info":   "error",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  200,
		"message": "Success",
		"data":    settlecarts,
	})

}
