package bot

import (
	"digi-bot/db"
	"digi-bot/messageCreator"
	"digi-bot/model"
	productService "digi-bot/service/product"
	"log"
	"strconv"
	"sync"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var Bot *tb.Bot

// todo: create interface and add other clients
func Run(group *sync.WaitGroup) {
	//var token = "1767506686:AAFP-w40DbrbhhVQLk6g2aGz3KqhF5oZugI"	// digiBot
	var token = "1700540701:AAGiNrhQNdha0FJVm9icPiv4VghZw7o1eE8" // goTestBot
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

	selector := &tb.ReplyMarkup{}
	btnGraph := selector.Data("نمودار قیمت", "graph")
	btnDelete := selector.Data("حذف", "delete")

	bot.Handle("/start", func(m *tb.Message) {
		userModel := model.ToUser(m.Sender)
		db.DB.Create(&userModel)

		_, _ = bot.Send(m.Sender, messageCreator.CreateHelpMsg(), &tb.SendOptions{
			ParseMode: "HTML",
		})
	})

	bot.Handle("/deleteall", func(m *tb.Message) {
		productService.DeleteAllUserProduct(m.Sender.ID)
		bot.Reply(m, "لیست کالا با موفقیت پاک شد")
	})

	bot.Handle("/add", func(m *tb.Message) {
		bot.Reply(m, "آدرس (url) کالا را وارد کنید")
		bot.Handle(tb.OnText, func(m *tb.Message) {
			product, productId, err := productService.AddProductToDB(m.Sender.ID, m.Text)
			if err != nil {
				_, _ = bot.Send(m.Sender, err.Error())
			} else {
				message := messageCreator.CreatePreviewMsg(product)
				_, _ = bot.Send(
					m.Sender,
					message,
					&tb.SendOptions{
						ParseMode:   "HTML",
						ReplyMarkup: getProductSelector(productId),
					})
			}
		})
	})

	bot.Handle("/help", func(m *tb.Message) {
		bot.Send(m.Sender, messageCreator.CreateHelpMsg(), &tb.SendOptions{
			ParseMode: "HTML",
		})
	})

	bot.Handle("/list", func(m *tb.Message) {
		_, err = bot.Send(m.Sender, productService.GetProductList(m.Sender.ID), &tb.SendOptions{
			ParseMode: "HTML",
		})
	})

	bot.Handle(&btnGraph, func(c *tb.Callback) {
		productId, _ := strconv.Atoi(c.Data)
		imagePath := productService.GetGraphPicName(productId)
		image := &tb.Photo{File: tb.FromDisk(imagePath)}
		bot.Reply(c.Message, image)
	})

	bot.Handle(&btnDelete, func(c *tb.Callback) {
		msg := productService.DeleteProduct(c.Data, c.Sender.ID)
		bot.Reply(c.Message, msg, &tb.SendOptions{
			ParseMode: "HTML",
		})
	})

	bot.Start()

}

func SendUpdateForUsers(usersId []int, productId int, message string) {
	for _, userId := range usersId {
		user := db.GetUserById(userId)
		_, _ = Bot.Send(user.ToTbUser(), message, &tb.SendOptions{
			ParseMode:   "HTML",
			ReplyMarkup: getProductSelector(productId),
		})
	}
}

func getProductSelector(productId int) *tb.ReplyMarkup {
	selector := &tb.ReplyMarkup{}
	btnGraph := selector.Data("نمودار قیمت", "graph", strconv.Itoa(productId))
	btnDelete := selector.Data("حذف", "delete", strconv.Itoa(productId))

	selector.Inline(
		selector.Row(btnGraph, btnDelete),
	)
	return selector
}
