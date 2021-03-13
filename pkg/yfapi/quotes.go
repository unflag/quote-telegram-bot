package yfapi

import (
	"fmt"
)

type Quote struct {
	Price        QuotePrice
	AssetProfile QuoteAssetProfile
	FundProfile  QuoteFundProfile
	Statistics   QuoteStatistics
	Financials   QuoteFinancials
	Earnings     QuoteEarnings
}

type QuoteResponse struct {
	Data QuoteSummary `json:"quoteSummary"`
}

type QuoteSummary struct {
	Data  QuoteData  `json:"result"`
	Error QueryError `json:"error"`
}

type QuoteData = []map[string]map[string]interface{}

type QueryError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

func (e *QueryError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Description)
}

type IndicatorValue struct {
	Raw float64 `mapstructure:"raw"`
	Fmt string  `mapstructure:"fmt"`
}

// https://query1.finance.yahoo.com/v11/finance/quoteSummary/${QUOTE}?modules=price
type QuotePrice struct {
	Symbol         string         `mapstructure:"symbol"`
	Name           string         `mapstructure:"shortName"`
	Type           string         `mapstructure:"quoteType"`
	Currency       string         `mapstructure:"currency"`
	CurrencySymbol string         `mapstructure:"currencySymbol"`
	Exchange       string         `mapstructure:"exchangeName"`
	MarketCap      IndicatorValue `mapstructure:"marketCap"`
	MarketPrice    IndicatorValue `mapstructure:"regularMarketPrice"`
}

// https://query1.finance.yahoo.com/v11/finance/quoteSummary/${QUOTE}?modules=assetProfile
type QuoteAssetProfile struct {
	Sector   string `mapstructure:"sector"`
	Industry string `mapstructure:"industry"`
	Website  string `mapstructure:"website"`
}

type QuoteFundProfile struct {
	Fees QuoteFundFees `mapstructure:"feesExpensesInvestment"`
}

type QuoteFundFees struct {
	ExpenseRatio IndicatorValue `mapstructure:"annualReportExpenseRatio"`
}

// https://query1.finance.yahoo.com/v11/finance/quoteSummary/${QUOTE}?modules=defaultKeyStatistics
type QuoteStatistics struct {
	EV           IndicatorValue `mapstructure:"enterpriseValue"`
	BV           IndicatorValue `mapstructure:"bookValue"`
	Beta         IndicatorValue `mapstructure:"beta"`
	Beta3Y       IndicatorValue `mapstructure:"beta3Year"`
	ReturnYTD    IndicatorValue `mapstructure:"ytdReturn"`
	Return3Y     IndicatorValue `mapstructure:"threeYearAverageReturn"`
	Return5Y     IndicatorValue `mapstructure:"fiveYearAverageReturn"`
	PriceToBook  IndicatorValue `mapstructure:"priceToBook"`
	PriceToSales IndicatorValue `mapstructure:"priceToSalesTrailing12Months"`
	EVToEBITDA   IndicatorValue `mapstructure:"enterpriseToEbitda"`
	EPS          IndicatorValue `mapstructure:"trailingEps"`
	Assets       IndicatorValue `mapstructure:"totalAssets"`
	Shares       IndicatorValue `mapstructure:"sharesOutstanding"`
	Category     string         `mapstructure:"category"`
}

// https://query1.finance.yahoo.com/v11/finance/quoteSummary/${QUOTE}?modules=financialData
type QuoteFinancials struct {
	Recommendation           string         `mapstructure:"recommendationKey"`
	NumberOfAnalystsOpinions IndicatorValue `mapstructure:"numberOfAnalystOpinions"`
	TotalCash                IndicatorValue `mapstructure:"totalCash"`
	EBITDA                   IndicatorValue `mapstructure:"ebitda"`
	TotalDebt                IndicatorValue `mapstructure:"totalDebt"`
	DebtToEquity             IndicatorValue `mapstructure:"debtToEquity"`
	ReturnOnAssets           IndicatorValue `mapstructure:"returnOnAssets"`
	ReturnOnEquity           IndicatorValue `mapstructure:"returnOnEquity"`
	FCF                      IndicatorValue `mapstructure:"freeCashflow"`
	Currency                 string         `mapstructure:"financialCurrency"`
}

// https://query1.finance.yahoo.com/v11/finance/quoteSummary/${QUOTE}?modules=earnings
type QuoteEarnings struct {
	Chart    FinancialsChart `mapstructure:"financialsChart"`
	Currency string          `mapstructure:"financialCurrency"`
}

type FinancialsChart struct {
	Quarterly []QuarterlyFinancialsChart `mapstructure:"quarterly"`
	Yearly    []YearlyFinancialsChart    `mapstructure:"yearly"`
}

type QuarterlyFinancialsChart struct {
	Date     string         `mapstructure:"date"`
	Revenue  IndicatorValue `mapstructure:"revenue"`
	Earnings IndicatorValue `mapstructure:"earnings"`
}

type YearlyFinancialsChart struct {
	Date     float64        `mapstructure:"date"`
	Revenue  IndicatorValue `mapstructure:"revenue"`
	Earnings IndicatorValue `mapstructure:"earnings"`
}

func (qc *QuarterlyFinancialsChart) Value(measurement string) float64 {
	var value float64
	switch measurement {
	case "earnings":
		value = qc.Earnings.Raw
	case "revenue":
		value = qc.Revenue.Raw
	}

	return value
}

func (yc *YearlyFinancialsChart) Value(measurement string) float64 {
	var value float64
	switch measurement {
	case "earnings":
		value = yc.Earnings.Raw
	case "revenue":
		value = yc.Revenue.Raw
	}

	return value
}

