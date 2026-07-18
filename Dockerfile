# Build stage
FROM golang:1.26.5 as build

WORKDIR /go/src/app

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Run vet and tests
RUN go vet -v
RUN go test -v

# Build the application
RUN CGO_ENABLED=0 go build -o /go/bin/go-ws-proxy

# Final stage from https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static-debian13

# Copy binary from builder
COPY --from=build /go/bin/go-ws-proxy /

# Expose default port
EXPOSE 8080

# Run the application
ENTRYPOINT ["/go-ws-proxy", "-listenHostAndPort", ":8080"]
CMD ["-tcpHostAndPort", "localhost:31415"]
