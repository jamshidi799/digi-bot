package bot

import (
	"digi-bot/config"
	"digi-bot/model"
	"fmt"
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

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

	//b.Handle("/hello", func(m *tb.Message) {
	//	userModel := model.ToUserModel(m.Sender)
	//	config.DB.Create(&userModel)
	//	fmt.Printf("%+v\n", userModel)
	//	fmt.Printf("%+v\n", m.Sender)
	//
	//	b.Send(m.Sender, "Hello World!")
	//})

	var user model.UserModel
	config.DB.First(&user)
	fmt.Printf("%+v\n", user)
	fmt.Printf("%+v\n", user.ToTbUser())
	_, err = b.Send(user.ToTbUser(), "fuck me")
	if err != nil {
		log.Fatal(err)
	}

	b.Start()

}
