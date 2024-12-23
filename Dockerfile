# Build stage
FROM node:20-alpine as build

WORKDIR /app

# Copy package files
COPY package*.json ./
COPY bun.lockb ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build the app
RUN npm run build

# Runtime stage
FROM nginx:alpine

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy built files from build stage
COPY --from=build /app/dist /usr/share/nginx/html

# Create a script to replace environment variables in the built files
RUN echo '#!/bin/sh' > /docker-entrypoint.sh && \
    echo 'for file in /usr/share/nginx/html/assets/*.js; do' >> /docker-entrypoint.sh && \
    echo '  if [ -f "$file" ]; then' >> /docker-entrypoint.sh && \
    echo '    envsubst '\''${VITE_BACKEND_URL} ${VITE_API_MODE}'\'' < "$file" > "$file.tmp" && mv "$file.tmp" "$file"' >> /docker-entrypoint.sh && \
    echo '  fi' >> /docker-entrypoint.sh && \
    echo 'done' >> /docker-entrypoint.sh && \
    echo 'nginx -g "daemon off;"' >> /docker-entrypoint.sh && \
    chmod +x /docker-entrypoint.sh

# Expose port 80
EXPOSE 80

# Set environment variables with defaults
ENV VITE_BACKEND_URL=http://localhost:8020
ENV VITE_API_MODE=real

# Run the entrypoint script
ENTRYPOINT ["/docker-entrypoint.sh"]
