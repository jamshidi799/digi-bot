package job

import (
	"digi-bot/crawler/digikalaCrawler"
	"digi-bot/db"
	"digi-bot/messageCreator"
	"digi-bot/model"
	productService "digi-bot/service/product"
	userService "digi-bot/service/user"
	"log"
	"time"
)

func ChangeDetectorJob() {
	for {
		log.Printf("--------------- Scheduler --------------------\n")

		updateCount := refresh()

		log.Printf("compare finished. num of update: %d \n", updateCount)
		log.Printf("--------------- done --------------------\n")

		time.Sleep(time.Hour * 2)
	}
}

func refresh() int {
	updateCount := 0
	products := db.GetAllProduct()

	for _, product := range products {
		newProduct, err := crawler.DigikalaCrawler{}.Crawl(product.Url)
		if err != nil {
			log.Println(newProduct)
			continue
		}

		if message, isChanged, changeLevel := compare(newProduct, product.ToDto()); isChanged {
			log.Printf("old price: %d, new price: %d",
				product.Price,
				newProduct.Price)
			userService.SendProductUpdateToUsers(product.ID, message, changeLevel)
			productService.UpdateProduct(product, newProduct)
			updateCount++
		}

		//break
		time.Sleep(time.Second * 3)
	}
	return updateCount
}

func compare(newProduct model.ProductDto, oldProduct model.ProductDto) (message string, isChanged bool, changeLevel int) {
	if newProduct.Price == oldProduct.Price {
		return "", false, 0
	}

	if newProduct.Price == 0 {
		return messageCreator.CreateNotAvailableMsg(newProduct), true, 1
	}

	changeLevel = 1
	if newProduct.Status == 2 {
		changeLevel = 2
	}

	return messageCreator.CreateNormalPriceChangeMsg(newProduct, newProduct.Price, oldProduct.Price), true, changeLevel
}
