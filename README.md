# Today Dashboard

[Previous content remains the same until the Error Handling section...]

## Error Handling and Resilience

The setup includes several layers of error handling and resilience features:

### 1. Dynamic DNS Resolution

The nginx configuration uses Docker's internal DNS resolver to dynamically resolve service addresses:

```nginx
# Enable Docker DNS resolver
resolver 127.0.0.11 valid=30s ipv6=off;

location / {
    # Use variables to force DNS re-resolution
    set $frontend_upstream "http://today-frontend";
    proxy_pass $frontend_upstream:80;
    # ... other configurations
}
```

This setup:
- Allows services to come and go without nginx restart
- Automatically detects when services become available
- Handles container recreation seamlessly
- Prevents nginx from failing when services are down

### 2. Service Level Resilience

Each service container includes:
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:80/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 10s
restart: unless-stopped
```

Benefits:
- Automatic container health monitoring
- Self-healing through automatic restarts
- Graceful startup and shutdown

### 3. Nginx Error Handling

Multiple layers of error protection:
```nginx
# Timeouts prevent hanging connections
proxy_connect_timeout 5s;
proxy_send_timeout 10s;
proxy_read_timeout 10s;

# Automatic retry on failures
proxy_next_upstream error timeout http_502 http_503 http_504;
proxy_next_upstream_timeout 0;
proxy_next_upstream_tries 0;

# Custom error pages
error_page 502 503 504 /error.html;
```

This ensures:
- Fast failure detection
- Graceful error handling
- User-friendly error pages
- No cascading failures

### 4. Network Resilience

The setup uses multiple networks:
```yaml
networks:
  default:    # Service-specific network
    name: today-network
    driver: bridge
  shared-web: # Shared nginx network
    external: true
    name: shared-web
```

Providing:
- Network isolation
- Service discovery
- Load balancing
- Failure isolation

## Maintenance and Recovery

### 1. Normal Operation
- Services can be started/stopped independently
- Nginx continues running even when services are down
- DNS resolution automatically handles service recovery

### 2. Service Updates
```bash
# Update a service without affecting others
docker-compose up -d --build frontend
```

### 3. Nginx Updates
```bash
cd nginx
docker-compose up -d --build
```

### 4. Full System Recovery
```bash
# Start shared infrastructure
cd nginx
docker-compose up -d

# Start application services
cd ..
docker-compose up -d
```

## Troubleshooting

### 1. Service Issues
If a service is down:
- Nginx will serve the error page
- Other services remain unaffected
- Service can be restarted independently
```bash
docker-compose restart frontend
```

### 2. DNS Resolution Issues
If services aren't being discovered:
```bash
# Check Docker DNS
docker exec shared-nginx ping today-frontend
```

### 3. Network Issues
```bash
# Verify network connectivity
docker network inspect shared-web

# Check if services are properly connected
docker network inspect today-network
```

### 4. Logs
```bash
# Check nginx logs
docker-compose -f nginx/docker-compose.yml logs -f

# Check service logs
docker-compose logs -f frontend
docker-compose logs -f backend
```

[Rest of the previous content remains the same...]
