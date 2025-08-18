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

# Production stage - use nginxinc/nginx-unprivileged for non-root operation
FROM nginxinc/nginx-unprivileged:alpine-slim

# Copy built assets from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Expose port 8080 (nginx-unprivileged uses 8080 instead of 80)
EXPOSE 8080

# Start nginx as non-root user
CMD ["nginx", "-g", "daemon off;"] 