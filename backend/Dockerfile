# Build Stage
FROM golang:1.24.0-alpine AS builder

# Build Dependencies (gcc and musl-dev needed for SQLite CGO)
RUN apk add --no-cache git gcc musl-dev sqlite-dev

# Directory
WORKDIR /app

# Copy the Go Modules and Sum Files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . . 

# Build
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o url-monitor .

# Runtime
FROM alpine:latest

# Install runtime dependencies (ca-certificates for HTTPS, sqlite-libs for database)
RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /app/url-monitor .

# Directory for SQLITE DB
RUN mkdir -p /app/db

# Expose the HTTP port
EXPOSE 8080

# Run
CMD ["./url-monitor"]