package yfapi

import (
	"fmt"
	"html"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (q *Quote) WebsiteButton() tgbot.InlineKeyboardButton {
	if q.AssetProfile.Website == "" {
		return tgbot.InlineKeyboardButton{}
	}

	return tgbot.NewInlineKeyboardButtonURL(q.AssetProfile.Website, q.AssetProfile.Website)
}

func (q *Quote) ChartsButton() tgbot.InlineKeyboardButton {
	if len(q.Earnings.Chart.Yearly) == 0 && len(q.Earnings.Chart.Quarterly) == 0 {
		return tgbot.InlineKeyboardButton{}
	}

	return tgbot.NewInlineKeyboardButtonData("Charts", fmt.Sprintf("%s|%s", q.Price.Symbol, "charts"))
}

func (q *Quote) StandardMessageInlineKeyboard() *tgbot.InlineKeyboardMarkup {
	yfURL := fmt.Sprintf("https://finance.yahoo.com/quote/%s", q.Price.Symbol)
	kb := tgbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbot.InlineKeyboardButton{
			{
				tgbot.NewInlineKeyboardButtonURL("Details(Yahoo Finance)", yfURL),
			},
		},
	}

	websiteBtn := q.WebsiteButton()
	if websiteBtn.Text != "" {
		kb.InlineKeyboard[0] = append([]tgbot.InlineKeyboardButton{websiteBtn}, kb.InlineKeyboard[0]...)
	}

	chartsButton := q.ChartsButton()
	if chartsButton.Text != "" {
		kb.InlineKeyboard = append(kb.InlineKeyboard, []tgbot.InlineKeyboardButton{chartsButton})
	}

	return &kb
}

func (q *Quote) StandardMessage() string {
	var msg string
	switch q.Price.Type {
	case "EQUITY":
		msg = fmt.Sprintf("*%s (%s) %s*\n"+
			"_%s_\n\n"+
			"```\n"+
			"MarketCap:      %s\n"+
			"EV:             %s\n"+
			"BV(per share):  %s\n"+
			"Beta:           %s\n\n"+
			"EPS:            %s\n"+
			"P/E:            %s\n"+
			"P/S:            %s\n"+
			"P/B:            %s\n"+
			"Debt/Equity:    %s\n"+
			"Debt/EBITDA:    %s\n\n"+
			"Total Debt:     %s\n"+
			"Total Cash:     %s\n\n"+
			"ROA:            %s\n"+
			"ROE:            %s\n"+
			"FCF:            %s\n"+
			"```",
			q.Name(),
			q.Symbol(),
			q.MarketPrice(),
			q.SectorIndustry(),
			q.MarketCap(),
			q.EnterpriseValue(),
			q.BookValuePerShare(),
			q.Beta(),
			q.EPS(),
			q.PToE(),
			q.PriceToSales(),
			q.PriceToBook(),
			q.DebtToEquity(),
			q.DebtToEBITDA(),
			q.TotalDebt(),
			q.TotalCash(),
			q.ROA(),
			q.ROE(),
			q.FCF(),
		)
	case "ETF":
		msg = fmt.Sprintf("*%s (%s) %s*\n"+
			"_%s_\n\n"+
			"```\n"+
			"Beta(3Y):        %s\n"+
			"Assets:          %s\n"+
			"Expense Ratio:   %s\n"+
			"Return(YTD)      %s\n"+
			"Return(Avg, 3Y)  %s\n"+
			"Return(Avg, 5Y)  %s\n"+
			"```",
			q.Name(),
			q.Symbol(),
			q.MarketPrice(),
			q.Category(),
			q.Beta(),
			q.Assets(),
			q.ExpenseRatio(),
			q.Return("YTD"),
			q.Return("3Y"),
			q.Return("5Y"),
		)
	case "CURRENCY":
		msg = fmt.Sprintf("*%s %s*\n",
			q.Name(),
			q.MarketPrice(),
		)
	}

	return msg
}

func HelpMessage(lang string) string {
	var msg string
	hand := html.UnescapeString("&#" + "128071" + ";")
	switch lang {
	case "ru":
		msg = "Я могу:\n" +
			"- найти базовые показатели и графики компании или фонда по тикеру(например AAPL или VOO)\n" +
			"- найти курс обмена валют (например RUB=X для курса USD/RUB, либо USDRUB=X/RUBUSD=X для конкретной пары)\n" +
			"Список бирж и их суффиксов: [yahoo finance knowledge base](https://help.yahoo.com/kb/exchanges-data-providers-yahoo-finance-sln2310.html)\n\n" +
			"Напиши мне " + hand
	default:
		msg = "I can:\n" +
			"- find basic financial indicators of arbitrary stock symbol(e.g. AAPL or VOO)\n" +
			"- find currency exchange ratio (e.g. RUB=X for USD/RUB pair, or USDRUB=X/RUBUSD=X for specific pair)\n" +
			"Exchanges and data providers list: [yahoo finance knowledge base](https://help.yahoo.com/kb/exchanges-data-providers-yahoo-finance-sln2310.html)\n\n" +
			"Write me " + hand
	}

	return msg
}
