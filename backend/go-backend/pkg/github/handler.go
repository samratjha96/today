package github

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	cache "github.com/patrickmn/go-cache"
)

var (
	// Create a cache with a default expiration time of 5 minutes and purge unused items every 10 minutes
	repoCache = cache.New(5*time.Minute, 10*time.Minute)
	cacheKey  = "trending_repos"
)

func fetchTrendingRepos() ([]Repository, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", "https://api.gitterapp.com/repositories?since=daily", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
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

// RegisterRoutes registers the GitHub routes with the given fiber app
func RegisterRoutes(router fiber.Router) {
	github := router.Group("/github")

	github.Get("/trending", func(c *fiber.Ctx) error {
		// Try to get data from cache first
		if cached, found := repoCache.Get(cacheKey); found {
			return c.JSON(cached.([]Repository))
		}

		// Cache miss - fetch new data
		repos, err := fetchTrendingRepos()
		if err != nil {
			// If fetch fails, try to get stale data from cache with no expiration
			if cached, found := repoCache.Get(cacheKey); found {
				return c.JSON(cached.([]Repository))
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Store in cache
		repoCache.Set(cacheKey, repos, cache.DefaultExpiration)
		return c.JSON(repos)
	})
}
