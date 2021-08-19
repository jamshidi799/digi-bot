package bot

import (
	"digi-bot/db"
	"digi-bot/messageCreator"
	"digi-bot/model"
	"digi-bot/service/pivot"
	productService "digi-bot/service/product"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var Bot *tb.Bot

// todo: create interface and add other clients
func Run(group *sync.WaitGroup) {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error loading .env file")
	}

	token := os.Getenv("BOT_TOKEN")
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
	btnSetting := selector.Data("تنظیمات", "setting")
	btnOne := selector.Data("1", "one")
	btnTwo := selector.Data("2", "two")

	bot.Handle("/start", func(m *tb.Message) {
		userModel := model.ToUser(m.Sender)
		db.DB.Create(&userModel)

		_, _ = bot.Send(m.Sender, messageCreator.CreateHelpMsg(), &tb.SendOptions{
			ParseMode: "HTML",
		})

		commandLogs("start", m.Sender.ID)
	})

	bot.Handle("/deleteall", func(m *tb.Message) {
		productService.DeleteAllUserProduct(m.Sender.ID)
		bot.Reply(m, "لیست کالا با موفقیت پاک شد")

		commandLogs("delete all", m.Sender.ID)
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

		commandLogs("add", m.Sender.ID)
	})

	bot.Handle("/help", func(m *tb.Message) {
		bot.Send(m.Sender, messageCreator.CreateHelpMsg(), &tb.SendOptions{
			ParseMode: "HTML",
		})

		commandLogs("help", m.Sender.ID)
	})

	bot.Handle("/list", func(m *tb.Message) {
		_, err = bot.Send(m.Sender, productService.GetProductList(m.Sender.ID), &tb.SendOptions{
			ParseMode: "HTML",
		})

		commandLogs("list", m.Sender.ID)
	})

	bot.Handle(&btnGraph, func(c *tb.Callback) {
		imagePath, err := productService.GetGraphPicName(c.Data)
		if err != nil {
			bot.Reply(c.Message, err)
		} else {
			image := &tb.Photo{File: tb.FromDisk(imagePath)}
			bot.Reply(c.Message, image)
		}
		log.Println(imagePath)
		commandLogs("graph", c.Sender.ID)
	})

	bot.Handle(&btnDelete, func(c *tb.Callback) {
		msg := productService.DeleteProduct(c.Data, c.Sender.ID)
		bot.Reply(c.Message, msg, &tb.SendOptions{
			ParseMode: "HTML",
		})

		commandLogs("delete", c.Sender.ID)
	})

	bot.Handle(&btnSetting, func(c *tb.Callback) {
		msg := messageCreator.CreateChangeSettingGuide()
		productId := c.Data
		bot.Reply(c.Message, msg, &tb.SendOptions{
			ParseMode:   "HTML",
			ReplyMarkup: getProductSettingSelector(productId),
		})

		commandLogs("setting", c.Sender.ID)
	})

	bot.Handle(&btnOne, func(c *tb.Callback) {
		productId := c.Data
		userId := c.Sender.ID
		msg := pivot.UpdateStatus(1, productId, userId)
		bot.Reply(c.Message, msg, &tb.SendOptions{
			ParseMode: "HTML",
		})
	})

	bot.Handle(&btnTwo, func(c *tb.Callback) {
		productId := c.Data
		userId := c.Sender.ID
		msg := pivot.UpdateStatus(2, productId, userId)

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
	productIdStr := strconv.Itoa(productId)
	selector := &tb.ReplyMarkup{}
	btnGraph := selector.Data("نمودار قیمت", "graph", productIdStr)
	btnDelete := selector.Data("حذف", "delete", productIdStr)
	btnSetting := selector.Data("تنظیمات", "setting", productIdStr)

	selector.Inline(
		selector.Row(btnGraph, btnDelete),
		selector.Row(btnSetting),
	)
	return selector
}

func getProductSettingSelector(productId string) *tb.ReplyMarkup {
	selector := &tb.ReplyMarkup{}
	btnOne := selector.Data("1", "one", productId)
	btnTwo := selector.Data("2", "two", productId)
	//btnThree := selector.Data("3", "three", productId)

	selector.Inline(
		selector.Row(btnOne, btnTwo),
	)
	return selector
}

func commandLogs(command string, userId int) {
	log.Printf("command: %s, userId: %d", command, userId)
}
