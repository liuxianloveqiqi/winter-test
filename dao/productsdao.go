package dao

import (
	"fmt"
	"winter-test/model"
)

// 主页根据排序规则搜索商品
func SearchProducts(q, s, o string) ([]model.Product, error) {
	var products []model.Product
	order := fmt.Sprintf("%s %s", s, o)
	if o == "" {
		order = s
	}
	err := db.Where("name LIKE ?", "%"+q+"%").Order(order).Find(&products).Error
	if err != nil {
		fmt.Println("***********", err)
		return nil, err
	}
	for i := range products {
		var seller model.Seller
		result := db.Model(&model.Seller{}).Where("name = ?", products[i].Seller).First(&seller)
		if result.Error != nil {
			return nil, result.Error
		}
		if err := db.Where("id = ?", seller.ID).First(&seller).Error; err != nil {
			fmt.Println("***********", err)
			return nil, err
		}
		products[i].Seller = seller.SellerName
	}
	return products, nil
}

// 根据分类搜索
func ShowCategory(c string) ([]model.Product, error) {
	var products []model.Product
	result := db.Where("category = ?", c).Find(&products)
	if result.Error != nil {
		fmt.Println("***********", result.Error)
		return nil, result.Error
	}
	for i := range products {
		var seller model.Seller
		result = db.Model(&model.Seller{}).Where("name = ?", products[i].Seller).First(&seller)
		if result.Error != nil {
			return nil, result.Error
		}
		result := db.Where("id = ?", seller.ID).First(&seller)
		if result.Error != nil {
			fmt.Println("***********", result.Error)
			return nil, result.Error
		}
		products[i].Seller = seller.SellerName
	}
	return products, nil
}

// 轮播展示商品
func ShowRotation() ([]model.RotationProduct, error) {
	var rotationProducts []model.RotationProduct
	if err := db.Find(&rotationProducts).Error; err != nil {
		fmt.Println("***********", err)
		return nil, err
	}
	return rotationProducts, nil
}

// 商品详情
func Productdata(id string) ([]model.Product, error) {
	rows, err := db.Query("select * from product where ID = ?", id)
	if err != nil {
		fmt.Println("***********", err)

		return nil, err
	}
	defer rows.Close()
	var products []model.Product
	var sellerID int
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Image, &product.Category, &product.Price, &product.Stock, &product.Sale, &product.Rating, &sellerID); err != nil {
			fmt.Println("***********", err)

			return nil, err
		}
		db.QueryRow("select seller_name from seller where id = ?", sellerID).Scan(&product.Seller)
		fmt.Println("999999", product.Seller)
		products = append(products, product)
	}
	return products, nil
}

// 商品款式
func GetStyles(product_id string) ([]model.Style, error) {
	var styles []model.Style

	rows, err := db.Query("select  * from style where product_id = ?", product_id)
	if err != nil {
		fmt.Println("***********", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var style model.Style
		if err := rows.Scan(&style.ID, &style.Name, &style.ProductID, &style.Stock, &style.StyleImage); err != nil {
			return nil, err
		}
		fmt.Println("***********", err)

		styles = append(styles, style)
	}
	return styles, nil
}

// 收藏商品
func AddFavorite(f *model.Favorite) error {
	// 判断是否已经收藏过该商品
	var count int
	err := db.QueryRow("select count(*) from favorite where user_name = ? and product_id = ?", f.UserName, f.ProductID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("该商品已经收藏了")
	}
	// 插入收藏商品记录
	_, err = db.Exec("insert into favorite (user_name, product_id) values (?, ?)", f.UserName, f.ProductID)
	if err != nil {
		return err
	}
	return nil
}

// 取消收藏的商品
func RemoveFavorite(f *model.Favorite) error {
	// 判断是否已经收藏过该商品
	var count int
	err := db.QueryRow("select count(*) from favorite where user_name = ? and product_id = ?", f.UserName, f.ProductID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("还未收藏该商品，无法取消收藏商品")
	}
	// 删除收藏商品记录
	_, err = db.Exec("delete from favorite where user_name = ? and product_id = ?", f.UserName, f.ProductID)
	if err != nil {
		return err
	}
	return nil
}

// 展示用户所有的收藏商品
func ShowFavorites(username string) ([]model.Product, error) {
	// 内连接查询
	rows, err := db.Query("select product.* from favorite join product on favorite.product_id = product.id where favorite.user_name = ?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []model.Product
	var sellerID int
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Image, &product.Category, &product.Price, &product.Stock, &product.Sale, &product.Rating, &sellerID); err != nil {
			return nil, err
		}
		db.QueryRow("select seller_name from seller where id = ?", sellerID).Scan(&product.Seller)
		fmt.Println("999999", product.Seller)
		products = append(products, product)
	}
	return products, nil
}

// 展示店铺
func ShowSeller(sellerID int) (model.Seller, error) {
	var seller model.Seller
	err := db.QueryRow("select * from seller where id = ?", sellerID).Scan(&seller.ID, &seller.SellerName, &seller.Announcement, &seller.Description, &seller.SellerImage, &seller.SellerGrade)
	if err != nil {
		return seller, err
	}
	fmt.Println(seller, "ooooooooooo")
	return seller, nil
}

// 店铺根据排序对象进行展示商品
func SortProducts(orderBy, sort string) ([]model.Product, error) {
	var products []model.Product
	query := "select * from product"

	// 根据排序规则组装SQL查询语句
	if orderBy != "" {
		query = query + " order by " + orderBy + " " + sort
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var sellerID int
	for rows.Next() {
		var product model.Product

		if err = rows.Scan(&product.ID, &product.Name, &product.Description, &product.Image, &product.Category, &product.Price, &product.Stock, &product.Sale, &product.Rating, &sellerID); err != nil {
			return nil, err
		}
		db.QueryRow("select seller_name from seller where id = ?", sellerID).Scan(&product.Seller)

		products = append(products, product)
	}

	return products, nil
}

// 店铺分类展示商品
func SellerShowCategory(c string, SellerID int) ([]model.Product, error) {
	rows, err := db.Query("select * from product where category = ? and seller_id = ?", c, SellerID)
	if err != nil {
		fmt.Println("***********", err)

		return nil, err
	}
	defer rows.Close()
	var products []model.Product
	var sellerID int
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Image, &product.Category, &product.Price, &product.Stock, &product.Sale, &product.Rating, &sellerID); err != nil {
			fmt.Println("***********", err)

			return nil, err
		}
		db.QueryRow("select seller_name from seller where id = ?", sellerID).Scan(&product.Seller)
		fmt.Println("999999", product.Seller)

		products = append(products, product)
	}
	fmt.Println(products)
	return products, nil
}
