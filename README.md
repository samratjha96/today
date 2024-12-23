# Today Dashboard

A real-time dashboard showing market data and tech news.

## Infrastructure Setup

### 1. Shared Caddy Setup

The infrastructure uses a shared Caddy v2 container that routes traffic to multiple services based on domain names.

1. Create the Caddy directory structure:
```bash
mkdir -p caddy/error
```

2. Start the shared Caddy service:
```bash
cd caddy
docker-compose up -d
```

The shared Caddy service configuration (`caddy/docker-compose.yml`):
```yaml
version: '3.8'

services:
  caddy:
    image: caddy:2-alpine
    container_name: shared-caddy
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile:ro
      - ./error:/usr/share/caddy/error:ro
      - caddy_data:/data
      - caddy_config:/config
    ports:
      - "80:80"
    networks:
      - shared-web
    healthcheck:
      test: ["CMD", "caddy", "validate", "--config", "/etc/caddy/Caddyfile"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
    restart: unless-stopped

volumes:
  caddy_data:
  caddy_config:

networks:
  shared-web:
    name: shared-web
    driver: bridge
```

### 2. Cloudflare Tunnel Setup

1. Install cloudflared:
```bash
# For Debian/Ubuntu
curl -L --output cloudflared.deb https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
sudo dpkg -i cloudflared.deb
```

2. Authenticate with Cloudflare:
```bash
cloudflared tunnel login
```

3. Create a tunnel:
```bash
cloudflared tunnel create shared-services
```

4. Configure the tunnel in `~/.cloudflared/config.yml`:
```yaml
tunnel: <YOUR-TUNNEL-ID>
credentials-file: /root/.cloudflared/<YOUR-TUNNEL-ID>.json

ingress:
  - hostname: today.techbrohomelab.xyz
    service: http://localhost:80
  - hostname: *.techbrohomelab.xyz
    service: http://localhost:80
  - service: http_status:404
```

5. Start the tunnel as a service:
```bash
sudo cloudflared service install
sudo systemctl start cloudflared
```

## Today App Setup

### 1. Environment Configuration

Create a `.env` file from the template:
```bash
cp .env.example .env
```

Configure the following variables:
```bash
# Backend API URL (for frontend)
VITE_BACKEND_URL=/api

# API Mode (real/mock)
VITE_API_MODE=real

# Allowed hosts for CORS
ALLOWED_HOSTS=today.techbrohomelab.xyz
```

### 2. Caddy Configuration

The Caddy configuration includes security headers and domain-based routing:

```caddy
{
    # Global options
    admin off
    
    servers {
        protocols h2c h2 h1
        timeouts {
            read_body 10s
            read_header 10s
            write 30s
            idle 2m
        }
    }
}

:80 {
    @today host today.techbrohomelab.xyz
    handle @today {
        # Security headers
        header {
            X-Frame-Options "SAMEORIGIN"
            X-XSS-Protection "1; mode=block"
            X-Content-Type-Options "nosniff"
            Content-Security-Policy "default-src 'self' https:; script-src 'self' 'unsafe-inline' 'unsafe-eval' https:; style-src 'self' 'unsafe-inline' https:; img-src 'self' data: https:; font-src 'self' https:; connect-src 'self' https:; frame-ancestors 'none'; base-uri 'self'; form-action 'self'"
            Referrer-Policy "strict-origin-when-cross-origin"
            Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
            Permissions-Policy "accelerometer=(), camera=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), payment=(), usb=()"
        }
        
        handle /api/* {
            reverse_proxy today-backend:8020
            uri strip_prefix /api
        }
        
        handle {
            reverse_proxy today-frontend:80
        }
    }

    handle_errors {
        rewrite * /error.html
        file_server {
            root /usr/share/caddy/error
        }
    }
}
```

Key features:
- Domain-based routing with matchers
- Comprehensive security headers
- Clean API path handling
- Custom error pages
- HTTP/2 support

### 3. App Services Configuration

