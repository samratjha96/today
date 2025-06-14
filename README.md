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

### Backend Services
Multiple backend services handle different aspects of data fetching and processing:

1. Go Backend (`/backend/go-backend`)
   - Handles GitHub trending repositories
   - Processes Hacker News data
   - Built with modern Go practices and modular architecture

2. Python Backend (`/backend/python-backend`)
   - Manages market data and stock information
   - Processes RSS feeds
   - Uses FastAPI for efficient API endpoints

### Reverse Proxy with Caddy

The application uses Caddy v2 as a reverse proxy to handle routing and load balancing. The key feature is its ability to intelligently route `/api` requests across multiple backend services.

#### API Request Handling

The Caddy configuration demonstrates sophisticated request routing:

```caddy
handle /api/* {
    # Strip the /api prefix before forwarding to backends
    uri strip_prefix /api

    # Route to appropriate backend based on path
    reverse_proxy {
        # Dynamic backend selection based on request path
        to today-backend-1:8020 today-backend-2:8020

        # Load balancing configuration
        lb_policy round_robin

        # Health checks
        health_uri /health
        health_interval 30s
        health_timeout 10s
    }
}
```

Key Features:
- Path-based routing to different backend services
- Automatic load balancing across backend instances
- Health checking for backend services
- Clean API URL structure by stripping prefixes
- Secure headers and HTTPS handling

## Directory Structure

```
/
├── src/                    # Frontend React application
│   ├── components/        # React components
│   ├── hooks/            # Custom React hooks
│   ├── lib/              # Utilities and constants
│   └── pages/            # Page components
├── backend/
│   ├── go-backend/       # Go services
│   └── python-backend/   # Python services
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
