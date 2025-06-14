# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Frontend (React + TypeScript)

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Build for development
npm run build:dev

# Lint code
npm run lint

# Preview production build
npm run preview
```

### Go Backend

```bash
# Navigate to Go backend directory
cd backend/go-backend

# Run the Go backend
go run main.go

# Build the Go backend
go build -o app main.go
```

### Docker

```bash
# Start all services with Docker Compose
docker-compose up -d

# Build and start specific service
docker-compose up -d --build <service-name>

# View logs for all services
docker-compose logs -f

# View logs for specific service
docker-compose logs -f <service-name>

# Stop all services
docker-compose down
```

## System Architecture

The Today Dashboard is a microservice-based application with two main components:

1. **Frontend**: React/TypeScript SPA built with Vite
   - Uses React Query for data fetching and state management
   - Implements responsive UI with Tailwind CSS
   - Uses shadcn/ui components for UI elements

2. **Go Backend**: Consolidated backend service handling all data sources
   - Built with Go Fiber framework
   - Uses SQLite for data persistence
   - Implements periodic data fetching with a job scheduler
   - Endpoints:
     - `/github/trending`: GitHub trending repositories
     - `/hackernews/top`: HackerNews top stories
     - `/tickers`: Financial market data (ticker information)
     - `/news`: Tech news from various RSS feeds

3. **Caddy Reverse Proxy**: Handles routing between services
   - Routes all `/api/*` requests to the Go backend
   - Routes all other requests to the frontend
   - Adds security headers and handles errors

## Data Flow

1. Frontend makes API requests to `/api/*` endpoints
2. Caddy routes requests to the Go backend service
3. Backend service fetches, processes, and returns data
4. React components render data with React Query handling caching and refetching

## Key Components

### Frontend Components

- **TickerTable**: Displays financial market data
- **GithubTrending**: Shows trending GitHub repositories
- **HackerNews**: Displays top stories from Hacker News
- **RSSFeed**: Shows tech news from various RSS feeds

### Backend Services

- **Go Backend**:
  - Job Scheduler: Periodically fetches data from various sources
  - Data Storage: SQLite database for persisting fetched data
  - API Handlers: Process and serve data through REST endpoints
  - Modules:
    - github: Handles GitHub trending repositories
    - hackernews: Handles HackerNews top stories
    - tickers: Handles financial market data
    - rss: Handles RSS feeds from tech news sources

## Configuration

- Environment variables control backend URLs and API modes:
  - `VITE_BACKEND_URL`: Backend API URL (defaults to "/api")
  - `VITE_API_MODE`: Can be set to "mock" for development without backends
  - `ALLOWED_HOSTS`: Comma-separated list of allowed hosts for CORS

## Docker Deployment

- Each component is containerized with its own Dockerfile
- Docker Compose orchestrates all services:
  - `frontend`: React application served with serve
  - `go-backend`: Consolidated Go service for all data sources
  - `caddy`: Reverse proxy for routing requests

- Shared network `shared-web` connects all services
- Health checks ensure service availability