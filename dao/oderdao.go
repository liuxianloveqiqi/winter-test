package dao

import (
	"time"
	"winter-test/model"
)

func CreateOrder(carts []*model.Cart, totalPrice float64, userName string) (int64, error) {
	// 创建订单
	res, err := db.Exec("insert into orders (user_name, amount, step, created_time, updated_time) values (?, ?, ?, ?, ?)", userName, totalPrice, 1, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}
	orderID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// 创建订单商品关联
	for _, cart := range carts {
		_, err = db.Exec("insert into order_items (cart_id, order_id) values (?, ?)", cart.ID, orderID)
		if err != nil {
			return 0, err
		}
	}

	return orderID, nil
}
