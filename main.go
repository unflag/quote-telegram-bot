package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"quote-telegram-bot/pkg/helpers"
	"quote-telegram-bot/pkg/yfapi"
)

var (
	botToken string
	debug    bool
)

func init() {
	debugEnv, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	flag.StringVar(&botToken, "botToken", os.Getenv("BOT_TOKEN"), "Telegram bot token")
	flag.BoolVar(&debug, "debug", debugEnv, "Enable debug")
}

func main() {
	bot, err := tgbot.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}

	bot.Debug = debug
	updateConfig := tgbot.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.CallbackQuery != nil {
			data := strings.Split(update.CallbackQuery.Data, "|")
			if len(data) < 2 {
				continue
			}
			symbol, cmd := data[0], data[1]

			quote, err := yfapi.GetQuote(symbol)
			if err != nil {
				log.Println(err)
				continue
			}

			switch cmd {
			case "charts":
				chart, err := quote.ChartBytes("quarterly")
				if err != nil {
					log.Println(err)
					continue
				}
				graph := tgbot.NewPhotoUpload(update.CallbackQuery.Message.Chat.ID, chart)
				graph.ReplyMarkup = yfapi.ChartKeyboard(symbol)
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
				chart, err := quote.ChartBytes(cmd)
				if err != nil {
					log.Println(err)
					continue
				}
				params := yfapi.NewChartUpdateParams(update.CallbackQuery.Message, symbol)
				err = helpers.Retry(3, func() error {
					if _, err = bot.UploadFile("editMessageMedia", params, "charts.png", chart); err != nil {
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

		switch update.Message.IsCommand() {
		case true:
			msg.Text = yfapi.HelpMessage(update.Message.From.LanguageCode)
		default:
			text := helpers.Sanitize(update.Message.Text)
			quote, err := yfapi.GetQuote(text)
			if err != nil {
				msg.Text = fmt.Sprintf("Unable to get data for symbol: %s", text)
				log.Println(err)
				if qerr, ok := err.(*yfapi.QuoteError); ok {
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
