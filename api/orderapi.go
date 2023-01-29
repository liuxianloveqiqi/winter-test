package api

import (
	"github.com/gin-gonic/gin"
	"winter-test/model"
	"winter-test/service"
)

func OderRounte(r *gin.Engine) {
	// 订单功能
	order := r.Group("suning/order")
	{
		order.GET("/add", service.JwtAuthMiddleware(), Creatorder) //创造订单

	}

}

// 建立一个订单
func Creatorder(c *gin.Context) {
	un := c.MustGet("claims").(*model.MyClaims).UserName
	var usermessage model.UserMessage
	c.ShouldBind(&usermessage)

	service.CreatOrder(un, c)
}
