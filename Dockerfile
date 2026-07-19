# Build stage
FROM golang:1.26.5 as build

# Copy source code
WORKDIR /go/src/app
COPY . .

# Download dependencies
RUN go mod download

# Run vet and tests
RUN go vet -v
RUN go test -v

# Build the application with release tag embedded via ldflags
ARG RELEASE_TAG=dev
RUN CGO_ENABLED=0 go build -ldflags="-X main.releaseTag=${RELEASE_TAG}" -o /go/bin/go-ws-proxy

# Final stage from https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static-debian13

# Copy binary from builder
COPY --from=build /go/bin/go-ws-proxy /

# Expose default port
EXPOSE 80

# Run the application
ENTRYPOINT ["/go-ws-proxy", "-listenHostAndPort", ":80"]
CMD ["-tcpHostAndPort", "localhost:31415"]
