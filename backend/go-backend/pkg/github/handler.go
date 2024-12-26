package github

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
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

func (h *Handler) fetchTrendingRepos() ([]Repository, error) {
	req, err := http.NewRequest("GET", "https://api.gitterapp.com/repositories?since=daily", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repos []Repository
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	return repos, nil
}

func (h *Handler) GetTrendingRepos(c *fiber.Ctx) error {
	repos, err := h.fetchTrendingRepos()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(repos)
}

// RegisterRoutes registers the GitHub routes with the given fiber app
func (h *Handler) RegisterRoutes(app *fiber.App) {
	// Cache middleware configuration
	cacheConfig := cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true" // Skip cache if refresh query param is true
		},
		Expiration:   5 * time.Minute,
		CacheControl: true,
	}

	// Apply cache middleware only to the trending endpoint
	app.Get("/github/trending", cache.New(cacheConfig), h.GetTrendingRepos)
}
