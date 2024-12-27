package tickers

type TickerData struct {
	Ticker      string   `json:"ticker"`
	TodaysPrice *float64 `json:"todaysPrice"`
	DayChange   *float64 `json:"dayChange"`
	WeekChange  *float64 `json:"weekChange"`
	YearChange  *float64 `json:"yearChange"`
}

var DefaultTickers = []string{
	"SPY",
	"QQQ",
	"VTI",
	"VT",
	"SCHD",
	"REIT",
	"IAU",
}
