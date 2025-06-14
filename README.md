# Today Dashboard

A real-time dashboard showing market data and tech news.

## Architecture Overview

The Today Dashboard is built with a modern microservices architecture:

### Frontend
- React + TypeScript application built with Vite
- Located in `/src` directory
- Key components:
  - Market data visualization (TickerCard, TickerTable)
  - Tech news aggregation (HackerNews, GithubTrending, RSSFeed)
  - Reusable UI components in `/src/components/ui`
- Custom hooks for data fetching and state management in `/src/hooks`
- Utility functions and constants in `/src/lib`

### Backend Service

A consolidated Go backend handles all data fetching and processing:

* Go Backend (`/backend/go-backend`)
  - Built with Go Fiber framework
  - Uses SQLite for data persistence
  - Implements a job scheduler for periodic data fetching
  - Key modules:
    - `github`: Scrapes and serves GitHub trending repositories
    - `hackernews`: Fetches and processes Hacker News stories
    - `tickers`: Retrieves financial market data from Yahoo Finance
    - `rss`: Aggregates tech news from various RSS feeds

### Reverse Proxy with Caddy

The application uses Caddy v2 as a reverse proxy to handle routing and load balancing. The key feature is its ability to intelligently route `/api` requests across multiple backend services.

#### API Request Handling

The Caddy configuration provides efficient routing:

```caddy
handle /api/* {
    # Strip the /api prefix before forwarding to backend
    uri strip_prefix /api

    # Route to Go backend
    reverse_proxy today-go-backend:3001 {
        # Health checks
        health_uri /health
        health_interval 30s
        health_timeout 10s
    }
}
```

Key Features:
- Clean API URL structure by stripping prefixes
- Health checking for backend service
- Secure headers and HTTPS handling
- Efficient routing of all API requests to the Go backend

## Directory Structure

```
/
├── src/                    # Frontend React application
│   ├── components/        # React components
│   ├── hooks/            # Custom React hooks
│   ├── lib/              # Utilities and constants
│   └── pages/            # Page components
├── backend/
│   └── go-backend/       # Consolidated Go backend service
│       ├── pkg/          # Backend modules (github, hackernews, rss, tickers)
│       └── data/         # SQLite database directory
└── caddy/                # Caddy reverse proxy configuration
```

## Development & Deployment

This project uses a Makefile to simplify development and deployment workflows. Below are the available commands:

### Available Make Commands

```bash
# View all available commands
make help

# Set up development environment and start servers
make dev

# Only install dependencies without starting servers
make dev-deps

# Build and start all services in production mode
make prod

# Clean up all build artifacts and containers
make clean

# Run tests
make test
```

### Development

For development, simply run:

```bash
make dev
```

This will:
- Install all necessary dependencies (frontend and backend)
- Start the Go backend on port 3001
- Start the frontend development server on port 8019

The application will be available at http://localhost:8019 during development.

### Production Deployment

For production deployment, use:

```bash
make prod
```

This will:
- Create required Docker networks
- Build all containers
- Start the entire stack including Caddy reverse proxy

The application will be available at http://localhost in production mode with Caddy handling routing and HTTPS.
