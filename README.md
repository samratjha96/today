# Today Dashboard

A real-time dashboard showing market data and tech news.

## Infrastructure Setup

### 1. Shared Nginx Setup

The infrastructure uses a shared nginx container that routes traffic to multiple services based on domain names.

1. Create the shared nginx directory structure:
```bash
mkdir -p nginx/conf.d nginx/html
```

2. Start the shared nginx service:
```bash
cd nginx
docker-compose up -d
```

The shared nginx service configuration (`nginx/docker-compose.yml`):
```yaml
version: '3.8'

services:
  nginx:
    image: nginx:alpine
    container_name: shared-nginx
    volumes:
      - ./conf.d:/etc/nginx/conf.d
      - ./html:/usr/share/nginx/html
    ports:
      - "80:80"
    networks:
      - shared-web
    healthcheck:
      test: ["CMD", "nginx", "-t"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
    restart: unless-stopped

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

### 2. Nginx Configuration

The nginx configuration includes error handling and resilience features to prevent service disruptions:

```nginx
# Upstream definitions with failure handling
upstream today-frontend-upstream {
    server today-frontend:80 max_fails=3 fail_timeout=10s;
    server 127.0.0.1:1 backup;  # Fallback for graceful failure
}

upstream today-backend-upstream {
    server today-backend:8020 max_fails=3 fail_timeout=10s;
    server 127.0.0.1:1 backup;  # Fallback for graceful failure
}

server {
    listen 80;
    server_name today.techbrohomelab.xyz;
    
    # Custom error handling
    error_page 502 503 504 /error.html;
    
    location / {
        proxy_pass http://today-frontend-upstream;
        # Timeouts and error handling
        proxy_connect_timeout 5s;
        proxy_send_timeout 10s;
        proxy_read_timeout 10s;
        proxy_next_upstream error timeout http_502 http_503 http_504;
        proxy_intercept_errors on;
    }
    
    location /api/ {
        proxy_pass http://today-backend-upstream/;
        # Similar timeout and error handling settings
    }
}
```

Key resilience features:
- Upstream failure detection (`max_fails=3 fail_timeout=10s`)
- Graceful degradation with backup servers
- Custom error pages
- Connection timeouts
- Automatic retry on failures

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

### 4. Deployment

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

1. Create a new nginx configuration in `nginx/conf.d/`:
```nginx
# Example for a new service
upstream newapp-frontend-upstream {
    server newapp-frontend:80 max_fails=3 fail_timeout=10s;
    server 127.0.0.1:1 backup;
}

server {
    listen 80;
    server_name newapp.techbrohomelab.xyz;
    # ... rest of the configuration similar to today.conf
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
# Reload nginx to pick up new configuration
cd nginx
docker-compose restart

# Start your new service
cd /path/to/your/service
docker-compose up -d
```

## Troubleshooting

1. Check nginx logs:
```bash
cd nginx
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

4. Verify nginx configuration:
```bash
cd nginx
docker-compose exec nginx nginx -t
```

5. Service Unavailability:
- If a service is down, nginx will display a custom error page
- Check the specific service logs: `docker logs <container-name>`
- The shared nginx will continue running even if individual services are down
- Other services will remain unaffected by one service's failure

## Error Handling

The setup includes several layers of error handling:

1. Service Level:
- Health checks for both frontend and backend
- Automatic container restarts
- Graceful shutdown handling

2. Nginx Level:
- Custom error pages
- Upstream failure detection
- Connection timeouts
- Automatic request retries
- Graceful degradation

3. Network Level:
- Isolated service networks
- Shared proxy network
- Cloudflare SSL/TLS protection

This multi-layered approach ensures high availability and graceful degradation when issues occur.
