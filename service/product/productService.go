package product

import (
	"digi-bot/crawler"
	"digi-bot/db"
	"digi-bot/model"
	"errors"
	"fmt"
	"strings"
)

func AddProductToDB(senderId int, url string) (model.Product, error) {
	fmt.Printf("%+v %+v", senderId, url)
	if res := strings.Contains(url, "digikala.com"); !res {
		return model.Product{}, errors.New("ادرس نامعتبر است")
	}

	product, err := crawler.Crawl(url)
	if err != nil {
		return model.Product{}, err
	}

	fmt.Printf("%+v", product)
	productModel := product.ToProductModel(senderId)
	db.DB.Create(&productModel)
	return productModel.ToProduct(), nil
}

func UpdateProduct(product model.ProductModel, newProduct model.Product) {
	product.Price = newProduct.Price
	product.OldPrice = newProduct.OldPrice
	db.DB.Save(&product)
}
func DeleteAllUserProduct(userId int) {
	db.DB.Where("user_id = ?", userId).Delete(&model.ProductModel{})
}

func DeleteProductByName(name string) model.Product {
	var product model.ProductModel
	name = strings.TrimSpace(name)
	db.DB.Where("name = ?", name).First(&product).Delete(&model.ProductModel{})
	return product.ToProduct()
}
