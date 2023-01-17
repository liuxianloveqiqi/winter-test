package dao

import (
	"errors"
	"winter-test/model"
)

// 添加商品到购物车
func AddCart(c *model.SmallCart) error {
	// 查询是否已经加入购物车
	var count int
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
	} else {
		// 插入新的购物车记录
		_, err = db.Exec("insert into cart (user_name, product_id,style_id, quantity) values (?, ?, ?,?)", c.UserName, c.ProductID, c.StyleID, c.Quantity)
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
	err := db.QueryRow("select count(*),quantity from cart where user_name = ? and product_id = ? and style_id=? group by id", c.UserName, c.ProductID, c.StyleID).Scan(&count, &quantity)
	if err != nil {
		return err
	} else if quantity == 0 {
		return errors.New("该商品没有添加到购物车中，无法删除")
	} else if quantity < c.Quantity {
		return errors.New("删除购物车中的数量超出，无法删除")
	} else {
		_, err = db.Exec("update cart set quantity = quantity - ? where user_name = ? and product_id = ? and style_id=?", c.Quantity, c.UserName, c.ProductID, c.StyleID)
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
                                 product.price, style.name as style_name
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
