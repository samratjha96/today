# Quick Test Guide

## Test on EC2

1. Start containers:
```bash
docker-compose up -d
```

2. Check if containers are running:
```bash
docker ps
```

3. Check container logs if there are issues:
```bash
docker logs today-frontend
docker logs today-backend
```

4. Test the app:
```bash
# Get your EC2 public IP
curl http://169.254.169.254/latest/meta-data/public-ipv4

# Test in browser
http://YOUR_EC2_PUBLIC_IP

# Or test with curl
curl localhost
```

## Cleanup if needed

```bash
# Stop everything
docker-compose down

# Start fresh
docker-compose up -d
```

## Common Issues

1. If containers won't start:
```bash
# Check logs
docker logs today-frontend
docker logs today-backend

# Rebuild containers
docker-compose up -d --build
```

2. If frontend can't connect to backend:
```bash
# Check backend is running
curl localhost:8020/health

# Check backend logs
docker logs today-backend
