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

	// Get top story IDs
	resp, err := h.client.Get(hackerNewsTopStoriesURL)
	if err != nil {
		log.Printf("[HackerNews] Failed to fetch top story IDs: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch top stories: %v", err),
		})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[HackerNews] Failed to read top stories response: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to read response body: %v", err),
		})
	}

	var storyIDs []int
	if err := json.Unmarshal(body, &storyIDs); err != nil {
		log.Printf("[HackerNews] Failed to parse story IDs: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to parse story IDs: %v", err),
		})
	}

	log.Printf("[HackerNews] Successfully fetched %d story IDs, processing top 10", len(storyIDs))

	// Get details for top 10 stories
	stories := make([]Story, 0, 10)
	stored := 0
	for _, id := range storyIDs[:10] {
		storyURL := fmt.Sprintf(hackerNewsStoryURL, id)
		resp, err := h.client.Get(storyURL)
		if err != nil {
			log.Printf("[HackerNews] Failed to fetch story %d: %v", id, err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
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
		log.Printf("[HackerNews] Failed to fetch any valid stories")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch any stories",
		})
	}

	log.Printf("[HackerNews] Successfully processed %d stories, stored %d in database", len(stories), stored)
	return c.JSON(stories)
}
