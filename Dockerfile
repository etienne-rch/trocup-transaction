# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy only the dependency files first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application with additional flags
RUN CGO_ENABLED=0 GOOS=linux go build -o app

# Final stage
FROM alpine:latest

# Set constants for user/group IDs and name
ENV APP_USER=svc_nonroot
ENV APP_UID=10001
ENV APP_GID=10001

# Update Alpine and add necessary packages
RUN apk --no-cache upgrade && \
    apk --no-cache add ca-certificates wget

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/app .

# Create user with fixed name and UID
RUN addgroup -S -g ${APP_GID} appgroup && \
    adduser -S -g appgroup -u ${APP_UID} ${APP_USER} && \
    # Set proper permissions
    chown -R ${APP_UID}:${APP_GID} /app && \
    chmod -R 550 /app && \
    chmod 500 /app/app

# Use the numerical UID/GID instead of username
USER ${APP_UID}:${APP_GID}

EXPOSE 5003

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget --no-verbose --tries=1 --spider http://localhost:5003/api/health || exit 1

CMD ["./app"]