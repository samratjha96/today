# Today Dashboard

A real-time dashboard showing market data and tech news.

## Host Machine Setup (Ubuntu EC2)

### Install and Configure Nginx

1. Update system and install nginx:
```bash
# Update package list
sudo apt update

# Install nginx
sudo apt install nginx -y

# Start nginx and enable it to start on boot
sudo systemctl start nginx
sudo systemctl enable nginx

# Check status
sudo systemctl status nginx
```

2. Set up nginx configuration:
```bash
# Create cache directory
sudo mkdir -p /var/cache/nginx

# Backup default configuration
sudo mv /etc/nginx/nginx.conf /etc/nginx/nginx.conf.backup

# Copy our main configuration
sudo cp nginx.conf /etc/nginx/nginx.conf

# Remove default site configuration
sudo rm /etc/nginx/sites-enabled/default

# Copy our site configuration
sudo cp nginx-host.conf /etc/nginx/sites-available/samratjha.com
sudo ln -s /etc/nginx/sites-available/samratjha.com /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# If test passes, reload nginx
sudo systemctl reload nginx
```

3. Set up logging:
```bash
# Create log directory if it doesn't exist
sudo mkdir -p /var/log/nginx

# Set proper permissions
sudo chown -R www-data:adm /var/log/nginx
```

4. Configure firewall:
```bash
# Allow HTTP traffic
sudo ufw allow 80/tcp

# Allow HTTPS traffic (if using SSL)
sudo ufw allow 443/tcp

# Enable firewall if not already enabled
sudo ufw enable
```

5. (Optional) Set up SSL with Certbot:
```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx -y

# Get SSL certificate
sudo certbot --nginx -d today.samratjha.com

# Enable auto-renewal
sudo systemctl enable certbot.timer
sudo systemctl start certbot.timer
```

### Nginx Configuration Details

#### Main Configuration (nginx.conf)

The main nginx configuration file includes optimized settings for performance and security:

1. Worker Processes:
```nginx
worker_processes auto;  # Automatically sets worker count based on CPU cores
```
- For high-traffic sites, you might want to set this manually to CPU core count minus one

2. Event Settings:
```nginx
events {
    worker_connections 1024;  # Maximum connections per worker
    multi_accept on;         # Accept multiple connections per event
    use epoll;              # Efficient event processing method
}
```
- Increase worker_connections for high-traffic sites
- Default 1024 is good for most cases

3. HTTP Settings:
```nginx
http {
    sendfile on;           # Efficient file sending
    tcp_nopush on;        # Optimize packet sending
    tcp_nodelay on;       # Reduce latency
    keepalive_timeout 65; # Keep connections alive
}
```

4. SSL Configuration:
```nginx
ssl_protocols TLSv1.2 TLSv1.3;              # Modern SSL protocols
ssl_prefer_server_ciphers on;               # Prefer server ciphers
ssl_session_cache shared:SSL:10m;           # SSL session cache
```

5. Gzip Compression:
```nginx
gzip on;
gzip_types text/plain text/css application/json application/javascript;
```
- Reduces bandwidth usage
- Improves load times

6. Rate Limiting:
```nginx
limit_req_zone $binary_remote_addr zone=one:10m rate=10r/s;
```
- Protects against DDoS attacks
- Adjust rate based on your needs

### Performance Tuning

1. For high-traffic sites:
```nginx
worker_processes auto;
worker_rlimit_nofile 20000;
events {
    worker_connections 4096;
    multi_accept on;
}
```

2. For large file transfers:
```nginx
client_max_body_size 100M;
client_body_buffer_size 128k;
```

3. For better caching:
```nginx
open_file_cache max=1000 inactive=20s;
open_file_cache_valid 30s;
open_file_cache_min_uses 2;
open_file_cache_errors on;
```

### Security Best Practices

1. Hide nginx version:
```nginx
server_tokens off;
```

2. Secure headers:
```nginx
add_header X-Frame-Options "SAMEORIGIN";
add_header X-XSS-Protection "1; mode=block";
add_header X-Content-Type-Options "nosniff";
```

3. SSL configuration:
```nginx
ssl_protocols TLSv1.2 TLSv1.3;
ssl_prefer_server_ciphers on;
ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256;
```

[Rest of the existing README content remains the same...]
