package bot

import (
	"digi-bot/db"
	"digi-bot/messageCreator"
	"digi-bot/model"
	"digi-bot/service"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var telegramBot TelegramBot

type TelegramBot struct {
	bot *tb.Bot
}

func GetTelegramBot() TelegramBot {
	return telegramBot
}

func InitTelegramBot(group *sync.WaitGroup) {
	if telegramBot.bot != nil {
		log.Println("telebot: can't reinitialize bot")
	}

	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error loading .env file")
	}

	token := os.Getenv("BOT_TOKEN")

	telegramBot.bot, err = tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
	}

	telegramBot.callHandlers()

	group.Done()
	log.Println("bot started")
	telegramBot.bot.Start()
}

func (tlBot TelegramBot) callHandlers() {
	tlBot.handleStart()
	tlBot.handleDeleteAll()
	tlBot.handleAdd()
	tlBot.handleDelete()
	tlBot.handleHelp()
	tlBot.handleList()
	tlBot.handleGraph()
	tlBot.handleHistory()
	tlBot.handleSetting()
}

func (tlBot TelegramBot) handleStart() {
	tlBot.bot.Handle("/start", func(m *tb.Message) {
		userModel := model.ToUser(m.Sender)
		db.DB.Create(&userModel)

		_, _ = tlBot.bot.Send(m.Sender, messageCreator.CreateHelpMsg(), &tb.SendOptions{
			ParseMode: "HTML",
		})

		commandLogs("start", m.Sender.ID)
	})
}

func (tlBot TelegramBot) handleDeleteAll() {
	bot := tlBot.bot
	bot.Handle("/deleteall", func(m *tb.Message) {
		service.DeleteAllUserProduct(m.Sender.ID)
		bot.Reply(m, "لیست کالا با موفقیت پاک شد")

		commandLogs("delete all", m.Sender.ID)
	})
}

func (tlBot TelegramBot) handleAdd() {
	bot := tlBot.bot
	bot.Handle("/add", func(m *tb.Message) {
		bot.Reply(m, "آدرس (url) کالا را وارد کنید")
		bot.Handle(tb.OnText, func(m *tb.Message) {
			product, productId, err := service.AddProductToDB(m.Sender.ID, m.Text)
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
}

func (tlBot TelegramBot) handleHelp() {
	bot := tlBot.bot
	bot.Handle("/help", func(m *tb.Message) {
		bot.Send(m.Sender, messageCreator.CreateHelpMsg(), &tb.SendOptions{
			ParseMode: "HTML",
		})

		commandLogs("help", m.Sender.ID)
	})
}

func (tlBot TelegramBot) handleList() {
	bot := tlBot.bot
	bot.Handle("/list", func(m *tb.Message) {
		bot.Send(m.Sender, service.GetProductList(m.Sender.ID), &tb.SendOptions{
			ParseMode: "HTML",
		})

		commandLogs("list", m.Sender.ID)
	})
}

func (tlBot TelegramBot) handleDelete() {
	bot := tlBot.bot
	selector := &tb.ReplyMarkup{}
	btnDelete := selector.Data("حذف", "delete")

	bot.Handle(&btnDelete, func(c *tb.Callback) {
		msg := service.DeleteProduct(c.Data, c.Sender.ID)
		bot.Reply(c.Message, msg, &tb.SendOptions{
			ParseMode: "HTML",
		})

		commandLogs("delete", c.Sender.ID)
	})
}

func (tlBot TelegramBot) handleGraph() {
	bot := tlBot.bot

	selector := &tb.ReplyMarkup{}
	btnGraph := selector.Data("نمودار قیمت", "graph")

	bot.Handle(&btnGraph, func(c *tb.Callback) {
		imagePath, err := service.GetGraphPicName(c.Data)
		if err != nil {
			bot.Reply(c.Message, err)
		} else {
			image := &tb.Photo{File: tb.FromDisk(imagePath)}
			bot.Reply(c.Message, image)
		}
		log.Println(imagePath)
		commandLogs("graph", c.Sender.ID)
	})
}

func (tlBot TelegramBot) handleHistory() {
	bot := tlBot.bot

	selector := &tb.ReplyMarkup{}
	btnHistory := selector.Data("نمودار بلند‌مدت", "history")

	bot.Handle(&btnHistory, func(c *tb.Callback) {
		imagePath, err := service.GetHistoryPicName(c.Data)
		if err != nil {
			bot.Reply(c.Message, err)
		} else {
			image := &tb.Photo{File: tb.FromDisk(imagePath)}
			bot.Reply(c.Message, image)
		}
		log.Println(imagePath)
		commandLogs("history", c.Sender.ID)
	})
}

func (tlBot TelegramBot) handleSetting() {
	bot := tlBot.bot

	selector := &tb.ReplyMarkup{}
	btnSetting := selector.Data("تنظیمات", "setting")
	btnOne := selector.Data("1", "one")
	btnTwo := selector.Data("2", "two")

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
		msg := service.UpdateStatus(1, productId, userId)
		bot.Reply(c.Message, msg, &tb.SendOptions{
			ParseMode: "HTML",
		})
	})

	bot.Handle(&btnTwo, func(c *tb.Callback) {
		productId := c.Data
		userId := c.Sender.ID
		msg := service.UpdateStatus(2, productId, userId)

		bot.Reply(c.Message, msg, &tb.SendOptions{
			ParseMode: "HTML",
		})
	})
}

func (tlBot TelegramBot) SendUpdateForUsers(usersId []int, productId int, message string, available bool) {
	rand.Seed(time.Now().UnixNano())
	for _, userId := range usersId {
		user := db.GetUserById(userId)
		msg, _ := tlBot.bot.Send(user.ToTbUser(), message, &tb.SendOptions{
			ParseMode:   "HTML",
			ReplyMarkup: getProductSelector(productId),
		})
		if !available {
			random := rand.Intn(5)
			if random == 0 {
				tlBot.bot.Reply(msg, "💩")
			} else {
				gif := &tb.Animation{File: tb.FromDisk(fmt.Sprintf("assets/gif%d.mp4", random))}
				tlBot.bot.Reply(msg, gif)
			}
		}
	}
}

func getProductSelector(productId int) *tb.ReplyMarkup {
	productIdStr := strconv.Itoa(productId)
	selector := &tb.ReplyMarkup{}
	btnGraph := selector.Data("نمودار قیمت", "graph", productIdStr)
	btnDelete := selector.Data("حذف", "delete", productIdStr)
	btnSetting := selector.Data("تنظیمات", "setting", productIdStr)
	btnHistory := selector.Data("نمودار بلند‌مدت", "history", productIdStr)

	selector.Inline(
		selector.Row(btnGraph, btnDelete),
		selector.Row(btnSetting, btnHistory),
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
