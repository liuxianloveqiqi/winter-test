package dao

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"winter-test/model"
)

// 展示收货信息
func ShowMessage(username string) ([]model.Address, error) {
	rows, err := db.Query("select * from address where user_name = ?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var addresses []model.Address
	for rows.Next() {
		var address model.Address
		if err := rows.Scan(&address.ID, &address.UserName, &address.Place); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}
	return addresses, nil
}

// 创建订单
func CreateOrder(carts []*model.Cart, amount float64, userName string, address_id int, cartItemIDs []int) (model.Order, error) {

	var order model.Order
	res, err := db.Exec("insert into orders (user_name, amount, step, created_time, updated_time, address_id) values (?, ?, ?, ?, ?, ?)", userName, amount, 1, time.Now(), time.Now(), address_id)
	if err != nil {
		println(address_id)
		fmt.Println(err, "1111111111")
		return order, err
	}
	orderID, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err, "22222222222")

		return order, err
	}

	// 创建订单商品关联
	for _, cart := range carts {
		_, err = db.Exec("insert into order_items (cart_id, order_id) values (?, ?)", cart.ID, orderID)
		if err != nil {
			fmt.Println(err, "3333333333333")

			return order, err
		}
	}
	// 从用户账户中扣除总价
	_, err = db.Exec("update user set money = money - ? where username = ?", amount, userName)
	if err != nil {
		fmt.Println(err, "444444444444444")

		return order, err
	}
	// 扣钱完成从购物车中删除结算的商品
	var s []string
	for _, id := range cartItemIDs {
		s = append(s, strconv.Itoa(id))
	}
	idList := strings.Join(s, ",")
	// 先删掉订单关联表中的商品删除
	_, err = db.Exec(fmt.Sprintf("delete  from order_items where cart_id in (%s)", idList))
	if err != nil {
		fmt.Println("7777")
		return order, err
	}
	_, err = db.Exec(fmt.Sprintf("delete  from cart where id in (%s)", idList))
	if err != nil {
		fmt.Println("*******", cartItemIDs)
		fmt.Println(idList, "#####")
		fmt.Println(err, "55555555555")

		return order, err
	}
	row := db.QueryRow(`select orders.id ,orders.user_name, orders.amount, 
                               user.human_name, user.phone_number, address.place, orders.step,
                               orders.created_time, orders.updated_time 
                               from orders  
                               join address on orders.address_id
                               join user on orders.user_name = user.username
                               where orders.id = ?`, orderID)
	fmt.Println(row, "uuuu")
	err = row.Scan(&order.ID, &order.UserName, &order.Amount, &order.HumanName, &order.PhoneNumber, &order.Address, &order.Step, &order.CreatedTime, &order.UpdatedTime)
	if err != nil {
		fmt.Println(err, "66666666666666")
		fmt.Println(orderID, "KKKKKK")
		return order, err
	}

	return order, nil
}
