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
	tlBot.handleQuery()
	tlBot.handleDiscount()
}

func (tlBot TelegramBot) handleStart() {
	tlBot.bot.Handle("/start", func(c tele.Context) error {
		commandLogs("start", c.Sender().ID)
		db.SaveUser(c.Sender())

		return c.Send(service.CreateHelpMsg(), &tele.SendOptions{
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

			product.Id = productId
			//data, _ := json.Marshal(product)
			//kafka.Send("products", strconv.Itoa(productId), data)

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

func (tlBot TelegramBot) handleDiscount() {
	bot := tlBot.bot

	selector := &tele.ReplyMarkup{}
	btnDiscount := selector.Data("تخفیف", "discount")
	bot.Handle(&btnDiscount, func(c tele.Context) error {
		commandLogs("discount", c.Sender().ID)

		productId, _ := strconv.Atoi(c.Data())
		db.UpdateProductDiscount(productId, int(c.Sender().ID), 0)
		return c.Reply("با موفقیت انجام شد", &tele.SendOptions{
			ParseMode: "HTML",
		})
	})

	btnDiscount25 := selector.Data("تخفیف 25%", "25discount")
	bot.Handle(&btnDiscount25, func(c tele.Context) error {
		commandLogs("25%discount", c.Sender().ID)
		productId, _ := strconv.Atoi(c.Data())
		db.UpdateProductDiscount(productId, int(c.Sender().ID), 25)
		return c.Reply("با موفقیت انجام شد")
	})

	btnDiscount50 := selector.Data("تخفیف 50%", "50discount")
	bot.Handle(&btnDiscount50, func(c tele.Context) error {
		commandLogs("discount", c.Sender().ID)
		productId, _ := strconv.Atoi(c.Data())
		db.UpdateProductDiscount(productId, int(c.Sender().ID), 50)
		return c.Reply("با موفقیت انجام شد", &tele.SendOptions{
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

func (tlBot TelegramBot) handleQuery() {
	bot := tlBot.bot

	bot.Handle(tele.OnQuery, func(c tele.Context) error {

		//response := search.Query(c.Data())

		products, _ := db.GetAllProductByName(c.Data())
		results := make(tele.Results, len(products))
		for i, product := range products {
			message := service.CreatePreviewMsg(product)

			result := &tele.ArticleResult{
				URL:         product.Image,
				Text:        message,
				Title:       product.Name,
				Description: strconv.Itoa(product.Price),
				ThumbURL:    product.Image,
			}

			results[i] = result
			results[i].SetResultID(strconv.Itoa(i))
			results[i].SetReplyMarkup(getProductUpdateSelector(int(product.Id)))
			results[i].SetParseMode("HTML")
		}

		return c.Answer(&tele.QueryResponse{
			Results:    results,
			CacheTime:  60, // a minute
			IsPersonal: false,
		})
	})

}

func (tlBot TelegramBot) SendUpdateForUsers(productId int, message string, available bool, discount int) {
	pivots := db.GetAllPivotsByProductId(productId)

	rand.Seed(time.Now().UnixNano())
	for _, pivot := range pivots {

		if pivot.Discount > discount {
			continue
		}

		user := db.GetUserById(int64(pivot.UserId))
		msg, _ := tlBot.bot.Send(user.ToTbUser(), message, &tele.SendOptions{
			ParseMode:   "HTML",
			ReplyMarkup: getProductUpdateSelector(productId),
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

func getProductUpdateSelector(productId int) *tele.ReplyMarkup {
	productIdStr := strconv.Itoa(productId)
	selector := &tele.ReplyMarkup{}
	btnGraph := selector.Data("نمودار قیمت", "graph", productIdStr)
	btnDelete := selector.Data("حذف", "delete", productIdStr)

	selector.Inline(
		selector.Row(btnGraph, btnDelete),
	)
	return selector
}

func getProductSelector(productId int) *tele.ReplyMarkup {
	productIdStr := strconv.Itoa(productId)
	selector := &tele.ReplyMarkup{}
	btnDiscount := selector.Data("هر تغییری را اطلاع دهبد", "discount", productIdStr)
	btnDiscount50 := selector.Data("تخفیف50درصد", "50discount", productIdStr)
	btnDiscount25 := selector.Data("تخفیف25درصد", "25discount", productIdStr)
	selector.Inline(
		selector.Row(btnDiscount50, btnDiscount25),
		selector.Row(btnDiscount),
	)
	return selector
}

func commandLogs(command string, userId int64) {
	log.Printf("command: %s, userId: %d", command, userId)
}
