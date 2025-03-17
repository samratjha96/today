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
	cacheTime  = 5 * time.Minute
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

// ExtendCacheTime can be called to increase cache time when rate limiting is happening
func ExtendCacheTime(additional time.Duration) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if cache != nil {
		// Update the timestamp to extend the cache lifetime
		cache.timestamp = time.Now()
	}
}
