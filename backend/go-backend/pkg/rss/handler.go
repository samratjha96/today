package rss

import (
	"encoding/xml"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"go-backend/pkg/database"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Handler for RSS feed operations
type Handler struct {
	client *http.Client
}

// XML structures for RSS parsing
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string     `xml:"title"`
	Description string     `xml:"description"`
	Link        string     `xml:"link"`
	Items       []RSSEntry `xml:"item"`
}

type RSSEntry struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// NewHandler creates a new RSS handler
func NewHandler() *Handler {
	return &Handler{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// Initialize creates the RSS news table if it doesn't exist
func (h *Handler) Initialize() error {
	db := database.GetDB()

	// Create RSS news items table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS rss_news (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			source TEXT NOT NULL,
			title TEXT NOT NULL,
			link TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

// FetchRSSFeed fetches and parses an RSS feed from a given URL
func (h *Handler) FetchRSSFeed(url string) ([]RSSEntry, error) {
	log.Printf("[RSS] Fetching feed from %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[RSS] Failed to create request: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set user agent to avoid being blocked
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/rss+xml, application/xml, text/xml, */*")

	resp, err := h.client.Do(req)
	if err != nil {
		log.Printf("[RSS] Request to %s failed: %v", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[RSS] Request to %s returned non-OK status %d", url, resp.StatusCode)
		return nil, fmt.Errorf("request to %s returned status %s", url, resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[RSS] Failed to read response body: %v", err)
		return nil, err
	}

	var rss RSS
	err = xml.Unmarshal(bodyBytes, &rss)
	if err != nil {
		// Some feeds might have different formats, try alternative parsing here if needed
		log.Printf("[RSS] Failed to parse XML from %s: %v", url, err)
		return nil, err
	}

	return rss.Channel.Items, nil
}

// StoreRSSItems stores RSS items in the database
func (h *Handler) StoreRSSItems(source string, items []RSSEntry) (int, error) {
	db := database.GetDB()

	stored := 0
	for _, item := range items {
		if strings.TrimSpace(item.Link) == "" {
			continue // Skip items without links
		}

		_, err := db.Exec(`
			INSERT OR IGNORE INTO rss_news
			(source, title, link)
			VALUES (?, ?, ?)
		`,
			source,
			strings.TrimSpace(item.Title),
			strings.TrimSpace(item.Link),
		)
		if err != nil {
			log.Printf("[RSS DB] Failed to store RSS item '%s' in database: %v", item.Title, err)
			continue
		}
		stored++
	}

	log.Printf("[RSS DB] Successfully stored %d/%d items from %s in database", stored, len(items), source)
	return stored, nil
}

// FetchAllFeeds fetches all RSS feeds and updates the database
func (h *Handler) FetchAllFeeds() ([]RSSItem, error) {
	var allNews []RSSItem

	for source, url := range RSSFeedList {
		entries, err := h.FetchRSSFeed(url)
		if err != nil {
			log.Printf("[RSS] Error fetching feed from %s: %v", source, err)
			continue
		}

		// Limit to top 5 entries
		limit := 5
		if len(entries) < limit {
			limit = len(entries)
		}

		// Store in database
		_, err = h.StoreRSSItems(source, entries[:limit])
		if err != nil {
			log.Printf("[RSS] Error storing RSS items from %s: %v", source, err)
		}

		// Convert to response format
		for i := 0; i < limit; i++ {
			if strings.TrimSpace(entries[i].Link) != "" {
				allNews = append(allNews, RSSItem{
					Source: source,
					Title:  strings.TrimSpace(entries[i].Title),
					Link:   strings.TrimSpace(entries[i].Link),
				})
			}
		}
	}

	return allNews, nil
}

// GetNewsFromDB retrieves recent news items from the database
func (h *Handler) GetNewsFromDB() ([]RSSItem, error) {
	db := database.GetDB()

	rows, err := db.Query(`
		SELECT source, title, link
		FROM rss_news
		WHERE created_at >= datetime('now', '-1 hour')
		ORDER BY created_at DESC
		LIMIT 25
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var news []RSSItem
	for rows.Next() {
		var item RSSItem
		if err := rows.Scan(&item.Source, &item.Title, &item.Link); err != nil {
			log.Printf("[RSS] Failed to scan news item from database: %v", err)
			continue
		}
		news = append(news, item)
	}

	return news, nil
}

// GetNews retrieves news items, either from DB cache or freshly fetched
func (h *Handler) GetNews(c *fiber.Ctx) error {
	// Try to get from database first
	news, err := h.GetNewsFromDB()
	if err == nil && len(news) > 0 {
		log.Printf("[RSS] Cache hit: Returned %d news items from database", len(news))
		return c.JSON(news)
	}

	// If database query failed or returned no results, fetch fresh data
	if err != nil {
		log.Printf("[RSS] Cache miss due to DB query error: %v. Fetching fresh data.", err)
	} else {
		log.Printf("[RSS] Cache miss: No recent data found in DB. Fetching fresh data.")
	}

	// Fetch fresh data from RSS feeds
	freshNews, err := h.FetchAllFeeds()
	if err != nil {
		log.Printf("[RSS] Failed to fetch news after cache miss: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch news: %v", err),
		})
	}

	return c.JSON(freshNews)
}

// RegisterRoutes registers the RSS endpoints with the Fiber app
func (h *Handler) RegisterRoutes(app *fiber.App) {
	// Initialize database tables
	if err := h.Initialize(); err != nil {
		log.Fatalf("[RSS] Failed to initialize database: %v", err)
	}

	cacheConfig := cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		Expiration:   15 * time.Minute,
		CacheControl: true,
	}

	app.Get("/news", cache.New(cacheConfig), h.GetNews)
	log.Printf("[RSS] Routes registered with %v cache expiration", cacheConfig.Expiration)
}

// AddToJobScheduler adds periodic RSS feed fetching to the scheduler
func (h *Handler) AddToJobScheduler(addJob func(string, time.Duration, func() error)) {
	addJob("RSS Feeds", 30*time.Minute, func() error {
		_, err := h.FetchAllFeeds()
		return err
	})
}
