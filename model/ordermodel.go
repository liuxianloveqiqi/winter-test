package model

// 订单
type Order struct {
	ID          int     `json:"id" form:"id"`
	UserName    string  `json:"user_name" form:"user_name"`
	Amount      float64 `json:"amount" form:"amount"`
	HumanName   string  `json:"human_name" form:"human_name"`
	PhoneNumber int     `json:"phone_number" form:"phone_number"`
	Address     string  `json:"address" form:"address"`
	Step        int64   `json:"step" form:"step"` //1为确认订单，2未收货，3已收货
	CreatedTime string  `json:"created_time" form:"created_time"`
	UpdatedTime string  `json:"updated_time" form:"updated_time"`
}

// 订单商品关联
type OrderItem struct {
	ID      int `json:"id" form:"id"`
	CartId  int `json:"cart_id" form:"cart_id"`
	OrderId int `json:"order_id" form:"order_id"`
}
