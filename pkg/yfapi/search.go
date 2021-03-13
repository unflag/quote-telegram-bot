package yfapi

import (
	"fmt"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	searchTypes = map[string]struct{}{
		"S": {},
		"E": {},
		"C": {},
	}
)

type SearchResponse struct {
	Data ResultSet `json:"ResultSet"`
}

type ResultSet struct {
	Query  string         `json:"Query"`
	Result []SearchResult `json:"Result"`
}

type SearchResult struct {
	Symbol   string
	Name     string `json:"name"`
	Exch     string `json:"exch"`
	Type     string `json:"type"`
	ExchDisp string `json:"exchDisp"`
	TypeDisp string `json:"typeDisp"`
}

func (r *ResultSet) SearchMessage() string {
	validResultLen := 0
	for _, res := range r.Result {
		if _, ok := searchTypes[res.Type]; ok {
			validResultLen++
		}
	}

	if validResultLen == 0 {
		return fmt.Sprintf("Not Found: Quote not found for search message: %s", r.Query)
	}

	return fmt.Sprintf("Search result for: %s", r.Query)
}

func (r *ResultSet) SearchMessageInlineKeyboard() *tgbot.InlineKeyboardMarkup {
	rows := make([][]tgbot.InlineKeyboardButton, 0, len(r.Result)/4)
	buttons := make([]tgbot.InlineKeyboardButton, 0, len(r.Result))
	for i, res := range r.Result {
		if _, ok := searchTypes[res.Type]; ok {
			buttons = append(buttons,
				tgbot.NewInlineKeyboardButtonData(
					res.Symbol,
					res.Symbol,
				),
			)
		}
		if len(buttons) == 4 || i == len(r.Result)-1 {
			rows = append(rows, buttons)
			buttons = make([]tgbot.InlineKeyboardButton, 0, len(r.Result))
		}
	}

	return &tgbot.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}
