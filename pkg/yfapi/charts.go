package yfapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wcharczuk/go-chart/v2"
)

func (q *Quote) EarningsChart(interval string) ([]byte, error) {
	earnings := make([]chart.Value, 0, 4)

	switch interval {
	case "yearly":
		for _, e := range q.Earnings.Chart.Yearly {
			earnings = append(earnings, chart.Value{Value: e.Earnings.Raw, Label: fmt.Sprintf("%d", int(e.Date))})
		}
	case "quarterly":
		for _, e := range q.Earnings.Chart.Quarterly {
			earnings = append(earnings, chart.Value{Value: e.Earnings.Raw, Label: e.Date})
		}
	}

	graph := createBarChart(fmt.Sprintf("%s %s earnings", q.Price.Symbol, interval), earnings)

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (q *Quote) ChartBytes(param string) (tgbot.FileBytes, error) {
	b, err := q.EarningsChart(param)
	if err != nil {
		return tgbot.FileBytes{}, err
	}
	return tgbot.FileBytes{
		Name:  "charts.png",
		Bytes: b,
	}, nil
}

func createBarChart(name string, data []chart.Value) chart.BarChart {
	return chart.BarChart{
		Title:    name,
		Width:    512,
		Height:   384,
		BarWidth: 40,
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Bars: data,
	}
}

func ChartKeyboard(symbol string) *tgbot.InlineKeyboardMarkup {
	return &tgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbot.InlineKeyboardButton{
			{
				tgbot.NewInlineKeyboardButtonData("Quarterly", fmt.Sprintf("%s|%s", symbol, "quarterly")),
				tgbot.NewInlineKeyboardButtonData("Yearly", fmt.Sprintf("%s|%s", symbol, "yearly")),
			},
		},
	}
}

func NewChartUpdateParams(message *tgbot.Message, symbol string) map[string]string {
	media := struct {
		Type  string `json:"type"`
		Media string `json:"media"`
	}{Type: "photo", Media: "attach://charts.png"}

	mediaJSON, _ := json.Marshal(media)
	kbJSON, _ := json.Marshal(ChartKeyboard(symbol))

	params := map[string]string{
		"chat_id":      strconv.FormatInt(message.Chat.ID, 10),
		"message_id":   strconv.Itoa(message.MessageID),
		"media":        string(mediaJSON),
		"reply_markup": string(kbJSON),
	}

	return params
}
