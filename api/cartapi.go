package api

import (
	"github.com/gin-gonic/gin"
	"winter-test/dao"
	"winter-test/model"
	"winter-test/service"
)

func CartRoute(r *gin.Engine) {
	// 实现购物车功能
	cart := r.Group("/suning/cart")
	{
		cart.POST("/add", service.JwtAuthMiddleware(), AddCart)         //添加商品到购物车
		cart.GET("/list", service.JwtAuthMiddleware(), ListCart)        //展示购物车中的商品
		cart.DELETE("/remove", service.JwtAuthMiddleware(), RemoveCart) //删除购物车中的商品
		cart.POST("/settle", service.JwtAuthMiddleware(), SettleCart)   //选择商品进行结算
	}
}

// 将商品加入购物车
func AddCart(c *gin.Context) {
	var cart model.SmallCart
	if err := c.ShouldBind(&cart); err != nil {
		c.JSON(400, gin.H{
			"status": 400,
			"info":   "error",
			"data": gin.H{
				"error": "请求参数错误",
			},
		})
		return
	}
	cart.UserName = c.MustGet("claims").(*model.MyClaims).UserName

	if err := dao.AddCart(&cart); err != nil {
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
		"data":   "加入购物车成功",
	})
}

// 将商品删除购物车
func RemoveCart(c *gin.Context) {
	var cart model.SmallCart
	if err := c.ShouldBind(&cart); err != nil {
		c.JSON(400, gin.H{
			"status": 400,
			"info":   "error",
			"data": gin.H{
				"error": "请求参数错误",
			},
		})
		return
	}
	cart.UserName = c.MustGet("claims").(*model.MyClaims).UserName
	if err := dao.RemoveCart(&cart); err != nil {
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
		"data":   "删除购物车成功",
	})
}

// 列举用户的购物车
func ListCart(c *gin.Context) {
	claims := c.MustGet("claims").(*model.MyClaims)
	carts, err := dao.ListCart(claims.UserName)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"info":   "error",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  200,
		"message": "Success",
		"data":    carts,
	})
}

// 结算购物车中的部分商品
func SettleCart(c *gin.Context) {
	var request model.SettleCartRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(400, gin.H{
			"status": 400,
			"info":   "error",
			"error":  err.Error(),
		})
		return
	}

	service.SettleCart(c, request)
}
