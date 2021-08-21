package yfapi

import (
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	searchTypes = map[string]struct{}{
		"EQUITY":         {},
		"ETF":            {},
		"CURRENCY":       {},
		"CRYPTOCURRENCY": {},
	}
)

type SearchResponse struct {
	Result []SearchResult `json:"quotes"`
}

type SearchResult struct {
	Symbol   string `json:"symbol"`
	Name     string `json:"shortname"`
	Exchange string `json:"exchange"`
	Type     string `json:"quoteType"`
}

func (r *SearchResponse) SearchMessage() string {
	validResultLen := 0
	for _, res := range r.Result {
		if _, ok := searchTypes[res.Type]; ok {
			validResultLen++
		}
	}

	if validResultLen == 0 {
		return "Quote not found"
	}

	return "Search result:"
}

func (r *SearchResponse) SearchMessageInlineKeyboard() *tgbot.InlineKeyboardMarkup {
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
