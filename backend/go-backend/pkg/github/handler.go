package github

import (
	"encoding/json"
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

	// Store repos in database
	db := database.GetDB()
	for _, repo := range repos {
		builtByJSON, err := json.Marshal(repo.BuiltBy)
		if err != nil {
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
			continue
		}
	}

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
				continue
			}

			// Parse built_by JSON
			if err := json.Unmarshal(builtByJSON, &repo.BuiltBy); err != nil {
				continue
			}

			repos = append(repos, repo)
		}

		if len(repos) > 0 {
			return c.JSON(repos)
		}
	}

	// If no recent data in database or error occurred, fetch from API
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
