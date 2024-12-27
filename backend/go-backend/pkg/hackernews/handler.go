package hackernews

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"go-backend/pkg/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

const (
	hackerNewsTopStoriesURL = "https://hacker-news.firebaseio.com/v0/topstories.json"
	hackerNewsStoryURL      = "https://hacker-news.firebaseio.com/v0/item/%d.json"
)

type Handler struct {
	client *http.Client
}

func NewHandler() *Handler {
	return &Handler{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchTopStories fetches top stories from HackerNews API and stores them in the database
func (h *Handler) FetchTopStories() ([]Story, error) {
	// Get top story IDs
	resp, err := h.client.Get(hackerNewsTopStoriesURL)
	if err != nil {
		log.Printf("[HackerNews] Failed to fetch top story IDs: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[HackerNews] Failed to read top stories response: %v", err)
		return nil, err
	}

	var storyIDs []int
	if err := json.Unmarshal(body, &storyIDs); err != nil {
		log.Printf("[HackerNews] Failed to parse story IDs: %v", err)
		return nil, err
	}

	log.Printf("[HackerNews] Successfully fetched %d story IDs, processing top 10", len(storyIDs))

	db := database.GetDB()
	stories := make([]Story, 0, 10)
	stored := 0

	// Get details for top 10 stories
	for _, id := range storyIDs[:10] {
		storyURL := fmt.Sprintf(hackerNewsStoryURL, id)
		resp, err := h.client.Get(storyURL)
		if err != nil {
			log.Printf("[HackerNews] Failed to fetch story %d: %v", id, err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			log.Printf("[HackerNews] Failed to read story %d response: %v", id, err)
			continue
		}

		var story Story
		if err := json.Unmarshal(body, &story); err != nil {
			log.Printf("[HackerNews] Failed to parse story %d: %v", id, err)
			continue
		}

		// Store story in database
		_, err = db.Exec(`
			INSERT OR REPLACE INTO hackernews_stories 
			(id, by, descendants, score, time, title, type, url)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		`,
			story.ID,
			story.By,
			story.Descendants,
			story.Score,
			story.Time,
			story.Title,
			story.Type,
			story.URL,
		)
		if err != nil {
			log.Printf("[HackerNews] Failed to store story %d in database: %v", id, err)
			continue
		}
		stored++

		stories = append(stories, story)
	}

	if len(stories) == 0 {
		return nil, fmt.Errorf("failed to fetch any valid stories")
	}

	log.Printf("[HackerNews] Successfully processed %d stories, stored %d in database", len(stories), stored)
	return stories, nil
}

func (h *Handler) GetTopStories(c *fiber.Ctx) error {
	// Try to get stories from database first
	db := database.GetDB()
	rows, err := db.Query(`
		SELECT id, by, descendants, score, time, title, type, url
		FROM hackernews_stories
		WHERE created_at >= datetime('now', '-5 minutes')
		ORDER BY score DESC
		LIMIT 10
	`)
	if err == nil {
		defer rows.Close()

		var stories []Story
		for rows.Next() {
			var story Story
			err := rows.Scan(
				&story.ID,
				&story.By,
				&story.Descendants,
				&story.Score,
				&story.Time,
				&story.Title,
				&story.Type,
				&story.URL,
			)
			if err != nil {
				log.Printf("[HackerNews] Failed to scan story from database: %v", err)
				continue
			}
			stories = append(stories, story)
		}

		if len(stories) > 0 {
			log.Printf("[HackerNews] Cache hit: Returned %d stories from database", len(stories))
			return c.JSON(stories)
		}
	}

	log.Printf("[HackerNews] Cache miss: Fetching stories from API")
	stories, err := h.FetchTopStories()
	if err != nil {
		log.Printf("[HackerNews] Failed to fetch stories: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch stories: %v", err),
		})
	}

	return c.JSON(stories)
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	cacheConfig := cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		Expiration:   5 * time.Minute,
		CacheControl: true,
	}

	app.Get("/hackernews/top", cache.New(cacheConfig), h.GetTopStories)
	log.Printf("[HackerNews] Routes registered with %v cache expiration", cacheConfig.Expiration)
}
