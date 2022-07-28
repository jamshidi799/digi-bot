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
	tlBot.handleSetting()
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
			err = db.AddProductToDB(product, int(c.Sender().ID))
			if err != nil {
				return c.Send(err.Error())
			} else {
				message := service.CreatePreviewMsg(product)
				return c.Send(
					message,
					&tele.SendOptions{
						ParseMode:   "HTML",
						ReplyMarkup: getProductSelector(product.Id),
					})
			}
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
			err := c.Reply(err)
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

func (tlBot TelegramBot) handleSetting() {
	bot := tlBot.bot

	selector := &tele.ReplyMarkup{}
	btnSetting := selector.Data("تنظیمات", "setting")
	btnOne := selector.Data("1", "one")
	btnTwo := selector.Data("2", "two")

	bot.Handle(&btnSetting, func(c tele.Context) error {
		commandLogs("setting", c.Sender().ID)
		msg := service.CreateChangeSettingGuide()
		productId := c.Data()
		return c.Reply(msg, &tele.SendOptions{
			ParseMode:   "HTML",
			ReplyMarkup: getProductSettingSelector(productId),
		})
	})

	bot.Handle(&btnOne, func(c tele.Context) error {
		productId := c.Data()
		userId := c.Sender().ID
		msg := db.UpdateStatus(1, productId, userId)
		return c.Reply(msg, &tele.SendOptions{
			ParseMode: "HTML",
		})
	})

	bot.Handle(&btnTwo, func(c tele.Context) error {
		productId := c.Data()
		userId := c.Sender().ID
		msg := db.UpdateStatus(2, productId, userId)

		return c.Reply(msg, &tele.SendOptions{
			ParseMode: "HTML",
		})
	})
}

func (tlBot TelegramBot) SendUpdateForUsers(productId int, message string, available bool, changeLevel int) {
	usersId := db.GetAllUsersIdByProductId(productId, changeLevel)

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
	btnSetting := selector.Data("تنظیمات", "setting", productIdStr)

	selector.Inline(
		selector.Row(btnGraph, btnDelete),
		selector.Row(btnSetting),
	)
	return selector
}

func getProductSettingSelector(productId string) *tele.ReplyMarkup {
	selector := &tele.ReplyMarkup{}
	btnOne := selector.Data("1", "one", productId)
	btnTwo := selector.Data("2", "two", productId)
	//btnThree := selector.Data("3", "three", productId)

	selector.Inline(
		selector.Row(btnOne, btnTwo),
	)
	return selector
}

func commandLogs(command string, userId int64) {
	log.Printf("command: %s, userId: %d", command, userId)
}
