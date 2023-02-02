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
		p.GET("/search", SearchProducts)       // 主页搜索
		p.GET("/category/:name", ShowCategory) // 主页分类展示
		p.GET("/rotation", ShowRotation)       // 轮播页面展示
		// 商品的详情页
		p.GET("/:product_id", DetailedProduct) // 商品详情
		//商品收藏
		p.GET("/favorites/:product_id", service.JwtAuthMiddleware(), AddFavorite)       // 收藏商品
		p.DELETE("/favorites/:product_id", service.JwtAuthMiddleware(), RemoveFavorite) // 取消收藏商品
		// 店铺
		p.GET("/seller/show", ShowSeller)             // 店铺详情展示
		p.GET("/seller/sort", SortProducts)           // 店铺商品排序展示
		p.GET("/seller/category", SellerShowCategory) // 店铺商品分类展示
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

// 店铺展示
func ShowSeller(c *gin.Context) {
	// 请求店铺id
	id := c.Query("seller_id")
	sellerID, _ := strconv.Atoi(id)
	// 调用dao层获取店铺信息
	seller, err := dao.ShowSeller(sellerID)
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
		"data":   seller,
	})
}

// 店铺商品按照销量...排序
func SortProducts(c *gin.Context) {
	// 获取排序规则
	orderBy := c.DefaultQuery("order_by", "")
	sort := c.DefaultQuery("sort", "asc")

	// 调用dao层获取商品信息
	products, err := dao.SortProducts(orderBy, sort)
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
		"data":   products,
	})
}

// 店铺商品分类展示
func SellerShowCategory(c *gin.Context) {
	category := c.Query("name")
	id := c.Query("seller_id")
	sellerID, _ := strconv.Atoi(id)
	products, err := dao.SellerShowCategory(category, sellerID)
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
