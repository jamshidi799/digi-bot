package bot

import (
	"digi-bot/model/db"
	"digi-bot/service"
	"digi-bot/service/crawler"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	tele "gopkg.in/telebot.v3"
)

var telegramBot TelegramBot

type TelegramBot struct {
	bot *tele.Bot
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

	telegramBot.bot, err = tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
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
}

func (tlBot TelegramBot) handleStart() {
	tlBot.bot.Handle("/start", func(c tele.Context) error {
		commandLogs("start", c.Sender().ID)
		db.SaveUser(c.Sender())

		return c.Send(c.Sender(), service.CreateHelpMsg(), &tele.SendOptions{
			ParseMode: "HTML",
		})
	})
}

func (tlBot TelegramBot) handleDeleteAll() {
	bot := tlBot.bot
	bot.Handle("/deleteall", func(c tele.Context) error {
		commandLogs("delete all", c.Sender().ID)
		db.DeleteAllUserProduct(int(c.Sender().ID))

		return c.Reply(c.Message(), "لیست کالا با موفقیت پاک شد")
	})
}

func (tlBot TelegramBot) handleAdd() {
	bot := tlBot.bot

	bot.Handle("/add", func(c tele.Context) error {
		commandLogs("add", c.Sender().ID)
		err := c.Reply("آدرس (url) کالا را وارد کنید")

		bot.Handle(tele.OnText, func(c tele.Context) error {
			product, err := crawler.Crawl(c.Text())
			if err != nil {
				return c.Send(err.Error())
			}

			productId := db.AddProductToDB(product, int(c.Sender().ID))
			message := service.CreatePreviewMsg(product)
			return c.Send(
				message,
				&tele.SendOptions{
					ParseMode:   "HTML",
					ReplyMarkup: getProductSelector(productId),
				})
		})
		return err
	})
}

func (tlBot TelegramBot) handleHelp() {
	bot := tlBot.bot
	bot.Handle("/help", func(c tele.Context) error {
		commandLogs("help", c.Sender().ID)
		return c.Send(service.CreateHelpMsg(), &tele.SendOptions{
			ParseMode: "HTML",
		})
	})
}

func (tlBot TelegramBot) handleList() {
	bot := tlBot.bot
	bot.Handle("/list", func(c tele.Context) error {
		commandLogs("list", c.Sender().ID)
		return c.Send(db.GetProductList(int(c.Sender().ID)), &tele.SendOptions{
			ParseMode: "HTML",
		})
	})
}

func (tlBot TelegramBot) handleDelete() {
	bot := tlBot.bot
	selector := &tele.ReplyMarkup{}
	btnDelete := selector.Data("حذف", "delete")

	bot.Handle(&btnDelete, func(c tele.Context) error {
		commandLogs("delete", c.Sender().ID)
		msg := db.DeleteProduct(c.Data(), c.Sender().ID)
		return c.Reply(msg, &tele.SendOptions{
			ParseMode: "HTML",
		})
	})
}

func (tlBot TelegramBot) handleGraph() {
	bot := tlBot.bot

	selector := &tele.ReplyMarkup{}
	btnGraph := selector.Data("نمودار قیمت", "graph")

	bot.Handle(&btnGraph, func(c tele.Context) error {
		commandLogs("graph", c.Sender().ID)
		imagePath, err := db.GetGraphPicName(c.Data())
		if err != nil {
			err := c.Reply(err.Error())
			if err != nil {
				return err
			}
		} else {
			image := &tele.Photo{File: tele.FromDisk(imagePath)}
			err := c.Reply(image)
			if err != nil {
				return err
			}
		}
		log.Println(imagePath)
		return err
	})
}

func (tlBot TelegramBot) SendUpdateForUsers(productId int, message string, available bool) {
	usersId := db.GetAllUsersIdByProductId(productId)

	rand.Seed(time.Now().UnixNano())
	for _, userId := range usersId {
		user := db.GetUserById(userId)
		msg, _ := tlBot.bot.Send(user.ToTbUser(), message, &tele.SendOptions{
			ParseMode:   "HTML",
			ReplyMarkup: getProductSelector(productId),
		})
		if !available {
			random := rand.Intn(5)
			if random == 0 {
				tlBot.bot.Reply(msg, "💩")
			} else {
				gif := &tele.Animation{File: tele.FromDisk(fmt.Sprintf("assets/gif%d.mp4", random))}
				tlBot.bot.Reply(msg, gif)
			}
		}
	}
}

func getProductSelector(productId int) *tele.ReplyMarkup {
	productIdStr := strconv.Itoa(productId)
	selector := &tele.ReplyMarkup{}
	btnGraph := selector.Data("نمودار قیمت", "graph", productIdStr)
	btnDelete := selector.Data("حذف", "delete", productIdStr)

	selector.Inline(
		selector.Row(btnGraph, btnDelete),
	)
	return selector
}

func commandLogs(command string, userId int64) {
	log.Printf("command: %s, userId: %d", command, userId)
}
