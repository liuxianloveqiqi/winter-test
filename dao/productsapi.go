package dao

import (
	"database/sql"
	"fmt"
	"winter-test/model"
)

// 主页根据排序规则搜索商品
func SearchProducts(q, s, o string) ([]model.Product, error) {
	var rows *sql.Rows
	var err error
	if o == "asc" || o == "" {
		rows, err = db.Query("select * from product where name like ? order by ?", "%"+q+"%", s)
	} else {
		rows, err = db.Query("select * from product where name like ? order by ? desc", "%"+q+"%", s)

	}
	if err != nil {
		fmt.Println("***********", err)
		return nil, err
	}
	defer rows.Close()
	var products []model.Product
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Image, &product.Category, &product.Price, &product.Stock, &product.Sale, &product.Rating, &product.Seller); err != nil {
			fmt.Println("***********", err)
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil

}
func ShowCategory(c string) ([]model.Product, error) {
	rows, err := db.Query("select * from product where category = ?", c)
	if err != nil {
		fmt.Println("***********", err)

		return nil, err
	}
	defer rows.Close()
	var products []model.Product
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Image, &product.Category, &product.Price, &product.Stock, &product.Sale, &product.Rating, &product.Seller); err != nil {
			fmt.Println("***********", err)

			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

// 轮播展示商品
func ShowRotation() ([]model.RotationProduct, error) {
	rows, err := db.Query("select * from rotation_product")
	if err != nil {
		fmt.Println("***********", err)

		return nil, err
	}
	defer rows.Close()
	var rotationProducts []model.RotationProduct

	for rows.Next() {
		var rotationProduct model.RotationProduct
		if err := rows.Scan(&rotationProduct.ID, &rotationProduct.Name, &rotationProduct.Image, &rotationProduct.Description, &rotationProduct.URL); err != nil {
			fmt.Println("***********", err)

			return nil, err
		}
		rotationProducts = append(rotationProducts, rotationProduct)
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
	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Image, &product.Category, &product.Price, &product.Stock, &product.Sale, &product.Rating, &product.Seller); err != nil {
			fmt.Println("***********", err)

			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

// 商品款式
func GetStyles(product_id string) ([]model.Style, error) {
	var styles []model.Style
	rows, err := db.Query("select  * from style where product_id = ?", product_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var style model.Style
		if err := rows.Scan(&style.ID, &style.Name, style.ProductID, &style.Stock); err != nil {
			return nil, err
		}
		styles = append(styles, style)
	}
	return styles, nil
}
