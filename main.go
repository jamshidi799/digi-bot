package main

import (
	"digi-bot/bot"
	"digi-bot/crawler"
	"digi-bot/db"
	"digi-bot/messageCreator"
	"digi-bot/model"
	"github.com/prprprus/scheduler"
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
	group.Add(1)

	s, err := scheduler.NewScheduler(1000)
	if err != nil {
		panic(err)
	}
	s.Every().Second(10).Do(Scheduler)
	//s.Every().Hour(2).Do(Scheduler)
	//Scheduler()
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

		log.Printf("old price: %d, new price: %d",
			product.Price,
			newProduct.Price)

		if message, isChanged := changeDetector(newProduct, product.ToProduct()); isChanged {
			bot.SendUpdateForUser(product.UserId,
				message)

			updateCount++
		}

		//productService.UpdateProduct(product, newProduct)

		//break
		time.Sleep(time.Second * 3)
	}

	log.Printf("changeDetector finished. num of update: %d \n", updateCount)
	log.Printf("--------------- done --------------------\n")

}

func changeDetector(newProduct model.Product, oldProduct model.Product) (message string, isChanged bool) {
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
