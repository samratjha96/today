package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"go-backend/pkg/database"
	"go-backend/pkg/github"
	"go-backend/pkg/hackernews"
	"go-backend/pkg/rss"
	"go-backend/pkg/tickers"
)

type Job struct {
	name     string
	interval time.Duration
	handler  func() error
}

type JobScheduler struct {
	jobs []*Job
	wg   sync.WaitGroup
}

func NewJobScheduler() *JobScheduler {
	return &JobScheduler{
		jobs: make([]*Job, 0),
	}
}

func (s *JobScheduler) AddJob(name string, interval time.Duration, handler func() error) {
	s.jobs = append(s.jobs, &Job{
		name:     name,
		interval: interval,
		handler:  handler,
	})
}

func (s *JobScheduler) Start() {
	for _, job := range s.jobs {
		s.wg.Add(1)
		go s.runJob(job)
	}
}

func (s *JobScheduler) Stop() {
	s.wg.Wait()
}

func (s *JobScheduler) runJob(job *Job) {
	defer s.wg.Done()

	ticker := time.NewTicker(job.interval)
	defer ticker.Stop()

	log.Printf("Running job: %s", job.name)
	if err := job.handler(); err != nil {
		log.Printf("Error running job %s: %v", job.name, err)
	}

	for range ticker.C {
		log.Printf("Running job: %s", job.name)
		if err := job.handler(); err != nil {
			log.Printf("Error running job %s: %v", job.name, err)
		}
	}
}

func scheduledJobs(ghHandler *github.Handler, hnHandler *hackernews.Handler, rssHandler *rss.Handler) {
	scheduler := NewJobScheduler()

	scheduler.AddJob("GitHub Trending", time.Hour, func() error {
		_, err := ghHandler.FetchTrendingRepos()
		return err
	})

	scheduler.AddJob("HackerNews Top", 15*time.Minute, func() error {
		_, err := hnHandler.FetchTopStories()
		return err
	})

	// Add RSS feed job
	rssHandler.AddToJobScheduler(scheduler.AddJob)

	scheduler.Start()
}

func main() {
	if err := database.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	app := fiber.New()

	app.Use(logger.New())
	app.Use(recover.New())

	// Get allowed origins from environment variable or use default
	allowedHosts := os.Getenv("ALLOWED_HOSTS")
	if allowedHosts == "" {
		allowedHosts = "*"
	}

	// Configure CORS with proper prefixes if needed
	corsConfig := cors.Config{
		AllowHeaders: "Origin, Content-Type, Accept",
	}

	// If allowedHosts is "*", use that directly
	if allowedHosts == "*" {
		corsConfig.AllowOrigins = "*"
	} else {
		// For specific domains, ensure they are proper URLs
		corsConfig.AllowOrigins = "https://" + allowedHosts + ",http://" + allowedHosts
	}

	app.Use(cors.New(corsConfig))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	ghHandler := github.NewHandler()
	ghHandler.RegisterRoutes(app)

	hnHandler := hackernews.NewHandler()
	hnHandler.RegisterRoutes(app)

	tickerHandler := tickers.NewHandler()
	tickerHandler.RegisterRoutes(app)

	// Initialize RSS handler and register routes
	rssHandler := rss.NewHandler()
	rssHandler.RegisterRoutes(app)

	go scheduledJobs(ghHandler, hnHandler, rssHandler)

	log.Printf("Server starting on port %s\n", port)
	log.Fatal(app.Listen(":" + port))
}
