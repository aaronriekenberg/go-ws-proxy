# Build stage
FROM golang:1.26.5-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-ws-proxy .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS support
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/go-ws-proxy .

# Default environment variables
ENV LISTEN_HOST_AND_PORT=0.0.0.0:8080
ENV TCP_HOST_AND_PORT=localhost:31415

# Expose default port
EXPOSE 8080

# Run the application
ENTRYPOINT ["./go-ws-proxy", "-listenHostAndPort", "0.0.0.0:8080"]
CMD ["-tcpHostAndPort", "localhost:31415"]
