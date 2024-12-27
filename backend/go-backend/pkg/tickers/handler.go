package tickers

import (
	"log"
	"sort"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Get("/tickers", h.GetTickers)
}

func (h *Handler) GetTickers(c *fiber.Ctx) error {
	if data, ok := getCachedData(); ok {
		log.Printf("[Stocks] Returning data from cache")
		sortTickersByDayChange(data)
		return c.JSON(data)
	}

	log.Printf("[Stocks] Cache miss. Fetching data from API....")
	results := make(chan *TickerData, len(DefaultTickers))
	errors := make(chan error, len(DefaultTickers))

	var wg sync.WaitGroup
	wg.Add(len(DefaultTickers))

	for _, ticker := range DefaultTickers {
		go func(t string) {
			defer wg.Done()
			data, err := fetchTickerData(t)
			if err != nil {
				errors <- err
				results <- nil
				return
			}
			results <- data
		}(ticker)
	}

	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	tickerData := make([]TickerData, 0, len(DefaultTickers))
	for data := range results {
		if data != nil {
			tickerData = append(tickerData, *data)
		}
	}

	var errs []error
	for err := range errors {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(tickerData) == 0 && len(errs) > 0 {
		return fiber.NewError(fiber.StatusInternalServerError, errs[0].Error())
	}

	if len(tickerData) > 0 {
		sortTickersByDayChange(tickerData)
		updateCache(tickerData)
	}

	return c.JSON(tickerData)
}

func sortTickersByDayChange(data []TickerData) {
	sort.Slice(data, func(i, j int) bool {
		if data[i].DayChange == nil {
			return false
		}
		if data[j].DayChange == nil {
			return true
		}
		// Sort by absolute day change percentage descending
		// This will show biggest movers (both up and down) first
		return abs(*data[i].DayChange) > abs(*data[j].DayChange)
	})
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
