package main

import (
	"digi-bot/bot"
	"digi-bot/crawler"
	"digi-bot/db"
	"digi-bot/messageCreator"
	"digi-bot/model"
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
	Scheduler()
}

func Scheduler() {
	objects := db.GetAllProduct()

	//fmt.Printf("%+v", objects)

	for _, object := range objects {
		obj, err := crawler.Crawl(object.Url)
		if err != nil {
			log.Println(obj)
		}

		log.Printf("old price: %d, new price: %d",
			object.Price,
			obj.Price)

		if message, isChanged := changeDetector(obj, object.ToObject()); isChanged {
			bot.SendUpdateForUser(object.UserId,
				object.Image,
				message)
		}
		// todo: write new obj to db
		break
		time.Sleep(time.Second * 3)
	}
}

func changeDetector(newObj model.Object, oldObj model.Object) (message string, isChanged bool) {
	if newObj.Price == oldObj.Price {
		return "", false
	}

	if newObj.Price == 0 {
		if oldObj.Price == 0 {
			return "", false
		}
		return messageCreator.CreateNotAvailableMsg(newObj), true
	}

	return messageCreator.CreateNormalPriceChangeMsg(newObj, newObj.Price, oldObj.Price), true
}
