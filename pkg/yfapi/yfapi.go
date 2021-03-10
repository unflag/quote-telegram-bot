package yfapi

import (
	"encoding/json"
	"fmt"
	"time"

	http "github.com/hashicorp/go-retryablehttp"
	"github.com/mitchellh/mapstructure"
)

const (
	apiVersion                 = "v11"
	assetProfileModule         = "assetProfile"
	defaultKeyStatisticsModule = "defaultKeyStatistics"
	earningsModule             = "earnings"
	fundProfileModule          = "fundProfile"
	priceModule                = "price"
	financialDataModule        = "financialData"
)

func GetResponse(symbol string) (QuoteData, error) {
	url := fmt.Sprintf(
		"https://query1.finance.yahoo.com/%s/finance/quoteSummary/%s?modules=%s,%s,%s,%s,%s,%s",
		apiVersion,
		symbol,
		assetProfileModule,
		defaultKeyStatisticsModule,
		earningsModule,
		fundProfileModule,
		priceModule,
		financialDataModule,
	)

	client := http.NewClient()
	client.HTTPClient.Timeout = 3 * time.Second
	client.RetryMax = 5
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	parsedResp := &YFResponse{}
	if err = json.NewDecoder(resp.Body).Decode(parsedResp); err != nil {
		return nil, err
	}

	if parsedResp.Data.Error.Code != "" {
		return nil, &parsedResp.Data.Error
	}

	return parsedResp.Data.Data, err
}

func GetQuote(symbol string) (*Quote, error) {
	data, err := GetResponse(symbol)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, &QuoteError{
			Code:        "Not Found",
			Description: "Symbol contains of invalid symbols",
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
