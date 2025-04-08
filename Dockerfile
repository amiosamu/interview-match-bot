FROM golang:1.23-alpine AS builder

# Install git and ca-certificates for potential dependency fetching
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o interview-bot ./cmd/bot

# Create a minimal runtime image
FROM alpine:3.18

# Install ca-certificates for HTTPS and postgresql-client for db migrations
RUN apk --no-cache add ca-certificates postgresql-client

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/interview-bot .

# Set environment variable
ENV TELEGRAM_BOT_TOKEN=""
ENV DATABASE_URL=""

# Run the binary
CMD ["./interview-bot"]