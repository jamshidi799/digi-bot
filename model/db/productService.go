package db

import (
	"digi-bot/model"
	"digi-bot/service"
	"digi-bot/utils"
	"errors"
	"log"
	"strconv"
	"time"
)

func AddProductToDB(product model.ProductDto, senderId int) (err error) {
	productModel := product.ToProduct()

	result := database.Where(productModel).Find(&productModel)
	if result.RowsAffected == 0 {
		productModel.RealID = product.Id
		database.Create(&productModel)
	}

	pivot := model.Pivot{UserID: senderId, ProductID: productModel.ID, NotificationSetting: 1}
	database.Create(&pivot)

	log.Printf("new product added: %s\n", product.Name)
	return
}

func UpdateProduct(product model.Product, newProduct model.ProductDto) {
	product.Price = newProduct.Price
	product.Status = newProduct.Status
	database.Save(&product)

	commitPriceChange(newProduct.Price, product.ID)
}

func DeleteAllUserProduct(userId int) {
	database.Where("user_id = ?", userId).Delete(&model.Pivot{})
}

func DeleteProduct(productId string, userId int) string {
	var ids []int
	id, _ := strconv.Atoi(productId)

	database.
		Select("pivots.id").
		Model(&model.Pivot{}).
		Joins("JOIN products product on product.id = pivots.product_id").
		Where("product.id = ? AND pivots.user_id = ?", id, userId).
		Find(&ids)

	if len(ids) == 0 {
		return "پروداکت یافت نشد!"
	}

	var deletedPivot model.Pivot
	database.
		Model(&model.Pivot{}).
		Where("id IN ?", ids).
		First(&deletedPivot).
		Delete(&model.Pivot{})

	product := GetProductById(deletedPivot.ProductID)

	msg := service.CreateDeleteProductSuccessfulMsg(product.ToDto())

	return msg
}

func GetProductList(userId int) string {
	products := GetAllProductByUserId(userId)
	return service.CreateProductListMsg(products)
}

func commitPriceChange(price int, productID int) {
	database.Create(&model.History{Price: price, ProductID: productID, Date: time.Now()})
}

func GetGraphPicName(productId string) (string, error) {
	pid, _ := strconv.Atoi(productId)
	var prices []model.GraphData
	database.
		Model(&model.History{}).
		Joins("JOIN products product on product.id = histories.product_id").
		Where("product.id = ? AND histories.price > 0", pid).
		Find(&prices)

	if len(prices) < 3 {
		return "", errors.New("تعداد قیمت ثبت‌شده کمتر از ۳ هست")
	}

	imagePath, err := utils.LinearRegreasion(prices)

	if err != nil {
		log.Println(err)
		return "", errors.New("خطا در ساخت تصویر")
	}

	return imagePath, nil
}

func GetHistoryPicName(productId string) (string, error) {
	pid, _ := strconv.Atoi(productId)
	var prices []model.GraphData
	database.
		Model(&model.BulkHistory{}).
		Joins("JOIN products product on product.real_id = bulk_histories.source_product_id").
		Where("product.id = ? AND bulk_histories.price > 0", pid).
		Find(&prices)

	if len(prices) < 3 {
		return "", errors.New("تعداد قیمت ثبت‌شده کمتر از ۳ هست")
	}

	//imagePath, err := graph.LinearRegreasion(prices)
	imagePath, err := utils.StockAnalysis(prices)

	if err != nil {
		log.Println(err)
		return "", errors.New("خطا در ساخت تصویر")
	}

	return imagePath, nil
}
