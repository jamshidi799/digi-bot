package bot

import (
	"digi-bot/config"
	"digi-bot/crawler"
	"digi-bot/model"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var Bot *tb.Bot

func Run() {
	var token = "1700540701:AAGiNrhQNdha0FJVm9icPiv4VghZw7o1eE8"
	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	Bot = b

	b.Handle("/start", func(m *tb.Message) {
		userModel := model.ToUserModel(m.Sender)
		config.DB.Create(&userModel)
		fmt.Printf("%+v\n", userModel)
		fmt.Printf("%+v\n", m.Sender)

		b.Send(m.Sender, "آدرس کالا را وارد کنید")
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		err := addObjectToDB(m.Sender.ID, m.Text)
		if err != nil {
			b.Send(m.Sender, err.Error())
		} else {
			b.Send(m.Sender, "کالا با موفقیت ذخیره شد. برای اضافه کردن کالای جدید کافی است فقط آدرس آن را وارد کنید")
		}

	})

	b.Start()

}

func Send(chatId int, message string) {
	var user model.UserModel
	config.DB.First(&user) // todo: use chat id
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
	config.DB.Create(&objModel)
	return nil
}
