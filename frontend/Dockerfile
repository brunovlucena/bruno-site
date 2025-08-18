# Build stage
FROM node:22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy package files first for better caching
COPY package*.json ./

# Install all dependencies (including dev dependencies for build)
RUN npm ci --legacy-peer-deps && npm cache clean --force

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage - use nginx without user restrictions
FROM nginx:stable-alpine

# Switch to root user
USER root

# Create nginx cache directories with proper permissions
RUN mkdir -p /tmp/cache/nginx/{client_temp,proxy_temp,fastcgi_temp,uwsgi_temp,scgi_temp} && \
    chmod -R 755 /tmp/cache/nginx

# Copy nginx configuration first
COPY nginx.conf /etc/nginx/nginx.conf

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Expose port 8080
EXPOSE 8080

# Start nginx as root (required for binding to port 8080)
CMD ["nginx", "-g", "daemon off;"] 