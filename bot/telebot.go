package bot

import (
	"digi-bot/crawler"
	"digi-bot/db"
	"digi-bot/model"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var Bot *tb.Bot

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
		fmt.Printf("%+v\n", userModel)
		fmt.Printf("%+v\n", m.Sender)

		_, _ = bot.Send(m.Sender, "آدرس کالا را وارد کنید")
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		err := addObjectToDB(m.Sender.ID, m.Text)
		if err != nil {
			_, _ = bot.Send(m.Sender, err.Error())
		} else {
			// todo: send object and its current status
			_, _ = bot.Send(m.Sender, "کالا با موفقیت ذخیره شد. برای اضافه کردن کالای جدید کافی است فقط آدرس آن را وارد کنید")
		}
	})

	bot.Start()

}

func SendUpdateForUser(chatId int, imageUrl string, message string) {
	user := db.GetUserById(chatId)
	//photo := &tb.Photo{File: tb.FromURL(imageUrl)}
	//imageMsg, _ := Bot.Send(user.ToTbUser(), photo)
	//Bot.Reply(imageMsg, message, &tb.SendOptions{
	//	ParseMode: "HTML",
	//})

	Bot.Send(user.ToTbUser(), message, &tb.SendOptions{
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

func addObjectToDB(senderId int, url string) error {
	fmt.Printf("%+v %+v", senderId, url)
	if res := strings.Contains(url, "digikala.com"); !res {
		return errors.New("آدرس وارد شده نامعتبر می‌باشد")
	}

	obj, err := crawler.Crawl(url)
	if err != nil {
		return err
	}

	fmt.Printf("%+v", obj)
	objModel := obj.ToObjectModel(senderId, url)
	db.DB.Create(&objModel)
	return nil
}
