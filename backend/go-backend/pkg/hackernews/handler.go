package hackernews

import (
	"encoding/json"
	"fmt"
	"io"
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

// Handler handles HackerNews related requests
type Handler struct {
	client *http.Client
}

// NewHandler creates a new HackerNews handler
func NewHandler() *Handler {
	return &Handler{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// RegisterRoutes registers the HackerNews routes with caching
func (h *Handler) RegisterRoutes(app *fiber.App) {
	// Cache middleware configuration
	cacheConfig := cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true" // Skip cache if refresh query param is true
		},
		Expiration:   5 * time.Minute,
		CacheControl: true,
	}

	// Apply cache middleware only to the top stories endpoint
	app.Get("/hackernews/top", cache.New(cacheConfig), h.GetTopStories)
}

// GetTopStories returns the top HackerNews stories
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
				continue
			}
			stories = append(stories, story)
		}

		if len(stories) > 0 {
			return c.JSON(stories)
		}
	}

	// If no recent data in database or error occurred, fetch from API
	// Get top story IDs
	resp, err := h.client.Get(hackerNewsTopStoriesURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch top stories: %v", err),
		})
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to read response body: %v", err),
		})
	}

	var storyIDs []int
	if err := json.Unmarshal(body, &storyIDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to parse story IDs: %v", err),
		})
	}

	// Get details for top 10 stories
	stories := make([]Story, 0, 10)
	for _, id := range storyIDs[:10] {
		storyURL := fmt.Sprintf(hackerNewsStoryURL, id)
		resp, err := h.client.Get(storyURL)
		if err != nil {
			continue // Skip this story if there's an error
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		var story Story
		if err := json.Unmarshal(body, &story); err != nil {
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
			continue
		}

		stories = append(stories, story)
	}

	if len(stories) == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch any stories",
		})
	}

	return c.JSON(stories)
}
