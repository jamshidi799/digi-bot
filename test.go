package main

import (
	"digi-bot/crawler"
	"fmt"
	"log"
)

func main() {
	//db := db.Init()
	//test := model.TestModel{Username: "mohammad", FirstName: "jamshidi", LastName: "amoo mama"}
	//db.Create(&test)

	obj, err := crawler.Crawl("https://www.digikala.com/product/dkp-4645079/%D9%84%D9%BE-%D8%AA%D8%A7%D9%BE-14-%D8%A7%DB%8C%D9%86%DA%86%DB%8C-%D8%A7%DB%8C%D8%B3%D9%88%D8%B3-%D9%85%D8%AF%D9%84-zenbook-14-um433i")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", obj)
}
