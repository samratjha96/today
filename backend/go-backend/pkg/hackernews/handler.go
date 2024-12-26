package hackernews

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// Handler handles HackerNews related requests
type Handler struct{}

// NewHandler creates a new HackerNews handler
func NewHandler() *Handler {
	return &Handler{}
}

// GetTopStories returns the top HackerNews stories
func (h *Handler) GetTopStories(c *fiber.Ctx) error {
	// Mock data for now
	currentTime := time.Now().Unix()
	stories := []Story{
		{
			By:          "nateb2022",
			Descendants: 26,
			ID:          42512896,
			Score:       100,
			Time:        currentTime,
			Title:       "Blackcandy: Self hosted music streaming server",
			Type:        "story",
			URL:         "https://github.com/blackcandy-org/blackcandy",
		},
		{
			By:          "dhouston",
			Descendants: 71,
			ID:          42512897,
			Score:       150,
			Time:        currentTime - 3600, // 1 hour ago
			Title:       "Rust vs Go: A Systems Programming Showdown",
			Type:        "story",
			URL:         "https://example.com/rust-vs-go",
		},
		{
			By:          "pg",
			Descendants: 42,
			ID:          42512898,
			Score:       200,
			Time:        currentTime - 7200, // 2 hours ago
			Title:       "The Future of Web Development in 2024",
			Type:        "story",
			URL:         "https://example.com/web-dev-2024",
		},
		{
			By:          "sama",
			Descendants: 89,
			ID:          42512899,
			Score:       300,
			Time:        currentTime - 10800, // 3 hours ago
			Title:       "New Developments in Large Language Models",
			Type:        "story",
			URL:         "https://example.com/llm-developments",
		},
		{
			By:          "patio11",
			Descendants: 55,
			ID:          42512900,
			Score:       250,
			Time:        currentTime - 14400, // 4 hours ago
			Title:       "The Economics of Software Development",
			Type:        "story",
			URL:         "https://example.com/software-economics",
		},
		{
			By:          "tptacek",
			Descendants: 34,
			ID:          42512901,
			Score:       180,
			Time:        currentTime - 18000, // 5 hours ago
			Title:       "Advanced Cryptography Patterns in Modern Applications",
			Type:        "story",
			URL:         "https://example.com/crypto-patterns",
		},
		{
			By:          "janesmith",
			Descendants: 63,
			ID:          42512902,
			Score:       220,
			Time:        currentTime - 21600, // 6 hours ago
			Title:       "Understanding WebAssembly Performance",
			Type:        "story",
			URL:         "https://example.com/wasm-performance",
		},
	}

	return c.JSON(stories)
}
