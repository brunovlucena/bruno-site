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
RUN mkdir -p /tmp/nginx_proxy /tmp/nginx_fastcgi /tmp/nginx_uwsgi /tmp/nginx_scgi /var/cache/nginx && \
    chmod 755 /tmp/nginx_proxy /tmp/nginx_fastcgi /tmp/nginx_uwsgi /tmp/nginx_scgi /var/cache/nginx && \
    # Remove default nginx files to save space
    rm -rf /usr/share/nginx/html/* && \
    # Remove unnecessary nginx modules and files
    rm -rf /var/cache/apk/* /tmp/* /var/tmp/*

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Create non-root user for security (nginx user already exists in alpine-slim)
RUN addgroup -g 1001 -S nginx-user && \
    adduser -S -D -H -u 1001 -h /var/cache/nginx -s /sbin/nologin -G nginx-user -g nginx-user nginx-user

# Change ownership of nginx directories
RUN chown -R nginx-user:nginx-user /var/cache/nginx /var/log/nginx /etc/nginx/conf.d /usr/share/nginx/html

# Expose port
EXPOSE 80

# Switch to non-root user
USER nginx-user

# Start nginx
CMD ["nginx", "-g", "daemon off;"] 