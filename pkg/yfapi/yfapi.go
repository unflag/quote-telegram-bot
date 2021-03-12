package yfapi

import (
	"encoding/json"
	"fmt"
	"time"

	http "github.com/hashicorp/go-retryablehttp"
	"github.com/mitchellh/mapstructure"
)

const (
	quotesApiVersion           = "v11"
	assetProfileModule         = "assetProfile"
	defaultKeyStatisticsModule = "defaultKeyStatistics"
	earningsModule             = "earnings"
	fundProfileModule          = "fundProfile"
	priceModule                = "price"
	financialDataModule        = "financialData"

	chartsApiVersion = "v8"
	chartMeta        = "meta"
	chartTimestamps  = "timestamp"
	chartIndicators  = "indicators"
)

type YFClient struct {
	*http.Client
}

func NewYFClient() *YFClient {
	client := http.NewClient()
	client.HTTPClient.Timeout = 3 * time.Second
	client.RetryMax = 5

	return &YFClient{
		client,
	}
}

func (c *YFClient) getQuoteResponse(symbol string) (QuoteData, error) {
	url := fmt.Sprintf(
		"https://query1.finance.yahoo.com/%s/finance/quoteSummary/%s?modules=%s,%s,%s,%s,%s,%s",
		quotesApiVersion,
		symbol,
		assetProfileModule,
		defaultKeyStatisticsModule,
		earningsModule,
		fundProfileModule,
		priceModule,
		financialDataModule,
	)

	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	parsedResp := &QuoteResponse{}
	if err = json.NewDecoder(resp.Body).Decode(parsedResp); err != nil {
		return nil, err
	}

	if parsedResp.Data.Error.Code != "" {
		return nil, &parsedResp.Data.Error
	}

	return parsedResp.Data.Data, nil
}

func (c *YFClient) GetQuote(symbol string) (*Quote, error) {
	data, err := c.getQuoteResponse(symbol)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, &QueryError{
			Code:        "Not Found",
			Description: "Symbol consists of invalid characters",
		}
	}

	quote := Quote{}
	for k, v := range data[0] {
		switch k {
		case defaultKeyStatisticsModule:
			if err = mapstructure.Decode(v, &quote.Statistics); err != nil {
				return nil, err
			}
		case assetProfileModule:
			if err = mapstructure.Decode(v, &quote.AssetProfile); err != nil {
				return nil, err
			}
		case fundProfileModule:
			if err = mapstructure.Decode(v, &quote.FundProfile); err != nil {
				return nil, err
			}
		case earningsModule:
			if err = mapstructure.Decode(v, &quote.Earnings); err != nil {
				return nil, err
			}
		case financialDataModule:
			if err = mapstructure.Decode(v, &quote.Financials); err != nil {
				return nil, err
			}
		case priceModule:
			if err = mapstructure.Decode(v, &quote.Price); err != nil {
				return nil, err
			}
		}
	}

	return &quote, nil
}

func (c *YFClient) getChartResponse(symbol string, period string) (ChartData, error) {
	var interval string
	switch period {
	case "1d":
		interval = "1h"
	case "1mo":
		interval = "1d"
	case "3mo":
		interval = "5d"
	case "6mo":
		interval = "5d"
	case "1y":
		interval = "1mo"
	default:
		period = "1d"
		interval = "1h"
	}

	url := fmt.Sprintf("https://query1.finance.yahoo.com/%s/finance/chart/%s?period1=0&period2=9999999999&interval=%s&range=%s",
		chartsApiVersion,
		symbol,
		interval,
		period,
	)

	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	parsedResp := &ChartResponse{}
	if err = json.NewDecoder(resp.Body).Decode(parsedResp); err != nil {
		return nil, err
	}

	if parsedResp.Data.Error.Code != "" {
		return nil, &parsedResp.Data.Error
	}

	return parsedResp.Data.Data, nil
}

func (c *YFClient) GetChart(symbol string, period string) (*Chart, error) {
	data, err := c.getChartResponse(symbol, period)
	if err != nil {
		return nil, err
	}

	chart := Chart{}
	for k, v := range data[0] {
		switch k {
		case chartMeta:
			if err = mapstructure.Decode(v, &chart.Meta); err != nil {
				return nil, err
			}
		case chartIndicators:
			if err = mapstructure.Decode(v, &chart.Indicators); err != nil {
				return nil, err
			}
		case chartTimestamps:
			if err = mapstructure.Decode(v, &chart.Timestamps); err != nil {
				return nil, err
			}
		}
	}

	return &chart, nil
}
