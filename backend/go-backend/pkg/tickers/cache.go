package tickers

import (
	"sync"
	"time"
)

type cachedData struct {
	data      []TickerData
	timestamp time.Time
}

var (
	cache      *cachedData
	cacheMutex sync.RWMutex
	cacheTime  = 2 * time.Minute
)

func getCachedData() ([]TickerData, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()

	if cache == nil {
		return nil, false
	}

	if time.Since(cache.timestamp) > cacheTime {
		return nil, false
	}

	return cache.data, true
}

func updateCache(data []TickerData) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	cache = &cachedData{
		data:      data,
		timestamp: time.Now(),
	}
}
