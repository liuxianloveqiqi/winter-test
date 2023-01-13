package api

import (
	"github.com/gin-gonic/gin"
	"winter-test/dao"
	"winter-test/service"
)

func ProductsRoute(r *gin.Engine) {
	// 商品路由
	p := r.Group("/suning/products")
	{
		//主页的一些设置
		p.GET("/search", SearchProducts)       //主页搜索
		p.GET("/category/:name", ShowCategory) //主页分类展示
		p.GET("/rotation", ShowRotation)       //轮播页面展示
		// 商品的详情页
		p.GET("/product/:id", DetailedProduct) //商品详情
		//商品收藏
		r.POST("/favorites", service.JwtAuthMiddleware(), AddFavorite)
		//r.DELETE("/favorites/:product_id",service.JwtAuthMiddleware(), RemoveFavorite)
		//r.GET("/favorites/:user_id",  service.JwtAuthMiddleware(),ShowFavorites)

	}

}

// 主页根据排序规则搜索商品
func SearchProducts(c *gin.Context) {
	query := c.Query("query")
	sort := c.Query("sort")
	order := c.Query("order")
	products, err := dao.SearchProducts(query, sort, order)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"info":   "error",
			"data": gin.H{
				"error": err,
			},
		})
	} else {
		c.JSON(200, gin.H{
			"status": 200,
			"info":   "success",
			"data":   products,
		})
	}
}

// 分类展示商品
func ShowCategory(c *gin.Context) {
	category := c.Param("name")
	products, err := dao.ShowCategory(category)
	if err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"info":   "error",
			"data": gin.H{
				"error": err,
			},
		})
	} else {
		c.JSON(200, gin.H{
			"status": 200,
			"info":   "success",
			"data":   products,
		})
	}
}

// 展示轮播商品
func ShowRotation(c *gin.Context) {
	rotationProducts, err := dao.ShowRotation()
	if err != nil {
		c.JSON(500, gin.H{
			"status": 500,
			"info":   "error",
			"data": gin.H{
				"error": err,
			},
		})
		return
	}
	c.JSON(200, gin.H{
		"status": 200,
		"info":   "success",
		"data":   rotationProducts,
	})

}

// 商品详情页面
func DetailedProduct(c *gin.Context) {
	// 获取商品id
	id := c.Param("id")
	// 获取商品详细信息
	service.Productdata(c, id)
	// 获取该商品所有款式信息
	service.GetStyles(c, id)
}
func AddFavorite(c *gin.Context) {

}
