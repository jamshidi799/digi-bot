package product

import (
	"digi-bot/crawler"
	"digi-bot/db"
	"digi-bot/messageCreator"
	"digi-bot/model"
	"errors"
	"log"
	"strings"
)

func AddProductToDB(senderId int, url string) (model.ProductDto, error) {
	//fmt.Printf("%+v %+v\n", senderId, url)
	if res := strings.Contains(url, "digikala.com"); !res {
		return model.ProductDto{}, errors.New("ادرس نامعتبر است")
	}

	product, err := crawler.Crawl(url)
	if err != nil {
		return model.ProductDto{}, err
	}

	//fmt.Printf("%+v", product)
	productModel := product.ToProduct()

	result := db.DB.Where(productModel).Find(&productModel)
	if result.RowsAffected == 0 {
		productModel = product.ToProduct()
		db.DB.Create(&productModel)
	}

	pivot := model.Pivot{UserID: senderId, ProductID: productModel.ID}
	db.DB.Create(&pivot)

	log.Printf("new product added: %s\n", product.Name)
	return product, nil
}

func UpdateProduct(product model.Product, newProduct model.ProductDto) {
	product.Price = newProduct.Price
	db.DB.Save(&product)
}

func DeleteAllUserProduct(userId int) {
	db.DB.Where("user_id = ?", userId).Delete(&model.Pivot{})
}

func DeleteProductByName(name string, userId int) string {
	var ids []int
	name = strings.TrimSpace(name)

	db.DB.
		Select("pivots.id").
		Model(&model.Product{}).
		Joins("JOIN products product on product.id = pivots.product_id").
		Where("product.name = ? AND pivots.user_id = ?", name, userId).
		Find(&ids)

	if len(ids) == 0 {
		return "پروداکت یافت نشد!"
	}

	var deletedPivot model.Pivot
	db.DB.
		Model(&model.Pivot{}).
		Where("id IN ?", ids).
		First(&deletedPivot).
		Delete(&model.Pivot{})

	product := db.GetProductById(deletedPivot.ProductID)

	msg := messageCreator.CreateDeleteProductSuccessfulMsg(product.ToDto())

	return msg
}

func GetProductList(userId int) string {
	products := db.GetAllProductByUserId(userId)
	return messageCreator.CreateProductListMsg(products)
}
