package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"quote-telegram-bot/pkg/helpers"
	"quote-telegram-bot/pkg/yfapi"
)

var (
	botToken string
	debug    bool
	version  bool

	Name    string
	Version string
	Date    string
)

func init() {
	debugEnv, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	flag.StringVar(&botToken, "botToken", os.Getenv("BOT_TOKEN"), "Telegram bot token")
	flag.BoolVar(&debug, "debug", debugEnv, "Enable debug")
	flag.BoolVar(&version, "version", false, "Print version")
	flag.Parse()
}

func main() {
	if version {
		fmt.Printf("%s %s %s", Name, Version, Date)
		return
	}

	bot, err := tgbot.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}

	bot.Debug = debug
	updateConfig := tgbot.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	yfc := yfapi.NewYFClient()

	for update := range updates {
		// update.CallbackQuery used to process buttons and sending back charts
		if update.CallbackQuery != nil {
			// generate button press metadata
			params, err := yfapi.NewChartParams(update.CallbackQuery.Data)
			if err != nil {
				log.Println(err)
				continue
			}

			// price and earnings/revenue charts have different sources and formats, but same interface
			var data yfapi.Chartable
			switch params.Measurement {
			case "price":
				data, err = yfc.GetPriceChart(params.Symbol, params.Interval)
			case "earnings", "revenue":
				data, err = yfc.GetQuote(params.Symbol)
			default:
				continue
			}

			if err != nil {
				log.Println(err)
				continue
			}

			switch params.Cmd {
			// initial received only once when user press button "Charts" under quote info message.
			// Price chart sent on this event.
			case "initial":
				chart, err := data.ChartBytes(params)
				if err != nil {
					log.Println(err)
					continue
				}
				graph := tgbot.NewPhotoUpload(update.CallbackQuery.Message.Chat.ID, chart)
				graph.ReplyMarkup = yfapi.ChartKeyboard(params, data.Intervals())
				err = helpers.Retry(3, func() error {
					if _, err := bot.Send(graph); err != nil {
						return err
					}
					return nil
				})
				if err != nil {
					log.Println(err)
				}
			default:
				// any chart updates are processing here
				chart, err := data.ChartBytes(params)
				if err != nil {
					log.Println(err)
					continue
				}
				p := yfapi.NewMediaUpdateParams(update.CallbackQuery.Message, params, data.Intervals())
				err = helpers.Retry(3, func() error {
					if _, err = bot.UploadFile("editMessageMedia", p, "charts.png", chart); err != nil {
						return err
					}
					return nil
				})
				if err != nil {
					log.Println(err)
				}
			}

			continue
		}

		msg := tgbot.NewMessage(update.Message.Chat.ID, "")
		msg.ParseMode = tgbot.ModeMarkdown
		msg.DisableWebPagePreview = true

		// text messages from user are processing here
		switch update.Message.IsCommand() {
		// bot has no user commands, so any command causes help message send
		case true:
			msg.Text = yfapi.HelpMessage(update.Message.From.LanguageCode)
		// any text messages are processing here
		default:
			text := helpers.Sanitize(update.Message.Text)
			quote, err := yfc.GetQuote(text)
			if err != nil {
				msg.Text = fmt.Sprintf("Unable to get data for symbol: %s", text)
				log.Println(err)
				if qerr, ok := err.(*yfapi.QueryError); ok {
					msg.Text = qerr.Error()
				}
			} else {
				msg.Text = quote.StandardMessage()
				msg.ReplyMarkup = quote.StandardMessageInlineKeyboard()
			}

			if msg.Text == "" {
				msg.Text = fmt.Sprintf("No data found for symbol: %s", text)
			}
		}

		err := helpers.Retry(3, func() error {
			if _, err := bot.Send(msg); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			log.Println(err)
		}
	}
}
