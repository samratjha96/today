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

## Getting Started

### Prerequisites

Before running the application, ensure you have the following installed:

**Required for Production:**
- [Docker](https://docs.docker.com/get-docker/) (20.10+)
- [Docker Compose](https://docs.docker.com/compose/install/) (v2+)

**Required for Development:**
- All of the above, plus:
- [Node.js](https://nodejs.org/) (20+)
- [Go](https://go.dev/doc/install) (1.23+)
- npm (comes with Node.js)

### Quick Start (Fresh Install)

```bash
# Check if all prerequisites are installed
make check-prereqs

# Complete fresh installation setup
make install

# For development mode
make dev

# OR for production mode
make prod
```

That's it! The Makefile handles everything automatically.

---

## Development & Deployment

This project uses a comprehensive Makefile to simplify all workflows. The Makefile is designed to work on completely fresh installations.

### 🔍 Check Prerequisites

Before starting, verify all required tools are installed:

```bash
make check-prereqs
```

This checks for Docker, Docker Compose, Node.js, npm, and Go.

### 🚀 Fresh Installation

For a complete fresh setup (first time):

```bash
make install
```

This command:
1. ✅ Checks all prerequisites
2. ✅ Creates environment files from `.env.example`
3. ✅ Creates Docker network (`shared-web`)
4. ✅ Prepares the system for deployment

### 💻 Development Mode

Run the application locally for development:

```bash
make dev
```

This will:
- Install all npm and Go dependencies
- Start Go backend on `http://localhost:3001`
- Start React frontend on `http://localhost:5173`
- Enable hot-reload for both services

**Individual Services:**
```bash
make dev-frontend  # Frontend only (port 5173)
make dev-backend   # Backend only (port 3001)
```

### 🐳 Production Deployment

Deploy using Docker containers:

```bash
# Complete production deployment from scratch
make prod

# OR step-by-step:
make install      # Initial setup
make build        # Build Docker images
make up           # Start services
```

**Production URLs:**
- Frontend: `http://localhost`
- API: `http://localhost/api`

**Service Management:**
```bash
make up          # Start all services
make down        # Stop all services
make restart     # Restart all services
make status      # Check service status
make health      # Check health of services
```

### 📊 Monitoring & Logs

View logs from running services:

```bash
make logs              # All services
make logs-frontend     # Frontend only
make logs-backend      # Backend only
make logs-caddy        # Caddy only
```

### 🧪 Testing

Run tests across the application:

```bash
make test              # All tests
make test-backend      # Backend tests only
```

### 🔨 Rebuilding

Rebuild specific services without full cleanup:

```bash
make rebuild           # Rebuild everything
make rebuild-frontend  # Frontend only
make rebuild-backend   # Backend only
make rebuild-caddy     # Caddy only
```

### 🗄️ Database Management

Backup and restore the SQLite database:

```bash
# Backup database
make db-backup

# Restore from backup
make db-restore BACKUP_FILE=backups/today-20231201-120000.db
```

### 🧹 Cleanup

Remove build artifacts and containers:

```bash
make clean         # Remove containers and build artifacts
make clean-all     # Complete cleanup including volumes and images
```

### 🔧 Advanced Commands

```bash
make shell-frontend    # Open shell in frontend container
make shell-backend     # Open shell in backend container
make shell-caddy       # Open shell in Caddy container
make validate-caddy    # Validate Caddyfile configuration
make update-deps       # Update all dependencies
```

### 📋 All Available Commands

View the complete list of commands:

```bash
make help
```

---

## Architecture Deep Dive

### Request Flow

```
Client Request
     ↓
[Caddy Reverse Proxy] :80/:443
     ↓
     ├─→ /api/* → [Go Backend] :3001 (internal)
     │               ↓
     │          [SQLite Database]
     │
     └─→ /* → [React Frontend] :80 (internal)
```

### How Caddy Routes Requests

The Caddy configuration (`caddy/Caddyfile`) implements intelligent routing:

1. **API Requests** (`/api/*`):
   - Strips the `/api` prefix
   - Forwards to `today-go-backend:3001`
   - Example: `http://localhost/api/tickers` → `http://today-go-backend:3001/tickers`

2. **Frontend Requests** (everything else):
   - Forwards to `today-frontend:80`
   - Serves the React SPA
   - Handles client-side routing

3. **Security Headers**:
   - X-Content-Type-Options
   - X-XSS-Protection
   - X-Frame-Options
   - Referrer-Policy

### Docker Network Architecture

All services communicate through a shared Docker network:

```
shared-web (Docker network)
├── shared-caddy (caddy:2-alpine)
├── today-frontend (node:20-alpine)
└── today-go-backend (golang:1.23-alpine → alpine:latest)
```

**Key Points:**
- Services communicate via container names
- Internal ports are not exposed to host
- Only Caddy exposes ports 80/443
- SQLite database persists in Docker volume

### Data Flow & Job Scheduler

The Go backend runs periodic jobs to fetch data:

- **GitHub Trending**: Every 1 hour
- **Hacker News**: Every 15 minutes
- **RSS Feeds**: Configurable intervals
- **Ticker Data**: On-demand (cached)

Data is stored in SQLite and served through REST APIs.
