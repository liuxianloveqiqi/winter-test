package main

import (
	"github.com/gin-gonic/gin"
	"winter-test/api"
	"winter-test/dao"
)

func main() {
	//打开数据库
	dao.Opendata()
	//打开路由
	r := gin.Default()
	api.UserRoute(r)
	api.ProductsRoute(r)
	api.CartRoute(r)
	api.OderRoute(r)
	api.CommentRoute(r)
	r.Run(":9090")
}
