package api

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"winter-test/dao"
	"winter-test/model"
	"winter-test/service"
)

func CommentRoute(r *gin.Engine) {
	comment := r.Group("suning/comment")
	{
		comment.GET("/create", service.JwtAuthMiddleware(), CreateComment) // 创造一条评论
		comment.GET("/show", GetComment)                                   // 展示商品评论
		comment.DELETE("/remove/:id", service.JwtAuthMiddleware(), DeleteComment)
	}
}

// 用户新建评论
func CreateComment(c *gin.Context) {
	// 获取当前登录用户的用户名
	username = c.MustGet("claims").(*model.MyClaims).UserName
	// 获取请求参数中的评论信息
	comment := c.Query("comment")
	id := c.Query("product_id")
	productID, _ := strconv.Atoi(id)
	// 不请求parent_id默认为0，代表无父评论
	id2 := c.DefaultQuery("parent_id", "0")
	parentID, _ := strconv.Atoi(id2)
	service.CreateComment(c, username, comment, productID, parentID)

}

// 展示商品评论
func GetComment(c *gin.Context) {
	// 获取请求参数中的产品 ID
	id := c.Query("product_id")
	productID, _ := strconv.Atoi(id)
	// 不请求parent_id默认为0，代表无父评论
	id2 := c.DefaultQuery("parent_id", "0")
	parentID, _ := strconv.Atoi(id2)
	// 从数据库中获取评论信息
	comments, err := dao.GetComment(productID, parentID)
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
		"data": gin.H{
			"comments": comments,
		},
	})
}

// 删除评论
func DeleteComment(c *gin.Context) {
	// 获取当前登录用户的用户名
	username := c.MustGet("claims").(*model.MyClaims).UserName
	// 获取请求参数中的评论ID
	commentIDStr := c.Param("id")
	commentID, _ := strconv.Atoi(commentIDStr)

	// 检查评论是否属于当前用户
	un, err := dao.GetCommentByID(commentID)
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
	if username != un {
		c.JSON(403, gin.H{
			"status": 403,
			"info":   "error",
			"data": gin.H{
				"error": "没权限",
			},
		})
		return
	}

	// 删除评论
	err = dao.DeleteComment(commentID)
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
		"data":   "删除成功",
	})
}
