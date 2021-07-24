package product

import (
	"digi-bot/crawler"
	"digi-bot/db"
	"digi-bot/model"
	"fmt"
	"strings"
)

func AddProductToDB(senderId int, url string) (model.Product, error) {
	fmt.Printf("%+v %+v", senderId, url)
	if res := strings.Contains(url, "digikala.com"); !res {
		return model.Product{}, nil
	}

	product, err := crawler.Crawl(url)
	if err != nil {
		return model.Product{}, nil
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
func DeleteAllUserProduct(userId int) {}
func DeleteProductByName(name string) {}
