package product

import (
	"digi-bot/crawler"
	"digi-bot/db"
	"digi-bot/graph"
	"digi-bot/messageCreator"
	"digi-bot/model"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

func AddProductToDB(senderId int, url string) (model.ProductDto, int, error) {
	//fmt.Printf("%+v %+v\n", senderId, url)
	if res := strings.Contains(url, "digikala.com"); !res {
		return model.ProductDto{}, 0, errors.New("ادرس نامعتبر است")
	}

	product, err := crawler.Crawl(url)
	if err != nil {
		return model.ProductDto{}, 0, err
	}

	//fmt.Printf("%+v", product)
	productModel := product.ToProduct()

	result := db.DB.Where(productModel).Find(&productModel)
	if result.RowsAffected == 0 {
		productModel = product.ToProduct()
		db.DB.Create(&productModel)
	}

	pivot := model.Pivot{UserID: senderId, ProductID: productModel.ID, NotificationSetting: 1}
	db.DB.Create(&pivot)

	log.Printf("new product added: %s\n", product.Name)
	return product, productModel.ID, nil
}

func UpdateProduct(product model.Product, newProduct model.ProductDto) {
	product.Price = newProduct.Price
	product.Status = newProduct.Status
	db.DB.Save(&product)

	commitPriceChange(newProduct.Price, product.ID)
}

func DeleteAllUserProduct(userId int) {
	db.DB.Where("user_id = ?", userId).Delete(&model.Pivot{})
}

func DeleteProduct(productId string, userId int) string {
	var ids []int
	id, _ := strconv.Atoi(productId)

	db.DB.
		Select("pivots.id").
		Model(&model.Pivot{}).
		Joins("JOIN products product on product.id = pivots.product_id").
		Where("product.id = ? AND pivots.user_id = ?", id, userId).
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

func commitPriceChange(price int, productID int) {
	db.DB.Create(&model.History{Price: price, ProductID: productID, Date: time.Now()})
}

func GetGraphPicName(productId string) (string, error) {
	pid, _ := strconv.Atoi(productId)
	var prices []model.History
	db.DB.
		Model(&model.History{}).
		Joins("JOIN products product on product.id = histories.product_id").
		Where("product.id = ?", pid).
		Find(&prices)

	if len(prices) < 3 {
		return "", errors.New("تعداد قیمت ثبت‌شده کمتر از ۳ هست")
	}

	imagePath, err := graph.LinearRegreasion(prices)

	if err != nil {
		log.Println(err)
		return "", errors.New("خطا در ساخت تصویر")
	}

	return imagePath, nil
}
