package dao

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"winter-test/model"
)

// 添加商品到购物车
func AddCart(c *model.SmallCart) error {
	// 查询是否已经加入购物车
	var count int
	var simplePrice float64
	err := db.QueryRow("select count(*) from cart where user_name = ? and product_id = ? and style_id=?", c.UserName, c.ProductID, c.StyleID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		// 更新购物车记录中的quantity
		_, err = db.Exec("update cart set quantity = quantity + ? where user_name = ? and product_id = ? and style_id=?", c.Quantity, c.UserName, c.ProductID, c.StyleID)
		if err != nil {
			return err
		}
		// 查询出商品单价
		if err1 := db.QueryRow("select price from product where id=?", c.ProductID).Scan(&simplePrice); err1 != nil {
			return err1
		}
		// 更新购物车的某商品总价钱=单价*数量
		_, err = db.Exec("update cart set price= price + ? where user_name = ? and product_id = ? and style_id=?", float64(c.Quantity)*simplePrice, c.UserName, c.ProductID, c.StyleID)
		if err != nil {
			return err
		}
	} else {
		// 插入新的购物车记录
		_, err = db.Exec("insert into cart (user_name, product_id,style_id, quantity,price) values (?, ?, ?,?,?)", c.UserName, c.ProductID, c.StyleID, c.Quantity, float64(c.Quantity)*simplePrice)
		if err != nil {
			return err
		}
	}
	return nil
}

// 删除购物车中的商品
func RemoveCart(c *model.SmallCart) error {
	// 判断该用户是否已经添加过该商品
	var count, quantity int
	var simplePrice float64
	// 查询出商品单价
	if err1 := db.QueryRow("select price from product where id=?", c.ProductID).Scan(&simplePrice); err1 != nil {
		return err1
	}
	err := db.QueryRow("select count(*),quantity from cart where user_name = ? and product_id = ? and style_id=? group by id", c.UserName, c.ProductID, c.StyleID).Scan(&count, &quantity)
	if err != nil {
		return err
	} else if quantity == 0 {
		return errors.New("该商品没有添加到购物车中，无法删除")
	} else if quantity < c.Quantity {
		return errors.New("删除购物车中的数量超出，无法删除")
	} else {
		_, err = db.Exec("update cart set quantity = quantity - ?, price=price-? where user_name = ? and product_id = ? and style_id=?", c.Quantity, float64(c.Quantity)*simplePrice, c.UserName, c.ProductID, c.StyleID)
		if err != nil {
			return err
		}
	}
	return nil
}

// 展示用户的购物车
func ListCart(username string) ([]*model.Cart, error) {
	// 查询购物车中所有商品信息
	rows, err := db.Query(`select cart.id, cart.quantity, product.name, style.style_image,
                                 cart.price, style.name as style_name
                                 from cart
                                 join product on cart.product_id = product.ID
                                 join style on cart.style_id = style.id
                                 where cart.user_name = ?`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// 创建结果切片
	var carts []*model.Cart
	for rows.Next() {
		var cart model.Cart
		if err := rows.Scan(&cart.ID, &cart.Quantity, &cart.ProductName, &cart.ProductImage, &cart.ProductPrice, &cart.StyleName); err != nil {
			return nil, err
		}
		carts = append(carts, &cart)
	}
	return carts, nil
}

// 结算购物车中的部分商品
func SettleCart(cartItemIDs []int) ([]*model.Cart, float64, error) {
	var carts []*model.Cart
	var cart model.Cart
	var totalPrice float64
	var s []string
	for _, id := range cartItemIDs {
		s = append(s, strconv.Itoa(id))
	}
	idList := strings.Join(s, ",")

	query := fmt.Sprintf("select c.id, c.user_name, p.name, c.quantity, p.image, p.price, s.name FROM cart c LEFT JOIN product p ON c.product_id = p.ID LEFT JOIN style s ON c.style_id = s.id WHERE c.id IN (%s)", idList)
	rows, err := db.Query(query)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&cart.ID, &cart.UserName, &cart.ProductName, &cart.Quantity, &cart.ProductImage, &cart.ProductPrice, &cart.StyleName)
		if err != nil {
			return nil, 0, err
		}
		carts = append(carts, &cart)
		// 计算总价
		totalPrice += cart.ProductPrice
	}
	//// 从用户账户中扣除总价
	//_, err = db.Exec("update user set money = money - ? where username = ?", totalPrice, cart.UserName)
	//if err != nil {
	//	return nil, 0, err
	//}
	// 从购物车中删除结算的商品
	//_, err = db.Exec(fmt.Sprintf("delete  from cart where id in (%s)", idList))
	//if err != nil {
	//	return nil, 0, err
	//}
	return carts, totalPrice, nil
}
