package tickers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type YahooFinanceResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				RegularMarketPrice   float64 `json:"regularMarketPrice"`
				ChartPreviousClose   float64 `json:"chartPreviousClose"`
				Currency             string  `json:"currency"`
				Symbol               string  `json:"symbol"`
				RegularMarketTime    int64   `json:"regularMarketTime"`
				RegularMarketDayHigh float64 `json:"regularMarketDayHigh"`
				RegularMarketDayLow  float64 `json:"regularMarketDayLow"`
				RegularMarketVolume  float64 `json:"regularMarketVolume"`
			} `json:"meta"`
			Timestamp []int64 `json:"timestamp"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}

func fetchTickerData(ticker string) (*TickerData, error) {
	dayData, err := fetchYahooData(ticker, "1d", "1d")
	if err != nil {
		return nil, err
	}

	weekData, err := fetchYahooData(ticker, "5d", "1d")
	if err != nil {
		return nil, err
	}

	yearData, err := fetchYahooData(ticker, "1y", "1d")
	if err != nil {
		return nil, err
	}

	currentPrice := dayData.Chart.Result[0].Meta.RegularMarketPrice
	dayChange := calculatePercentChange(currentPrice, dayData.Chart.Result[0].Meta.ChartPreviousClose)
	weekChange := calculatePercentChange(currentPrice, weekData.Chart.Result[0].Meta.ChartPreviousClose)
	yearChange := calculatePercentChange(currentPrice, yearData.Chart.Result[0].Meta.ChartPreviousClose)

	return &TickerData{
		Ticker:      ticker,
		TodaysPrice: &currentPrice,
		DayChange:   &dayChange,
		WeekChange:  &weekChange,
		YearChange:  &yearChange,
	}, nil
}

func fetchYahooData(ticker, timeRange, interval string) (*YahooFinanceResponse, error) {
	url := "https://query1.finance.yahoo.com/v8/finance/chart/" + ticker + "?range=" + timeRange + "&interval=" + interval

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data YahooFinanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if data.Chart.Error != nil {
		return nil, fmt.Errorf("API error: %v", data.Chart.Error)
	}

	if len(data.Chart.Result) == 0 {
		return nil, fmt.Errorf("no data returned for ticker %s", ticker)
	}

	return &data, nil
}

func calculatePercentChange(current, previous float64) float64 {
	if previous == 0 {
		return 0
	}
	return ((current - previous) / previous) * 100
}
