package tickers

import (
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
		return c.JSON(data)
	}

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
		updateCache(tickerData)
	}

	return c.JSON(tickerData)
}
