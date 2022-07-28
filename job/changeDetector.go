package job

import (
	"digi-bot/model"
	"digi-bot/model/db"
	"digi-bot/service"
	"digi-bot/service/bot"
	"digi-bot/service/crawler"
	"log"
	"time"
)

func ChangeDetectorJob() {
	for _, c := range crawler.Crawlers {
		go run(c)
	}
}

func run(c crawler.Crawler) {
	for {
		log.Printf("--------------- Scheduler for %v --------------------\n", c.GetDomain())
		updateCount := Detect(c)
		log.Printf("compare finished. num of updates: %d \n", updateCount)
		log.Printf("--------------- %v done --------------------\n", c.GetDomain())
		time.Sleep(time.Hour * 2)
	}
}

func getProducts(domain string) []model.Product {
	return db.GetAllProductByDomain(domain)
}

func Detect(c crawler.Crawler) int {
	counter := 0
	for _, product := range getProducts(c.GetDomain()) {
		newProduct, err := c.Crawl(product.Url)

		if err != nil {
			log.Println(err, newProduct)
			continue
		}

		if handleChange(newProduct, product) {
			counter++
		}

		time.Sleep(time.Second * 1)
	}
	return counter
}

func handleChange(newProduct model.ProductDto, product model.Product) (isChanged bool) {
	message, isChanged, changeLevel := compare(newProduct, product.ToDto())
	if !isChanged {
		return
	}
	log.Printf("old price: %d, new price: %d", product.Price, newProduct.Price)
	available := isChanged && newProduct.Status != 0
	bot.GetTelegramBot().SendUpdateForUsers(product.ID, message, available, changeLevel)
	db.UpdateProduct(product, newProduct)
	return
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
