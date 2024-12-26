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

## Development

1. Install dependencies:
```bash
npm install
```

2. Create environment configuration:
```bash
cp .env.example .env
```

3. Start the development server:
```bash
npm run dev
```

The application will be available at http://localhost:5173 during development.

## Production Deployment

1. Build the frontend:
```bash
npm run build
```

2. Start all services:
```bash
docker-compose up -d
```

The application uses Docker Compose for orchestrating all services in production, with Caddy handling routing and HTTPS.