func (q *Quote) Name() string {
	return q.Price.Name
}

func (q *Quote) Symbol() string {
	return q.Price.Symbol
}

func (q *Quote) MarketPrice() string {
	if q.Price.MarketPrice.Fmt == "" {
		return "N/A"
	}

	return q.Price.CurrencySymbol + q.Price.MarketPrice.Fmt
}

func (q *Quote) MarketCap() string {
	if q.Price.MarketCap.Fmt == "" {
		return "N/A"
	}

	return q.Price.CurrencySymbol + q.Price.MarketCap.Fmt
}

func (q *Quote) EnterpriseValue() string {
	if q.Statistics.EV.Fmt == "" {
		return "N/A"
	}

	return q.Price.CurrencySymbol + q.Statistics.EV.Fmt
}

func (q *Quote) BookValuePerShare() string {
	if q.Statistics.BV.Fmt == "" {
		return "N/A"
	}

	return q.Price.CurrencySymbol + q.Statistics.BV.Fmt
}

func (q *Quote) PriceToSales() string {
	if q.Statistics.PriceToSales.Fmt == "" {
		return "N/A"
	}

	return q.Statistics.PriceToSales.Fmt
}

func (q *Quote) PriceToBook() string {
	if q.Statistics.PriceToBook.Fmt == "" {
		return "N/A"
	}

	return q.Statistics.PriceToBook.Fmt
}

func (q *Quote) EnterpriseValueToEBITDA() string {
	if q.Statistics.EVToEBITDA.Fmt == "" {
		return "N/A"
	}

	return q.Statistics.EVToEBITDA.Fmt
}

func (q *Quote) DebtToEquity() string {
	if q.Financials.DebtToEquity.Fmt == "" {
		return "N/A"
	}

	return q.Financials.DebtToEquity.Fmt
}

func (q *Quote) DebtToEBITDA() string {
	if q.Financials.TotalDebt.Fmt == "" || q.Financials.EBITDA.Fmt == "" {
		return "N/A"
	}

	dte := q.Financials.TotalDebt.Raw / q.Financials.EBITDA.Raw
	if dte < 0 {
		return "N/A"
	}

	return fmt.Sprintf("%.2f", dte)
}

func (q *Quote) TotalDebt() string {
	if q.Financials.TotalDebt.Fmt == "" {
		return "N/A"
	}

	return q.Financials.TotalDebt.Fmt + " (" + q.Financials.Currency + ")"
}

func (q *Quote) TotalCash() string {
	if q.Financials.TotalCash.Fmt == "" {
		return "N/A"
	}

	return q.Financials.TotalCash.Fmt + " (" + q.Financials.Currency + ")"
}

func (q *Quote) EPS() string {
	if q.Statistics.EPS.Fmt == "" {
		return "N/A"
	}

	return q.Price.CurrencySymbol + q.Statistics.EPS.Fmt
}

func (q *Quote) ROA() string {
	if q.Financials.ReturnOnAssets.Fmt == "" {
		return "N/A"
	}

	return q.Financials.ReturnOnAssets.Fmt
}

func (q *Quote) ROE() string {
	if q.Financials.ReturnOnEquity.Fmt == "" {
		return "N/A"
	}

	return q.Financials.ReturnOnEquity.Fmt
}

func (q *Quote) FCF() string {
	if q.Financials.FCF.Fmt == "" {
		return "N/A"
	}

	return q.Financials.FCF.Fmt + " (" + q.Financials.Currency + ")"
}

func (q *Quote) PToE() string {
	if q.Price.MarketPrice.Fmt == "" || q.Statistics.EPS.Fmt == "" {
		return "N/A"
	}

	pe := q.Price.MarketPrice.Raw / q.Statistics.EPS.Raw
	if pe < 0 {
		return "N/A"
	}

	return fmt.Sprintf("%.2f", pe)
}

func (q *Quote) Beta() string {
	if q.Statistics.Beta.Fmt == "" && q.Statistics.Beta3Y.Fmt == "" {
		return "N/A"
	}

	var beta string
	switch q.Price.Type {
	case "EQUITY":
		beta = q.Statistics.Beta.Fmt
	case "ETF":
		beta = q.Statistics.Beta3Y.Fmt
	default:
		beta = "N/A"
	}

	return beta
}

func (q *Quote) Assets() string {
	if q.Statistics.Assets.Fmt == "" {
		return "N/A"
	}

	return q.Statistics.Assets.Fmt
}

func (q *Quote) ExpenseRatio() string {
	if q.FundProfile.Fees.ExpenseRatio.Fmt == "" {
		return "N/A"
	}

	return q.FundProfile.Fees.ExpenseRatio.Fmt
}

func (q *Quote) Return(period string) string {
	var ret string
	switch period {
	case "YTD":
		ret = q.Statistics.ReturnYTD.Fmt
	case "3Y":
		ret = q.Statistics.Return3Y.Fmt
	case "5Y":
		ret = q.Statistics.Return5Y.Fmt
	}

	if ret == "" {
		ret = "N/A"
	}

	return ret
}

func (q *Quote) SectorIndustry() string {
	if q.AssetProfile.Sector == "" || q.AssetProfile.Industry == "" {
		return "Unknown sector/industry"
	}

	return q.AssetProfile.Sector + " (" + q.AssetProfile.Industry + ")"
}

func (q *Quote) Category() string {
	if q.Statistics.Category == "" {
		return "Unknown category"
	}

	return q.Statistics.Category
}

func (q *Quote) Intervals() []string {
	return []string{"quarterly", "yearly"}
}
