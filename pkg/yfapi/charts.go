package yfapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/wcharczuk/go-chart/v2"
)

type Chart struct {
	Meta       ChartMeta       `mapstructure:"meta"`
	Indicators ChartIndicators `mapstructure:"indicators"`
	Timestamps []int           `mapstructure:"timestamp"`
}

type ChartResponse struct {
	Data ChartSummary `json:"chart"`
}

type ChartSummary struct {
	Data  ChartData  `json:"result"`
	Error QueryError `json:"error"`
}

type ChartData = []map[string]interface{}

type ChartMeta struct {
	Currency        string   `mapstructure:"currency"`
	Symbol          string   `mapstructure:"symbol"`
	DataGranularity string   `mapstructure:"dataGranularity"`
	Range           string   `mapstructure:"range"`
	ValidIntervals  []string `mapstructure:"validRanges"`
}

type ChartIndicators struct {
	Quote []ChartQuote `mapstructure:"quote"`
}

type ChartQuote struct {
	High   []float64 `mapstructure:"high"`
	Low    []float64 `mapstructure:"low"`
	Open   []float64 `mapstructure:"open"`
	Close  []float64 `mapstructure:"close"`
	Volume []int     `mapstructure:"volume"`
}

type ChartParams struct {
	Symbol      string
	Interval    string
	Measurement string
	Cmd         string
	Type        string
}

type Chartable interface {
	ChartBytes(p *ChartParams) (tgbot.FileBytes, error)
	Intervals() []string
}

func NewChartParams(callbackData string) (*ChartParams, error) {
	data := strings.Split(callbackData, "|")
	minLen := 5
	if len(data) < minLen {
		return nil, fmt.Errorf("provided data is too short(%d < %d): %s", len(data), minLen, callbackData)
	}

	return &ChartParams{
		Symbol:      data[0],
		Interval:    data[1],
		Measurement: data[2],
		Cmd:         data[3],
		Type:        data[4],
	}, nil
}

func NewMediaUpdateParams(message *tgbot.Message, p *ChartParams, i []string) map[string]string {
	media := struct {
		Type  string `json:"type"`
		Media string `json:"media"`
	}{Type: "photo", Media: "attach://charts.png"}

	mediaJSON, _ := json.Marshal(media)
	kbJSON, _ := json.Marshal(ChartKeyboard(p, i))

	updateParams := map[string]string{
		"chat_id":      strconv.FormatInt(message.Chat.ID, 10),
		"message_id":   strconv.Itoa(message.MessageID),
		"media":        string(mediaJSON),
		"reply_markup": string(kbJSON),
	}

	return updateParams
}

func (c *Chart) Intervals() []string {
	intervals := make([]string, 0, len(c.Meta.ValidIntervals))
	for _, interval := range c.Meta.ValidIntervals {
		if _, ok := priceIntervals[interval]; ok {
			intervals = append(intervals, interval)
		}
	}

	return intervals
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

func (c *Chart) ChartBytes(p *ChartParams) (tgbot.FileBytes, error) {
	b, err := c.priceChart(p)
	if err != nil {
		return tgbot.FileBytes{}, err
	}

	return tgbot.FileBytes{
		Name:  "charts.png",
		Bytes: b,
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

func (c *Chart) priceChart(p *ChartParams) ([]byte, error) {
	dates := make([]time.Time, 0, len(c.Timestamps))
	for _, ts := range c.Timestamps {
		dates = append(dates, time.Unix(int64(ts), 0))
	}

	graph := createTSChart(fmt.Sprintf("%s %s (%s)", p.Symbol, p.Measurement, p.Interval),
		dates,
		c.Indicators.Quote[0].High,
	)

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func createTSChart(name string, x []time.Time, y []float64) chart.Chart {
	return chart.Chart{
		Title:  name,
		Height: 512,
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		XAxis: chart.XAxis{
			TickPosition: chart.TickPositionUnderTick,
			TickStyle: chart.Style{
				FontSize: 12,
			},
		},
		YAxis: chart.YAxis{
			TickStyle: chart.Style{
				FontSize: 14,
			},
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name: name,
				Style: chart.Style{
					StrokeColor: chart.ColorBlue,
					StrokeWidth: 2.5,
					FillColor:   chart.ColorLightGray,
				},
				XValues: x,
				YValues: y,
			},
		},
	}

}

func ChartKeyboard(p *ChartParams, i []string) *tgbot.InlineKeyboardMarkup {
	firstRow := []tgbot.InlineKeyboardButton{
		tgbot.NewInlineKeyboardButtonData("Price",
			fmt.Sprintf("%s|%s|%s|%s|%s",
				p.Symbol,
				"1d",
				"price",
				"update",
				p.Type,
			),
		),
		tgbot.NewInlineKeyboardButtonData("Earnings",
			fmt.Sprintf("%s|%s|%s|%s|%s",
				p.Symbol,
				"quarterly",
				"earnings",
				"update",
				p.Type,
			),
		),
		tgbot.NewInlineKeyboardButtonData("Revenue",
			fmt.Sprintf("%s|%s|%s|%s|%s",
				p.Symbol,
				"quarterly",
				"revenue",
				"update",
				p.Type,
			),
		),
	}

	kb := make([][]tgbot.InlineKeyboardButton, 0, 2)
	if p.Type == "hasEarnings" {
		kb = append(kb, firstRow)
	}

	kb = append(kb, chartKeyboardSecondRow(p, i))

	return &tgbot.InlineKeyboardMarkup{
		InlineKeyboard: kb,
	}
}

func chartKeyboardSecondRow(p *ChartParams, intervals []string) []tgbot.InlineKeyboardButton {
	row := make([]tgbot.InlineKeyboardButton, 0, len(intervals))

	for _, interval := range intervals {
		row = append(row,
			tgbot.NewInlineKeyboardButtonData(
				interval,
				fmt.Sprintf("%s|%s|%s|%s|%s",
					p.Symbol,
					interval,
					p.Measurement,
					"update",
					p.Type,
				),
			),
		)
	}

	if len(row) <= 1 {
		return []tgbot.InlineKeyboardButton{}
	}

	return row
}
