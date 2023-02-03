package dao

import (
	"errors"
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
	// 创建了一个事务
	tx, err := db.Begin()
	if err != nil {
		return order, err
	}
	// 确保保事务总是被提交或回滚
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			// 没有错误发生，提交事务
			err = tx.Commit()
		}
	}()

	// 插入订单数据
	res, err := tx.Exec("insert into orders (user_name, amount, step, created_time, updated_time, address_id) values (?, ?, ?, ?, ?, ?)", userName, amount, 1, time.Now(), time.Now(), address_id)
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
		_, err = tx.Exec("insert into order_items (cart_id, order_id) values (?, ?)", cart.ID, orderID)
		if err != nil {
			fmt.Println(err, "3333333333333")

			return order, err
		}
	}
	// 从用户账户中扣除总价
	_, err = tx.Exec("update user set money = money - ? where username = ?", amount, userName)
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
	_, err = tx.Exec(fmt.Sprintf("delete  from order_items where cart_id in (%s)", idList))
	if err != nil {
		fmt.Println("7777")
		fmt.Println(err)
		return order, err
	}
	_, err = tx.Exec(fmt.Sprintf("delete  from cart where id in (%s)", idList))
	if err != nil {
		fmt.Println("*******", cartItemIDs)
		fmt.Println(idList, "#####")
		fmt.Println(err, "55555555555")

		return order, err
	}
	// 查询出订单信息
	row := tx.QueryRow(`select orders.id ,orders.user_name, orders.amount,
		                          user.human_name, user.phone_number, address.place, orders.step,
		                          orders.created_time, orders.updated_time
		                          from orders
		                          join address on orders.address_id=address.id
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

// 分类展示订单
func ShowOrdersByStep(username string, step int) ([]model.Order, error) {
	rows, err := db.Query(`select orders.id ,orders.user_name, orders.amount,
		                          user.human_name, user.phone_number, address.place, orders.step,
		                          orders.created_time, orders.updated_time
		                          from orders
		                          join address on orders.address_id=address.id
		                          join user on orders.user_name = user.username
		                          where orders.user_name= ? and orders.step = ?`, username, step)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var order model.Order
		err = rows.Scan(&order.ID, &order.UserName, &order.Amount, &order.HumanName, &order.PhoneNumber, &order.Address, &order.Step, &order.CreatedTime, &order.UpdatedTime)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	if len(orders) == 0 {
		return nil, errors.New("没有订单")
	}
	fmt.Println(orders)
	return orders, nil
}

// 改变订单状态
func UpdateOrderStep(id int, un string, step int) error {
	query := `update orders set step = ? where id = ? and user_name = ?`
	_, err := db.Exec(query, step, id, un)
	if err != nil {
		return err
	}
	return nil
}
