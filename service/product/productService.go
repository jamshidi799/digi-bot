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

func AddProductToDB(senderId int, url string) (model.Product, error) {
	//fmt.Printf("%+v %+v\n", senderId, url)
	if res := strings.Contains(url, "digikala.com"); !res {
		return model.Product{}, errors.New("ادرس نامعتبر است")
	}

	product, err := crawler.Crawl(url)
	if err != nil {
		return model.Product{}, err
	}

	//fmt.Printf("%+v", product)
	productModel := product.ToProductModel()

	result := db.DB.Where(productModel).Find(&productModel)
	if result.RowsAffected == 0 {
		productModel = product.ToProductModel()
		db.DB.Create(&productModel)
	}

	pivot := model.PivotModel{UserId: senderId, ProductId: productModel.ID}
	db.DB.Create(&pivot)

	log.Printf("new product added: %s\n", product.Name)
	return product, nil
}

func UpdateProduct(product model.ProductModel, newProduct model.Product) {
	product.Price = newProduct.Price
	db.DB.Save(&product)
}

func DeleteAllUserProduct(userId int) {
	db.DB.Where("user_id = ?", userId).Delete(&model.PivotModel{})
}

func DeleteProductByName(name string, userId int) string {
	var ids []int
	name = strings.TrimSpace(name)

	db.DB.
		Select("pivot_models.id").
		Model(&model.PivotModel{}).
		Joins("JOIN product_models product on product.id = pivot_models.product_id").
		Where("product.name = ? AND pivot_models.user_id = ?", name, userId).
		Find(&ids)

	if len(ids) == 0 {
		return "پروداکت یافت نشد!"
	}

	var deletedPivot model.PivotModel
	db.DB.
		Model(&model.PivotModel{}).
		Where("id IN ?", ids).
		First(&deletedPivot).
		Delete(&model.PivotModel{})

	product := db.GetProductById(deletedPivot.ProductId)

	msg := messageCreator.CreateDeleteProductSuccessfulMsg(product.ToProduct())

	return msg
}

func GetProductList(userId int) string {
	products := db.GetAllProductByUserId(userId)
	return messageCreator.CreateProductListMsg(products)
}
