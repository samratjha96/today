# EC2 Deployment Guide

## First Time Setup

1. SSH into your EC2 instance:
```bash
ssh ec2-user@YOUR_EC2_IP
```

2. Clone and prepare the repository:
```bash
# Clone the repository
git clone <repository-url>
cd today

# Make setup script executable
chmod +x get-ip.sh

# Run setup script
./get-ip.sh
```

3. Start the application:
```bash
# Build and start containers
docker-compose up -d

# Verify containers are running
docker-compose ps
```

4. Access the application:
- Open your browser
- Go to http://YOUR_EC2_IP
- Both frontend and backend should be working

## After EC2 Restart

If you restart your EC2 instance, you'll need to:
```bash
# 1. Update the backend URL with new IP
./get-ip.sh

# 2. Restart containers
docker-compose up -d
```

## Troubleshooting

1. Check container status:
```bash
docker-compose ps
```

2. View logs:
```bash
# All logs
docker-compose logs

# Frontend logs
docker-compose logs frontend

# Backend logs
docker-compose logs backend
```

3. Test backend directly:
```bash
# Get your EC2 IP
curl http://169.254.169.254/latest/meta-data/public-ipv4

# Test health endpoint
curl http://YOUR_EC2_IP:8020/health
```

4. Rebuild everything:
```bash
# Stop containers
docker-compose down

# Update IP and rebuild
./get-ip.sh
docker-compose up -d --build
```

## Common Issues

1. "ERR_BLOCKED_BY_CLIENT" in browser console:
   - This usually means the backend URL is wrong
   - Run ./get-ip.sh to update it
   - Rebuild and restart containers

2. Backend not responding:
   - Check if container is running: `docker-compose ps`
   - Check logs: `docker-compose logs backend`
   - Test directly: `curl localhost:8020/health`

3. Frontend not loading:
   - Check container: `docker-compose logs frontend`
   - Verify port 80: `curl localhost`
   - Check .env: `cat .env`
