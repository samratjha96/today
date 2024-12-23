# Quick Test Guide

## Configuration Files

1. Container Nginx Files:
   - `nginx.conf` - Base configuration for the frontend container
   - `default.conf` - Server block configuration for the frontend container

2. Host Nginx Files:
   - `nginx-host-main.conf` - Main configuration for the host machine
   - `nginx-host.conf` - Virtual host configuration for the host machine

## Testing Steps

1. Start the containers:
```bash
# Start containers in background
docker-compose up -d

# Verify containers are running
docker ps

# Check frontend container logs
docker logs today-frontend
```

2. Set up host nginx:
```bash
# Copy main nginx config
sudo cp nginx-host-main.conf /etc/nginx/nginx.conf

# Copy virtual host config
sudo cp nginx-host.conf /etc/nginx/sites-available/default
sudo ln -sf /etc/nginx/sites-available/default /etc/nginx/sites-enabled/default

# Create cache directory if it doesn't exist
sudo mkdir -p /var/cache/nginx

# Test config
sudo nginx -t

# Restart nginx
sudo systemctl restart nginx
```

3. Test the routing:
```bash
# Test nginx is responding
curl localhost

# Test frontend container directly
curl localhost:8019

# Test backend health
curl localhost:8020/health

# If you want to see headers
curl -v localhost

# If you want to see the full response
curl -i localhost
```

4. Browser Testing:
```bash
# Get your EC2 instance's public IP
curl http://169.254.169.254/latest/meta-data/public-ipv4

# Open in your browser
http://YOUR_EC2_PUBLIC_IP
```

The app should now work in your browser because:
- Host nginx is routing all traffic on port 80 to the frontend container (8019)
- Frontend container's nginx is serving the Vite app
- Backend CORS is temporarily allowing all origins
- Backend container is serving the API (8020)

5. Troubleshooting:

Container nginx issues:
```bash
# Check frontend container logs
docker logs today-frontend

# Inspect frontend container
docker exec -it today-frontend sh

# Check nginx config inside container
docker exec today-frontend nginx -t

# Check nginx config files inside container
docker exec today-frontend cat /etc/nginx/nginx.conf
docker exec today-frontend cat /etc/nginx/conf.d/default.conf
```

Host nginx issues:
```bash
# Check host nginx logs
sudo tail -f /var/log/nginx/error.log
sudo tail -f /var/log/nginx/access.log

# Check host nginx config
sudo nginx -t

# Check host nginx status
sudo systemctl status nginx
```

Other checks:
```bash
# Check container logs
docker logs today-backend

# Check listening ports
sudo netstat -tulpn | grep '80\|8019\|8020'

# Test API endpoints directly
curl http://localhost:8020/tickers
curl http://localhost:8020/news
```

6. If things go wrong, here's how to start fresh:
```bash
# Stop and remove containers
docker-compose down

# Remove host nginx configs
sudo rm /etc/nginx/nginx.conf
sudo rm /etc/nginx/sites-available/default
sudo rm /etc/nginx/sites-enabled/default

# Restore default host nginx config
sudo cp /etc/nginx/nginx.conf.backup /etc/nginx/nginx.conf

# Stop host nginx
sudo systemctl stop nginx

# Start fresh
sudo systemctl start nginx
docker-compose up -d
```

7. After Testing:
Once everything is working, you should:
1. Update backend/main.py to restrict CORS to your actual domain
2. Update nginx-host.conf with your domain configuration
3. Set up SSL with Certbot
4. Rebuild and redeploy the containers

The temporary CORS setting (`origins = ["*"]`) in backend/main.py should be replaced with your actual domains:
```python
origins = [
    "http://today.samratjha.com",
    "https://today.samratjha.com"
]
```

## Understanding the Nginx Setup

1. Container Nginx:
   - Base config (nginx.conf):
     * Sets up basic settings
     * Includes mime types
     * Configures logging
     * Enables gzip
   - Server block (default.conf):
     * Serves the built Vite app
     * Handles SPA routing
     * Manages static assets

2. Host Nginx:
   - Main config (nginx-host-main.conf):
     * System-wide settings
     * Worker processes
     * Global security
   - Virtual hosts (nginx-host.conf):
     * Domain routing
     * Proxy settings
     * SSL configuration (when added)
