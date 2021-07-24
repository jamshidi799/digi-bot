package main

import (
	"digi-bot/bot"
	"digi-bot/config"
	"digi-bot/crawler"
	"digi-bot/model"
	str "digi-bot/stringUtility"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	config.Init()
	var group sync.WaitGroup
	group.Add(1)
	go bot.Run(&group)
	group.Wait()
	log.Println("bot started")
	Scheduler()
}

func Scheduler() {
	var objects []model.ObjectModel
	config.DB.Find(&objects)

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
		return createNotAvailableMsg(newObj), true
	}

	return createNormalPriceChangeMsg(newObj, newObj.Price, oldObj.Price), true
}

func createNormalPriceChangeMsg(object model.Object, newPrice int, oldPrice int) string {
	output := createHeader(object).Append(fmt.Sprintf("%s -> %s", str.Number(oldPrice).AddComma(), str.Number(newPrice).AddComma()))
	return output.ToString()
}

func createNotAvailableMsg(obj model.Object) string {
	output := createHeader(obj).
		Append("ناموجود!")

	return output.ToString()
}

func createHeader(obj model.Object) str.String {
	output := str.
		String(obj.Name).
		Bold().
		ToLink(obj.Url).
		AddNewLine()

	return output
}
