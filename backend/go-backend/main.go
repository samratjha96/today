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
	"go-backend/pkg/tickers"
)

// Job represents a scheduled job with its configuration
type Job struct {
	name     string
	interval time.Duration
	handler  func() error
}

// JobScheduler manages multiple jobs with different intervals
type JobScheduler struct {
	jobs []*Job
	wg   sync.WaitGroup
}

// NewJobScheduler creates a new job scheduler
func NewJobScheduler() *JobScheduler {
	return &JobScheduler{
		jobs: make([]*Job, 0),
	}
}

// AddJob adds a new job to the scheduler
func (s *JobScheduler) AddJob(name string, interval time.Duration, handler func() error) {
	s.jobs = append(s.jobs, &Job{
		name:     name,
		interval: interval,
		handler:  handler,
	})
}

// Start starts all jobs in separate goroutines
func (s *JobScheduler) Start() {
	for _, job := range s.jobs {
		s.wg.Add(1)
		go s.runJob(job)
	}
}

// Stop waits for all jobs to complete
func (s *JobScheduler) Stop() {
	s.wg.Wait()
}

// runJob executes a single job on its specified interval
func (s *JobScheduler) runJob(job *Job) {
	defer s.wg.Done()

	ticker := time.NewTicker(job.interval)
	defer ticker.Stop()

	// Run immediately on startup
	log.Printf("Running job: %s", job.name)
	if err := job.handler(); err != nil {
		log.Printf("Error running job %s: %v", job.name, err)
	}

	// Then run on the specified interval
	for range ticker.C {
		log.Printf("Running job: %s", job.name)
		if err := job.handler(); err != nil {
			log.Printf("Error running job %s: %v", job.name, err)
		}
	}
}

func scheduledJobs(ghHandler *github.Handler, hnHandler *hackernews.Handler) {
	scheduler := NewJobScheduler()

	// Add GitHub job - runs every hour
	scheduler.AddJob("GitHub Trending", time.Hour, func() error {
		_, err := ghHandler.FetchTrendingRepos()
		return err
	})

	// Add HackerNews job - runs every 15 minutes
	scheduler.AddJob("HackerNews Top", 15*time.Minute, func() error {
		_, err := hnHandler.FetchTopStories()
		return err
	})

	// Start all jobs
	scheduler.Start()
}

func main() {
	// Initialize database
	if err := database.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// GitHub routes with database persistence
	ghHandler := github.NewHandler()
	ghHandler.RegisterRoutes(app)

	// HackerNews routes with database persistence
	hnHandler := hackernews.NewHandler()
	hnHandler.RegisterRoutes(app)

	// Tickers routes
	tickerHandler := tickers.NewHandler()
	tickerHandler.RegisterRoutes(app)

	// Start scheduled jobs in a goroutine
	go scheduledJobs(ghHandler, hnHandler)

	log.Printf("Server starting on port %s\n", port)
	log.Fatal(app.Listen(":" + port))
}
