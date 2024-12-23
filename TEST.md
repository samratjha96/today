# Quick Test Guide

Run these commands in order to test the setup:

1. Start the containers:
```bash
# Start containers in background
docker-compose up -d

# Verify containers are running
docker ps
```

2. Set up nginx:
```bash
# Copy nginx configs
sudo cp nginx.conf /etc/nginx/nginx.conf
sudo cp nginx-host.conf /etc/nginx/sites-available/default
sudo ln -sf /etc/nginx/sites-available/default /etc/nginx/sites-enabled/default

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
- nginx is routing all traffic on port 80 to the frontend
- Backend CORS is temporarily allowing all origins
- Frontend container is serving the Vite app
- Backend container is serving the API

5. Troubleshooting:
```bash
# Check nginx logs
sudo tail -f /var/log/nginx/error.log
sudo tail -f /var/log/nginx/access.log

# Check container logs
docker logs today-frontend
docker logs today-backend

# Check nginx status
sudo systemctl status nginx

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

# Remove nginx configs
sudo rm /etc/nginx/nginx.conf
sudo rm /etc/nginx/sites-available/default
sudo rm /etc/nginx/sites-enabled/default

# Restore default nginx config
sudo cp /etc/nginx/nginx.conf.backup /etc/nginx/nginx.conf

# Stop nginx
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
