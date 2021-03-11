package yfapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wcharczuk/go-chart/v2"
)

type ChartParams struct {
	Symbol      string
	Interval    string
	Measurement string
	Cmd         string
}

func NewChartParams(callbackData string) (*ChartParams, error) {
	data := strings.Split(callbackData, "|")
	minLen := 4
	if len(data) < minLen {
		return nil, fmt.Errorf("provided data is too short(%d < %d): %s", len(data), minLen, callbackData)
	}

	return &ChartParams{
		Symbol:      data[0],
		Interval:    data[1],
		Measurement: data[2],
		Cmd:         data[3],
	}, nil
}

func (q *Quote) earningsChart(p *ChartParams) ([]byte, error) {
	data := make([]chart.Value, 0, 4)

	switch p.Interval {
	case "yearly":
		for _, e := range q.Earnings.Chart.Yearly {
			data = append(data, chart.Value{Value: e.Value(p.Measurement), Label: fmt.Sprintf("%d", int(e.Date))})
		}
	case "quarterly":
		for _, e := range q.Earnings.Chart.Quarterly {
			data = append(data, chart.Value{Value: e.Value(p.Measurement), Label: e.Date})
		}
	}

	graph := createBarChart(fmt.Sprintf("%s %s %s", q.Price.Symbol, p.Interval, p.Measurement), data)

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (q *Quote) ChartBytes(p *ChartParams) (tgbot.FileBytes, error) {
	b, err := q.earningsChart(p)
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

func ChartKeyboard(p *ChartParams) *tgbot.InlineKeyboardMarkup {
	return &tgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbot.InlineKeyboardButton{
			{
				tgbot.NewInlineKeyboardButtonData("Earnings",
					fmt.Sprintf("%s|%s|%s|%s",
						p.Symbol,
						p.Interval,
						"earnings",
						"update",
					),
				),
				tgbot.NewInlineKeyboardButtonData("Revenue",
					fmt.Sprintf("%s|%s|%s|%s",
						p.Symbol,
						p.Interval,
						"revenue",
						"update",
					),
				),
			},
			{
				tgbot.NewInlineKeyboardButtonData("Quarterly",
					fmt.Sprintf("%s|%s|%s|%s",
						p.Symbol,
						"quarterly",
						p.Measurement,
						"update",
					),
				),
				tgbot.NewInlineKeyboardButtonData("Yearly",
					fmt.Sprintf("%s|%s|%s|%s",
						p.Symbol,
						"yearly",
						p.Measurement,
						"update",
					),
				),
			},
		},
	}
}

func NewChartUpdateParams(message *tgbot.Message, p *ChartParams) map[string]string {
	media := struct {
		Type  string `json:"type"`
		Media string `json:"media"`
	}{Type: "photo", Media: "attach://charts.png"}

	mediaJSON, _ := json.Marshal(media)
	kbJSON, _ := json.Marshal(ChartKeyboard(p))

	params := map[string]string{
		"chat_id":      strconv.FormatInt(message.Chat.ID, 10),
		"message_id":   strconv.Itoa(message.MessageID),
		"media":        string(mediaJSON),
		"reply_markup": string(kbJSON),
	}

	return params
}
