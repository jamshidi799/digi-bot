package main

import (
	"digi-bot/bot"
	"digi-bot/crawler"
	"digi-bot/db"
	"digi-bot/messageCreator"
	"digi-bot/model"
	productService "digi-bot/service/product"
	"log"
	"sync"
	"time"
)

func main() {
	db.Init()
	var group sync.WaitGroup
	group.Add(1)
	go bot.Run(&group)
	group.Wait()
	log.Println("bot started")

	//go crawler.BatchCrawlScheduler()
	go crawler.StartBatchCrawl()
	group.Add(1)
	for {
		Scheduler()
		time.Sleep(time.Hour * 2)
		//time.Sleep(time.Second * 10)
	}

	group.Wait()
}

func Scheduler() {
	products := db.GetAllProduct()

	log.Printf("--------------- Scheduler --------------------\n")
	updateCount := 0

	for _, product := range products {
		newProduct, err := crawler.Crawl(product.Url)
		if err != nil {
			log.Println(newProduct)
			continue
		}

		if message, isChanged := changeDetector(newProduct, product.ToDto()); isChanged {
			log.Printf("old price: %d, new price: %d",
				product.Price,
				newProduct.Price)

			usersId := db.GetAllUsersIdByProductId(product.ID)

			log.Printf("user affected: %d", len(usersId))

			bot.SendUpdateForUsers(usersId, product.ID, message)
			productService.UpdateProduct(product, newProduct)
			updateCount++
		}

		//break
		time.Sleep(time.Second * 3)
	}

	log.Printf("changeDetector finished. num of update: %d \n", updateCount)
	log.Printf("--------------- done --------------------\n")

}

func changeDetector(newProduct model.ProductDto, oldProduct model.ProductDto) (message string, isChanged bool) {
	if newProduct.Price == oldProduct.Price {
		return "", false
	}

	if newProduct.Price == 0 {
		if oldProduct.Price == 0 {
			return "", false
		}
		return messageCreator.CreateNotAvailableMsg(newProduct), true
	}

	return messageCreator.CreateNormalPriceChangeMsg(newProduct, newProduct.Price, oldProduct.Price), true
}
