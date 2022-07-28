package job

import (
	"digi-bot/model"
	"digi-bot/model/db"
	"digi-bot/service"
	"digi-bot/service/bot"
	"digi-bot/service/changeDetector"
	"log"
	"time"
)

func ChangeDetectorJob() {
	for {
		log.Printf("--------------- Scheduler --------------------\n")

		updateCount := refresh(&changeDetector.DigikalaChangeDetector{})

		log.Printf("compare finished. num of update: %d \n", updateCount)
		log.Printf("--------------- done --------------------\n")

		time.Sleep(time.Hour * 2)
	}
}

func refresh(c changeDetector.ChangeDetector) int {
	updateCount := 0
	c.Detect(func(newProduct model.ProductDto, product model.Product) {

		if message, isChanged, changeLevel := compare(newProduct, product.ToDto()); isChanged {
			log.Printf("old price: %d, new price: %d",
				product.Price,
				newProduct.Price)
			available := isChanged && newProduct.Status != 0
			bot.GetTelegramBot().SendUpdateForUsers(product.ID, message, available, changeLevel)
			db.UpdateProduct(product, newProduct)
			updateCount++
		}
	})

	return updateCount
}

func compare(newProduct model.ProductDto, oldProduct model.ProductDto) (message string, isChanged bool, changeLevel int) {
	if newProduct.Price == oldProduct.Price {
		return "", false, 0
	}

	if newProduct.Price == 0 {
		return service.CreateNotAvailableMsg(newProduct), true, 1
	}

	changeLevel = 1
	if newProduct.Status == 2 {
		changeLevel = 2
	}

	return service.CreateNormalPriceChangeMsg(newProduct, newProduct.Price, oldProduct.Price), true, changeLevel
}