The app's `docker-compose.yml`:
```yaml
version: '3.8'

services:
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    container_name: today-backend
    expose:
      - "8020"
    environment:
      - TZ=UTC
      - ALLOWED_HOSTS=${ALLOWED_HOSTS:-today.techbrohomelab.xyz}
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8020/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
    networks:
      - default
      - shared-web
    restart: unless-stopped

  frontend:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: today-frontend
    expose:
      - "80"
    environment:
      - VITE_BACKEND_URL=/api
      - VITE_API_MODE=${VITE_API_MODE:-real}
    depends_on:
      backend:
        condition: service_healthy
    networks:
      - default
      - shared-web
    restart: unless-stopped

networks:
  default:
    name: today-network
    driver: bridge
  shared-web:
    external: true
    name: shared-web
```

### 4. Performance Optimizations

#### Frontend Caching Strategy

The app uses React Query for efficient data fetching and caching:

1. RSS News Data (slower changing):
```typescript
{
  refetchInterval: 300000,  // Fetch new data every 5 minutes
  staleTime: 240000,       // Consider data fresh for 4 minutes
  gcTime: 600000,         // Keep in cache for 10 minutes
  refetchOnWindowFocus: false,
  refetchOnReconnect: true
}
```

2. Ticker Data (faster changing):
```typescript
{
  refetchInterval: 30000,  // Fetch new data every 30 seconds
  staleTime: 25000,       // Consider data fresh for 25 seconds
  gcTime: 120000,        // Keep in cache for 2 minutes
  refetchOnWindowFocus: false,
  refetchOnReconnect: true
}
```

Benefits:
- Reduced API calls
- Instant data display from cache
- Background updates without UI flicker
- Automatic retry with exponential backoff
- Optimized network usage

### 5. Deployment

1. Start the app services:
```bash
docker-compose up -d
```

2. Verify the deployment:
- Check container status: `docker ps`
- View logs: `docker-compose logs -f`
- Visit https://today.techbrohomelab.xyz

## Adding New Services

To add a new service to this infrastructure:

1. Add a new host matcher in the Caddyfile:
```caddy
@newapp host newapp.techbrohomelab.xyz
handle @newapp {
    reverse_proxy newapp-frontend:80
}
```

2. Add DNS record in Cloudflare:
- Type: CNAME
- Name: newapp
- Target: your-tunnel-domain.cfargotunnel.com
- Proxy status: Proxied

3. Configure your new service's docker-compose.yml:
- Connect to shared-web network
- Use unique container names
- Expose ports only internally

4. Deploy:
```bash
# Reload Caddy to pick up new configuration
cd caddy
docker-compose restart

# Start your new service
cd /path/to/your/service
docker-compose up -d
```

## Troubleshooting

1. Check Caddy logs:
```bash
cd caddy
docker-compose logs -f
```

2. Check Cloudflare Tunnel status:
```bash
cloudflared tunnel info <YOUR-TUNNEL-NAME>
```

3. Check container connectivity:
```bash
docker network inspect shared-web
```

4. Verify Caddy configuration:
```bash
cd caddy
docker-compose exec caddy caddy validate --config /etc/caddy/Caddyfile
```

5. Service Unavailability:
- If a service is down, Caddy will display a custom error page
- Check the specific service logs: `docker logs <container-name>`
- The shared Caddy will continue running even if individual services are down
- Other services will remain unaffected by one service's failure

## Error Handling

The setup includes several layers of error handling:

1. Service Level:
- Health checks for both frontend and backend
- Automatic container restarts
- Graceful shutdown handling
- React Query retry mechanism

2. Caddy Level:
- Custom error pages
- Security headers
- HTTP/2 support
- Zero-downtime config reloads
- Automatic HTTPS (when not using Cloudflare)

3. Network Level:
- Isolated service networks
- Shared proxy network
- Cloudflare SSL/TLS protection

4. Application Level:
- Data caching with React Query
- Exponential backoff for failed requests
- Graceful degradation with cached data
- Request timeout handling

This multi-layered approach ensures high availability, optimal performance, and graceful degradation when issues occur.
