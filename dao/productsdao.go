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
	err := DB.Where("name LIKE ?", "%"+q+"%").Order(order).Find(&products).Error
	if err != nil {
		fmt.Println("***********", err)
		return nil, err
	}
	for i := range products {
		var seller model.Seller
		result := DB.Model(&model.Seller{}).Where("name = ?", products[i].Seller).First(&seller)
		if result.Error != nil {
			return nil, result.Error
		}
		if err := DB.Where("id = ?", seller.ID).First(&seller).Error; err != nil {
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
	result := DB.Where("category = ?", c).Find(&products)
	if result.Error != nil {
		fmt.Println("***********", result.Error)
		return nil, result.Error
	}
	for i := range products {
		var seller model.Seller
		result = DB.Model(&model.Seller{}).Where("name = ?", products[i].Seller).First(&seller)
		if result.Error != nil {
			return nil, result.Error
		}
		result := DB.Where("id = ?", seller.ID).First(&seller)
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
	if err := DB.Find(&rotationProducts).Error; err != nil {
		fmt.Println("***********", err)
		return nil, err
	}
	return rotationProducts, nil
}

func Productdata(id string) ([]model.Product, error) {
	var products []model.Product
	result := DB.Where("id = ?", id).Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	for i := range products {
		var seller model.Seller
		DB.Model(&products[i]).Association("Seller").Find(&seller)
		products[i].Seller = seller.SellerName
	}
	return products, nil
}

// 商品款式
func GetStyles(product_id string) ([]model.Style, error) {
	var styles []model.Style
	result := DB.Preload("product").Where("product_id = ?", product_id).Find(&styles)
	if result.Error != nil {
		return nil, result.Error
	}
	return styles, nil
}

// 收藏商品
func AddFavorite(f *model.Favorite) error {
	// 判断是否已经收藏过该商品
	var count int64
	if err := DB.Model(&model.Favorite{}).Where("user_name = ? and product_id = ?", f.UserName, f.ProductID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("该商品已经收藏了")
	}

	// 插入收藏商品记录
	if err := DB.Create(f).Error; err != nil {
		return err
	}
	return nil
}

// 取消收藏的商品
func RemoveFavorite(f *model.Favorite) error {
	// 判断是否已经收藏过该商品
	var count int64
	if err := DB.Model(&model.Favorite{}).Where("user_name = ? and product_id = ?", f.UserName, f.ProductID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("还未收藏该商品，无法取消收藏商品")
	}
	// 删除收藏商品记录
	if err := DB.Where("user_name = ? and product_id = ?", f.UserName, f.ProductID).Delete(&model.Favorite{}).Error; err != nil {
		return err
	}
	return nil
}

// 展示用户所有的收藏商品
func ShowFavorites(username string) ([]model.Product, error) {
	var products []model.Product
	if err := DB.Model(&model.Favorite{}).Joins("JOIN product ON favorite.product_id = product.id").
		Where("favorite.user_name = ?", username).Scan(&products).Error; err != nil {
		return nil, err
	}
	for i := range products {
		// 获取卖家信息
		var seller model.Seller

		result := DB.Model(&model.Seller{}).Where("name = ?", products[i].Seller).First(&seller)
		if result.Error != nil {
			return nil, result.Error
		}
		if err := DB.Model(&model.Seller{}).Where("id = ?", seller.ID).First(&seller).Error; err != nil {
			return nil, err
		}
		products[i].Seller = seller.SellerName
	}
	return products, nil
}

// 展示店铺
func ShowSeller(sellerID int) (model.Seller, error) {
	var seller model.Seller
	err := DB.First(&seller, sellerID).Error
	if err != nil {
		return seller, err
	}
	return seller, nil
}

// 店铺根据排序对象进行展示商品
func SortProducts(orderBy, sort string) ([]model.Product, error) {
	var products []model.Product
	query := DB.Order(orderBy + " " + sort).Find(&products)
	if query.Error != nil {
		return nil, query.Error
	}
	for i := range products {
		DB.Model(&products[i]).Association("Seller").Find(&products[i].Seller)
	}
	return products, nil
}

// 店铺分类展示商品
func SellerShowCategory(c string, sellerID int) ([]model.Product, error) {
	var products []model.Product
	query := DB.Where("category = ? AND seller_id = ?", c, sellerID).Find(&products)
	if query.Error != nil {
		return nil, query.Error
	}
	for i := range products {
		DB.Model(&products[i]).Association("Seller").Find(&products[i].Seller)
	}
	return products, nil
}
