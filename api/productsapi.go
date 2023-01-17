package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"winter-test/dao"
	"winter-test/model"
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
		p.GET("/:product_id", DetailedProduct) //商品详情
		//商品收藏
		p.GET("/favorites/:product_id", service.JwtAuthMiddleware(), AddFavorite)       //收藏商品
		p.DELETE("/favorites/:product_id", service.JwtAuthMiddleware(), RemoveFavorite) //取消收藏商品
	}
	// 实现购物车功能
	cart := r.Group("/suning/cart")
	{
		cart.POST("/add", service.JwtAuthMiddleware(), AddCart) //添加商品到购物车
		cart.GET("/list", service.JwtAuthMiddleware(), ListCart)
		cart.DELETE("/remove", service.JwtAuthMiddleware(), RemoveCart) //删除购物车中的商品
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
	id := c.Param("product_id")
	// 获取商品详细信息
	service.Productdata(c, id)
	// 获取该商品所有款式信息
	service.GetStyles(c, id)
}

// 收藏商品
func AddFavorite(c *gin.Context) {
	var favorite model.Favorite
	// 获取商品id
	id := c.Param("product_id")
	favorite.ProductID, _ = strconv.Atoi(id)
	// 获取token里面的username
	favorite.UserName = c.MustGet("claims").(*model.MyClaims).UserName
	fmt.Println(favorite.UserName)
	if err := dao.AddFavorite(&favorite); err != nil {
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
		"data":   "收藏商品成功",
	})
}

// 删除收藏商品
func RemoveFavorite(c *gin.Context) {
	var favorite model.Favorite
	// 获取商品id
	id := c.Param("product_id")
	favorite.ProductID, _ = strconv.Atoi(id)
	// 获取token里面的username
	favorite.UserName = c.MustGet("claims").(*model.MyClaims).UserName
	fmt.Println(favorite.UserName)
	if err := dao.RemoveFavorite(&favorite); err != nil {
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
		"data":   "取消收藏商品成功",
	})
}

// 展示用户的全部收藏
func ShowFavorites(c *gin.Context) {
	// 获取token里面的username
	username1 := c.MustGet("claims").(*model.MyClaims).UserName
	favorites, err := dao.ShowFavorites(username1)
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
		"data":   favorites,
	})
}
