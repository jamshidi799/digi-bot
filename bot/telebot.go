package bot

import (
	"digi-bot/crawler"
	"digi-bot/db"
	"digi-bot/messageCreator"
	"digi-bot/model"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var Bot *tb.Bot

// todo: create interface and add other clients
func Run(group *sync.WaitGroup) {
	var token = "1700540701:AAGiNrhQNdha0FJVm9icPiv4VghZw7o1eE8"
	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	Bot = bot
	group.Done()

	bot.Handle("/start", func(m *tb.Message) {
		userModel := model.ToUserModel(m.Sender)
		db.DB.Create(&userModel)

		_, _ = bot.Send(m.Sender, "آدرس کالا را وارد کنید")
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		product, err := addObjectToDB(m.Sender.ID, m.Text)
		if err != nil {
			_, _ = bot.Send(m.Sender, err.Error())
		} else {
			message := messageCreator.CreatePreviewMsg(product)
			log.Println(message)
			_, _ = bot.Send(m.Sender, message, &tb.SendOptions{
				ParseMode: "HTML",
			})
		}
	})

	// todo: handle /delete and /help command

	bot.Start()

}

func SendUpdateForUser(chatId int, imageUrl string, message string) {
	user := db.GetUserById(chatId)
	//photo := &tb.Photo{File: tb.FromURL(imageUrl)}
	//imageMsg, _ := Bot.Send(user.ToTbUser(), photo)
	//Bot.Reply(imageMsg, message, &tb.SendOptions{
	//	ParseMode: "HTML",
	//})

	_, _ = Bot.Send(user.ToTbUser(), message, &tb.SendOptions{
		ParseMode: "HTML",
	})
}

func Send(chatId int, message string) {
	user := db.GetUserById(chatId)
	fmt.Printf("%+v\n", user)
	fmt.Printf("%+v\n", user.ToTbUser())
	_, err := Bot.Send(user.ToTbUser(), message)
	if err != nil {
		log.Fatal(err)
	}
}

func addObjectToDB(senderId int, url string) (model.Product, error) {
	fmt.Printf("%+v %+v", senderId, url)
	if res := strings.Contains(url, "digikala.com"); !res {
		return model.Product{}, nil
	}

	product, err := crawler.Crawl(url)
	if err != nil {
		return model.Product{}, nil
	}

	fmt.Printf("%+v", product)
	productModel := product.ToProductModel(senderId)
	db.DB.Create(&productModel)
	return productModel.ToProduct(), nil
}
