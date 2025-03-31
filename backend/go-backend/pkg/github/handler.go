package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go-backend/pkg/database"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

type Handler struct {
	client *http.Client
}

func NewHandler() *Handler {
	return &Handler{
		client: &http.Client{
			Timeout: 20 * time.Second,
		},
	}
}

// Helper to extract integer from string like "1,234 stars" or "553"
func extractIntFromString(text string) int {
	re := regexp.MustCompile(`[0-9,]+`)
	numStr := re.FindString(text)
	numStr = strings.ReplaceAll(numStr, ",", "")
	num, err := strconv.Atoi(numStr)
	if err != nil {
		// Log the error but return 0 if conversion fails
		log.Printf("[GitHub Helper] Failed to convert string '%s' (extracted from '%s') to int: %v", numStr, text, err)
		return 0
	}
	return num
}

// Helper to extract color from style attribute like "background-color: #3178c6"
func extractColorFromStyle(style string) string {
	re := regexp.MustCompile(`background-color:\s*(#[0-9a-fA-F]{6})`)
	matches := re.FindStringSubmatch(style)
	if len(matches) > 1 {
		return matches[1]
	}
	return "" // Return empty if no match
}

// FetchTrendingRepos fetches trending repositories from GitHub trending page and parses the HTML
func (h *Handler) FetchTrendingRepos() ([]Repository, error) {
	log.Printf("[GitHub] Fetching trending repositories from github.com/trending")
	// Using ?since=daily explicitly, although it might be the default
	trendingURL := "https://github.com/trending?since=daily"

	req, err := http.NewRequest("GET", trendingURL, nil)
	if err != nil {
		log.Printf("[GitHub] Failed to create request: %v", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Mimic a browser User-Agent, GitHub might block default Go clients
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")

	resp, err := h.client.Do(req)
	if err != nil {
		log.Printf("[GitHub] Request to %s failed: %v", trendingURL, err)
		return nil, fmt.Errorf("request to %s failed: %w", trendingURL, err)
	}
	defer resp.Body.Close()

	log.Printf("[GitHub] Response status from %s: %s", trendingURL, resp.Status)

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body) // Read body for error context
		log.Printf("[GitHub] Request to %s returned non-OK status %d: %s", trendingURL, resp.StatusCode, string(bodyBytes))
		return nil, fmt.Errorf("request to %s returned status %s", trendingURL, resp.Status)
	}

	// Parse the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("[GitHub] Failed to parse HTML from %s: %v", trendingURL, err)
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var repos []Repository
	baseURL, _ := url.Parse("https://github.com")

	// Find each repository box
	doc.Find("article.Box-row").Each(func(i int, s *goquery.Selection) {
		var repo Repository

		// --- Extract Repo Name, Author, and URL ---
		repoLink := s.Find("h2 a.Link") // More specific selector for the repo link
		repoFullName := strings.TrimSpace(repoLink.Text())
		repoPath, exists := repoLink.Attr("href")
		if !exists {
			log.Printf("[GitHub Parser] Could not find href for repo %d", i)
			return // Skip this repo if essential info is missing
		}
		repo.URL = baseURL.ResolveReference(&url.URL{Path: strings.TrimSpace(repoPath)}).String()

		// Extract author and name from full name like "author / name"
		parts := strings.Split(strings.ReplaceAll(repoFullName, " ", ""), "/")
		if len(parts) == 2 {
			repo.Author = parts[0]
			repo.Name = parts[1]
		} else {
			log.Printf("[GitHub Parser] Could not parse author/name from '%s' for repo %d", repoFullName, i)
			// Attempt fallback using the href path
			pathParts := strings.Split(strings.Trim(repoPath, "/"), "/")
			if len(pathParts) >= 2 {
				repo.Author = pathParts[0]
				repo.Name = pathParts[1]
			} else {
				return // Skip if name/author cannot be determined
			}
		}

		// --- Extract Owner Avatar URL ---
		// Construct the owner avatar URL based on convention
		repo.Avatar = fmt.Sprintf("https://github.com/%s.png?size=40", repo.Author)

		// --- Extract Description ---
		// Target the <p> tag directly following the h2, careful about structure changes
		repo.Description = strings.TrimSpace(s.Find("p.col-9").First().Text())

		// --- Extract Language and Color ---
		langSpan := s.Find("span[itemprop='programmingLanguage']")
		repo.Language = strings.TrimSpace(langSpan.Text())
		if repo.Language != "" {
			colorSpan := langSpan.Parent().Find(".repo-language-color")
			style, _ := colorSpan.Attr("style")
			repo.LanguageColor = extractColorFromStyle(style)
		}

		// --- Extract Stars and Forks ---
		starLink := s.Find("a[href$='/stargazers']") // Find link ending with /stargazers
		forkLink := s.Find("a[href$='/forks']")      // Find link ending with /forks

		repo.Stars = extractIntFromString(starLink.Text())
		repo.Forks = extractIntFromString(forkLink.Text())

		// --- Extract Current Period Stars ---
		// This is often in a span like: <span class="d-inline-block float-sm-right"> ... X stars today </span>
		starsTodaySpan := s.Find("span.float-sm-right") // Try the specific floated span first
		if starsTodaySpan.Length() == 0 {
			// Fallback: find any span containing "stars today" within the .f6 div
			s.Find(".f6 span").EachWithBreak(func(_ int, spanNode *goquery.Selection) bool {
				if strings.Contains(spanNode.Text(), "stars today") {
					starsTodaySpan = spanNode
					return false // Stop searching
				}
				return true
			})
		}
		repo.CurrentPeriodStars = extractIntFromString(starsTodaySpan.Text())

		// --- Extract Built By ---
		s.Find("span:contains('Built by') a").Each(func(_ int, contributorLink *goquery.Selection) {
			href, hrefExists := contributorLink.Attr("href")
			img := contributorLink.Find("img")
			avatar, avatarExists := img.Attr("src")
			username := ""
			if hrefExists {
				href = strings.TrimPrefix(href, "/") // Remove leading slash
				parts := strings.Split(href, "/")
				if len(parts) > 0 {
					username = parts[0] // Assuming first part is username
				}
			}
			// Fallback to alt text if username couldn't be parsed from href
			if username == "" {
				alt, altExists := img.Attr("alt")
				if altExists {
					username = strings.TrimPrefix(alt, "@")
				}
			}

			if hrefExists && avatarExists && username != "" {
				contributor := Contributor{
					Username: username,
					Href:     baseURL.ResolveReference(&url.URL{Path: href}).String(),
					Avatar:   avatar,
				}
				repo.BuiltBy = append(repo.BuiltBy, contributor)
			} else {
				log.Printf("[GitHub Parser] Could not extract full contributor info for repo %s/%s (href:%v, avatar:%v, user:%s)", repo.Author, repo.Name, hrefExists, avatarExists, username)
			}
		})

		repos = append(repos, repo)
	})

	log.Printf("[GitHub] Successfully parsed %d repositories from HTML", len(repos))

	// Store repos in database (Keep this part unchanged)
	db := database.GetDB()
	stored := 0
	failedToStore := 0
	for _, repo := range repos {
		builtByJSON, err := json.Marshal(repo.BuiltBy)
		if err != nil {
			log.Printf("[GitHub DB] Failed to marshal builtBy for repo %s/%s: %v", repo.Author, repo.Name, err)
			failedToStore++
			continue // Skip storing this repo
		}

		// Ensure language color is not null if language is empty
		langColor := repo.LanguageColor
		if repo.Language == "" {
			langColor = "" // Or potentially a default value if your DB requires non-null
		}

		_, err = db.Exec(`
			INSERT OR REPLACE INTO github_repositories
			(author, name, avatar, url, description, language, language_color, stars, forks, current_period_stars, built_by)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			repo.Author,             // author
			repo.Name,               // name
			repo.Avatar,             // avatar (owner)
			repo.URL,                // url
			repo.Description,        // description
			repo.Language,           // language
			langColor,               // language_color
			repo.Stars,              // stars
			repo.Forks,              // forks
			repo.CurrentPeriodStars, // current_period_stars
			builtByJSON,             // built_by
		)
		if err != nil {
			log.Printf("[GitHub DB] Failed to store repo %s/%s in database: %v", repo.Author, repo.Name, err)
			failedToStore++
			continue // Skip to next repo
		}
		stored++
	}
	if failedToStore > 0 {
		log.Printf("[GitHub DB] Stored %d repositories, failed to store %d", stored, failedToStore)
	} else {
		log.Printf("[GitHub DB] Successfully stored %d/%d repositories in database", stored, len(repos))
	}

	return repos, nil
}

// GetTrendingRepos tries to get recent data from DB, otherwise fetches fresh data.
func (h *Handler) GetTrendingRepos(c *fiber.Ctx) error {
	// Try to get from database first, checking freshness directly in SQL
	db := database.GetDB()
	cacheDuration := "-15 minutes" // Define cache duration (adjust as needed)
	log.Printf("[GitHub] Checking database for repositories updated within the last %s", strings.Replace(cacheDuration, "-", "", 1))

	rows, err := db.Query(`
		SELECT author, name, avatar, url, description, language, language_color,
		       stars, forks, current_period_stars, built_by
		FROM github_repositories
		WHERE created_at >= datetime('now', ?)
		ORDER BY current_period_stars DESC, stars DESC
		LIMIT 25
	`, cacheDuration)

	var reposFromDB []Repository
	var queryError error = err // Store potential query error

	if queryError == nil {
		defer rows.Close()
		for rows.Next() {
			var repo Repository
			var builtByJSON []byte
			// Note: We don't need to scan created_at anymore

			scanErr := rows.Scan(
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
			if scanErr != nil {
				log.Printf("[GitHub] Failed to scan repository from database: %v", scanErr)
				continue // Skip this row
			}

			if err := json.Unmarshal(builtByJSON, &repo.BuiltBy); err != nil {
				log.Printf("[GitHub] Failed to unmarshal builtBy JSON for repo %s/%s: %v", repo.Author, repo.Name, err)
				// Continue with the repo even if builtBy fails to unmarshal
			}

			reposFromDB = append(reposFromDB, repo)
		}
		// Check for errors that might have occurred during row iteration
		if err = rows.Err(); err != nil {
			queryError = fmt.Errorf("error iterating database rows: %w", err)
			log.Printf("[GitHub] %v", queryError)
		}
	} else {
		log.Printf("[GitHub] Initial DB query failed: %v", queryError)
	}

	// --- Decision Logic ---
	// If the query succeeded AND returned rows, we have a cache hit.
	if queryError == nil && len(reposFromDB) > 0 {
		log.Printf("[GitHub] Cache hit: Returned %d repositories from database", len(reposFromDB))
		return c.JSON(reposFromDB)
	}

	// Cache miss (query failed OR query succeeded but returned 0 rows because data was too old)
	if queryError != nil {
		log.Printf("[GitHub] Cache miss due to DB query error: %v. Fetching fresh data.", queryError)
	} else {
		log.Printf("[GitHub] Cache miss: No recent data found in DB. Fetching fresh data.")
	}

	// Fetch fresh data using the scraper
	repos, fetchErr := h.FetchTrendingRepos()
	if fetchErr != nil {
		log.Printf("[GitHub] Failed to fetch repositories after cache miss: %v", fetchErr)
		// Important: If fetch fails AFTER a cache miss, we MUST return an error.
		// We cannot return reposFromDB because it's either empty or contains stale data.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to fetch trending repositories: %v", fetchErr),
		})
	}

	// Return freshly fetched data
	return c.JSON(repos)
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	cacheConfig := cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		Expiration:   60 * time.Minute,
		CacheControl: true,
	}

	app.Get("/github/trending", cache.New(cacheConfig), h.GetTrendingRepos)
	log.Printf("[GitHub] Routes registered with %v cache expiration", cacheConfig.Expiration)
}
