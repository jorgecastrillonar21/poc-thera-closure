# TheraClosure Geolocation Service Dockerfile
# Multi-stage build with Python support for comprehensive world data seeding

# Stage 1: Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o geolocation-service \
    ./cmd/main.go

# Stage 2: Production stage
FROM alpine:3.18

# Install runtime dependencies including Python for world data seeding
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    python3 \
    py3-pip \
    curl \
    bash \
    && rm -rf /var/cache/apk/*

# Set timezone
ENV TZ=UTC
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Create app user for security
RUN addgroup -g 1001 -S theraclosure && \
    adduser -u 1001 -S theraclosure -G theraclosure

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/geolocation-service .

# Copy data seeding scripts and configurations
COPY --from=builder /app/data ./data

# Install Python dependencies for world data seeding
COPY data/seeds/requirements.txt /tmp/requirements.txt
RUN pip3 install --no-cache-dir -r /tmp/requirements.txt && \
    rm /tmp/requirements.txt

# Make seeding script executable
RUN chmod +x /app/data/seeds/run_seeder.sh

# Create logs directory and set permissions
RUN mkdir -p /app/logs && \
    chown -R theraclosure:theraclosure /app

# Switch to non-root user
USER theraclosure

# Expose service port
EXPOSE 3003

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:3003/health || exit 1

# Command to run the service
CMD ["./geolocation-service"]