# Build stage
FROM node:22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy package files first for better caching
COPY package*.json ./

# Install all dependencies (including dev dependencies for build)
RUN npm ci && npm cache clean --force

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage - use nginx:alpine-slim for smaller size
FROM nginx:alpine-slim

# Create necessary directories and set permissions
RUN mkdir -p /tmp/nginx_proxy /tmp/nginx_fastcgi /tmp/nginx_uwsgi /tmp/nginx_scgi /var/cache/nginx /run && \
    chmod 755 /tmp/nginx_proxy /tmp/nginx_fastcgi /tmp/nginx_uwsgi /tmp/nginx_scgi /var/cache/nginx /run && \
    chown -R nginx:nginx /run && \
    # Remove default nginx files to save space
    rm -rf /usr/share/nginx/html/* && \
    # Remove unnecessary nginx modules and files
    rm -rf /var/cache/apk/* /tmp/* /var/tmp/*

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port
EXPOSE 80

# Start nginx as nginx user
USER nginx
CMD ["nginx", "-g", "daemon off;"] 