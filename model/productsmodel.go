package model

import (
	"database/sql"
)

// 商品
type Product struct {
	ID          int            `json:"id" form:"id"`
	Name        string         `json:"name" form:"name"`
	Description sql.NullString `json:"description" form:"description"`
	Image       sql.NullString `json:"image" form:"image"`
	Category    string         `json:"category" form:"category"`
	Price       float64        `json:"price" form:"price"`
	Stock       int            `json:"stock" form:"stock"`
	Sale        int            `json:"sale" form:"sale"`
	Rating      float64        `json:"rating" form:"rating"`
	Seller      string         `json:"seller" form:"seller"`
}

// 轮播商品
type RotationProduct struct {
	ID          int            `json:"id" form:"id"`
	Name        string         `json:"name" form:"name"`
	Description sql.NullString `json:"description" form:"description"`
	Image       sql.NullString `json:"image" form:"image"`
	URL         sql.NullString `json:"url" form:"url"` //商品跳转地址
}

// 商品款式
type Style struct {
	ID         string         `json:"id" form:"id"`
	Name       string         `json:"name" form:"name"`
	ProductID  string         `json:"product_id" form:"product_id"`
	Stock      int            `json:"stock" form:"stock"`
	StyleImage sql.NullString `json:"style_image" form:"style_image"`
}

// 收藏商品
type Favorite struct {
	UserName  string `json:"user_name" form:"user_name"`
	ProductID int    `json:"product_id" form:"product_id"`
}

// 店铺
type Seller struct {
	ID           int     `json:"id" form:"id"`
	SellerName   string  `json:"seller_name" form:"seller_name"`
	Announcement string  `json:"announcement" form:"announcement"`
	Description  string  `json:"description" form:"description"`
	SellerImage  string  `json:"seller_image" form:"seller_image"`
	SellerGrade  float64 `json:"seller_grade" form:"seller_grade"`
}
