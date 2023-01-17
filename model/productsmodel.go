package model

import "database/sql"

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

type SmallCart struct {
	ID        int    `json:"id" form:"id"`
	UserName  string `json:"user_name" form:"user_name"`
	ProductID int    `json:"product_id" form:"product_id"`
	Quantity  int    `json:"quantity" form:"quantity"`
	StyleID   int    `json:"style_id" form:"style_id"`
}

// 商品购物车
type Cart struct {
	ID           int            `json:"id" form:"id"`
	UserName     string         `json:"user_name" form:"user_name"`
	ProductName  string         `json:"product_name" form:"product_name"`
	Quantity     int            `json:"quantity" form:"quantity"`
	ProductImage sql.NullString `json:"product_image" form:"product_image"`
	ProductPrice float64        `json:"product_price" form:"product_price"`
	StyleName    string         `json:"style_name" form:"style_name"`
}
