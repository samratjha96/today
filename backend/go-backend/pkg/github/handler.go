package github

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

// FetchTrendingRepos fetches trending repositories from GitHub API and stores them in the database
func (h *Handler) FetchTrendingRepos() ([]Repository, error) {
	log.Printf("[GitHub] Fetching trending repositories from API")
	req, err := http.NewRequest("GET", "https://api.gitterapp.com/repositories?since=daily", nil)
	if err != nil {
		log.Printf("[GitHub] Failed to create request: %v", err)
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		log.Printf("[GitHub] API request failed: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Log the response status code
	log.Printf("[GitHub] API response status: %s", resp.Status)

	// Read the response body into a buffer so we can inspect it
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[GitHub] Failed to read response body: %v", err)
		return nil, err
	}

	// Check if response is not a success status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("[GitHub] API returned error status %d: %s", resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("api returned status %d: %s", resp.StatusCode, resp.Status)
	}

	// Preview the response (only show first 100 characters if very long)
	preview := string(bodyBytes)
	if len(preview) > 100 {
		preview = preview[:100] + "..."
	}
	log.Printf("[GitHub] Response preview: %s", preview)

	var repos []Repository
	if err := json.Unmarshal(bodyBytes, &repos); err != nil {
		log.Printf("[GitHub] Failed to decode API response: %v", err)
		log.Printf("[GitHub] Raw response: %s", string(bodyBytes))
		return nil, err
	}

	log.Printf("[GitHub] Successfully fetched %d repositories from API", len(repos))

	// Store repos in database
	db := database.GetDB()
	stored := 0
	for _, repo := range repos {
		builtByJSON, err := json.Marshal(repo.BuiltBy)
		if err != nil {
			log.Printf("[GitHub] Failed to marshal builtBy for repo %s/%s: %v", repo.Author, repo.Name, err)
			continue
		}

		_, err = db.Exec(`
			INSERT OR REPLACE INTO github_repositories 
			(author, name, avatar, url, description, language, language_color, stars, forks, current_period_stars, built_by)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			repo.Author,
			repo.Name,
			repo.Avatar,
			repo.URL,
			repo.Description,
			repo.Language,
			repo.LanguageColor,
			repo.Stars,
			repo.Forks,
			repo.CurrentPeriodStars,
			builtByJSON,
		)
		if err != nil {
			log.Printf("[GitHub] Failed to store repo %s/%s in database: %v", repo.Author, repo.Name, err)
			continue
		}
		stored++
	}
	log.Printf("[GitHub] Successfully stored %d/%d repositories in database", stored, len(repos))

	return repos, nil
}

func (h *Handler) GetTrendingRepos(c *fiber.Ctx) error {
	// Try to get from database first
	db := database.GetDB()
	rows, err := db.Query(`
		SELECT author, name, avatar, url, description, language, language_color, 
		       stars, forks, current_period_stars, built_by
		FROM github_repositories
		WHERE created_at >= datetime('now', '-5 minutes')
		ORDER BY stars DESC
	`)
	if err == nil {
		defer rows.Close()

		var repos []Repository
		for rows.Next() {
			var repo Repository
			var builtByJSON []byte
			err := rows.Scan(
				&repo.Author,
				&repo.Name,
				&repo.Avatar,
				&repo.URL,
				&repo.Description,
				&repo.Language,
				&repo.LanguageColor,
				&repo.Stars,
				&repo.Forks,
				&repo.CurrentPeriodStars,
				&builtByJSON,
			)
			if err != nil {
				log.Printf("[GitHub] Failed to scan repository from database: %v", err)
				continue
			}

			if err := json.Unmarshal(builtByJSON, &repo.BuiltBy); err != nil {
				log.Printf("[GitHub] Failed to unmarshal builtBy JSON for repo %s/%s: %v", repo.Author, repo.Name, err)
				continue
			}

			repos = append(repos, repo)
		}

		if len(repos) > 0 {
			log.Printf("[GitHub] Cache hit: Returned %d repositories from database", len(repos))
			return c.JSON(repos)
		}
	}

	log.Printf("[GitHub] Cache miss: Fetching repositories from API")
	repos, err := h.FetchTrendingRepos()
	if err != nil {
		log.Printf("[GitHub] Failed to fetch repositories: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(repos)
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	cacheConfig := cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		Expiration:   5 * time.Minute,
		CacheControl: true,
	}

	app.Get("/github/trending", cache.New(cacheConfig), h.GetTrendingRepos)
	log.Printf("[GitHub] Routes registered with %v cache expiration", cacheConfig.Expiration)
}
