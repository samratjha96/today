package tickers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

// UserAgents list to rotate through when making requests
var UserAgents = []string{
	// Chrome
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",

	// Firefox
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:135.0) Gecko/20100101 Firefox/135.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14.7; rv:135.0) Gecko/20100101 Firefox/135.0",
	"Mozilla/5.0 (X11; Linux i686; rv:135.0) Gecko/20100101 Firefox/135.0",

	// Safari
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_7_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.3 Safari/605.1.15",

	// Edge
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36 Edg/131.0.2903.86",
}

type UserAgent struct {
	sync.Once
	name string
}

var agent UserAgent

func selectAgent() {
	agent.name = UserAgents[rand.Intn(len(UserAgents))]
}

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
	// Get day data
	dayData, err := fetchYahooData(ticker, "1d", "1d")
	if err != nil {
		return nil, err
	}

	// Get week data
	weekData, err := fetchYahooData(ticker, "5d", "1d")
	if err != nil {
		return nil, err
	}

	// Get year data
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

	// Log the URL being requested
	log.Printf("[Stocks] Requesting URL: %s", url)

	// Create a new request with User-Agent header
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Initialize the User-Agent once and reuse it
	agent.Do(func() { selectAgent() })
	req.Header.Add("User-Agent", agent.name)

	// Log the User-Agent being used
	log.Printf("[Stocks] Using User-Agent: %s", agent.name)

	// Make the request with the configured headers
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		// If we hit rate limit, extend cache time
		log.Printf("[Stocks] Rate limited by Yahoo API. Extending cache time.")
		ExtendCacheTime(30 * time.Minute)

		// Sleep before retry
		time.Sleep(2 * time.Second)

		// Log retry attempt with URL
		log.Printf("[Stocks] Retrying request to: %s", url)

		// Create a new request for the retry with possibly a different User-Agent
		retryReq, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create retry request: %v", err)
		}

		// Select a different User-Agent for the retry
		agent.name = UserAgents[rand.Intn(len(UserAgents))]
		retryReq.Header.Add("User-Agent", agent.name)
		log.Printf("[Stocks] Retry using User-Agent: %s", agent.name)

		// Try again with the new request
		resp, err = http.DefaultClient.Do(retryReq)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch data after retry: %v", err)
		}
		defer resp.Body.Close()
	}

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
