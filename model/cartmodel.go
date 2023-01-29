package model

import "database/sql"

type SmallCart struct {
	ID        int    `json:"id" form:"id"`
	UserName  string `json:"user_name" form:"user_name"`
	ProductID int    `json:"product_id" form:"product_id"`
	Quantity  int    `json:"quantity" form:"quantity"`
	StyleID   int    `json:"style_id" form:"style_id"`
}

// 商品购物车
type Cart struct {
	ID           int64          `json:"id" form:"id"`
	UserName     string         `json:"user_name" form:"user_name"`
	ProductName  string         `json:"product_name" form:"product_name"`
	Quantity     int            `json:"quantity" form:"quantity"`
	ProductImage sql.NullString `json:"product_image" form:"product_image"`
	ProductPrice float64        `json:"product_price" form:"product_price"`
	StyleName    string         `json:"style_name" form:"style_name"`
}

// 结算购物车传入的ids
type SettleCartRequest struct {
	IDs []int `json:"ids" form:"ids"`
}
