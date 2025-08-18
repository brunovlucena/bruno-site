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

# Production stage - use regular nginx but configure for non-root
FROM nginx:alpine

# Create nginx user and group if they don't exist, or modify existing ones
RUN if ! getent group nginx > /dev/null 2>&1; then \
        addgroup -g 101 -S nginx; \
    else \
        groupmod -g 101 nginx; \
    fi && \
    if ! getent passwd nginx > /dev/null 2>&1; then \
        adduser -S -D -H -u 101 -h /var/cache/nginx -s /sbin/nologin -G nginx -g nginx nginx; \
    else \
        usermod -u 101 -g 101 nginx; \
    fi

# Create necessary directories with correct permissions
RUN mkdir -p /var/cache/nginx/client_temp \
    /var/cache/nginx/fastcgi_temp \
    /var/cache/nginx/proxy_temp \
    /var/cache/nginx/scgi_temp \
    /var/cache/nginx/uwsgi_temp \
    /var/log/nginx \
    /var/run \
    /tmp/nginx_proxy \
    /tmp/nginx_fastcgi \
    /tmp/nginx_uwsgi \
    /tmp/nginx_scgi && \
    chown -R nginx:nginx /var/cache/nginx \
    /var/log/nginx \
    /var/run \
    /tmp/nginx_proxy \
    /tmp/nginx_fastcgi \
    /tmp/nginx_uwsgi \
    /tmp/nginx_scgi && \
    chmod -R 755 /var/cache/nginx \
    /var/log/nginx \
    /var/run \
    /tmp/nginx_proxy \
    /tmp/nginx_fastcgi \
    /tmp/nginx_uwsgi \
    /tmp/nginx_scgi

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html
RUN chown -R nginx:nginx /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf
RUN chown nginx:nginx /etc/nginx/nginx.conf

# Expose port 80
EXPOSE 80

# Switch to nginx user and start nginx
USER nginx
CMD ["nginx", "-g", "daemon off;"] 