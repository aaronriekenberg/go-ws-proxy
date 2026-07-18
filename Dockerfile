# Build stage
FROM golang:1.26.5 as build

WORKDIR /go/src/app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Run vet and tests
RUN go vet -v
RUN go test -v

# Build the application
RUN CGO_ENABLED=0 go build -o /go/bin/go-ws-proxy .

# Final stage based on https://github.com/GoogleContainerTools/distroless/blob/main/examples/go/Dockerfile
FROM gcr.io/distroless/static-debian12

# Copy binary from builder
COPY --from=build /go/bin/go-ws-proxy /go-ws-proxy

# Expose default port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/go-ws-proxy", "-listenHostAndPort", "0.0.0.0:8080"]
CMD ["-tcpHostAndPort", "localhost:31415"]
