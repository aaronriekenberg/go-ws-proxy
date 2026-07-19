# go-ws-proxy

A quick binary WebSocket to TCP proxy.

## Usage

### Binary

```
go-ws-proxy -listenHostAndPort localhost:8080 -tcpHostAndPort localhost:31415
```

### Docker

Pull the image from Docker Hub:

```bash
docker pull aaronriekenberg/go-ws-proxy:latest
```

Run a container:

```bash
docker run -d \
  -p 8080:8080 \
  --name ws-proxy \
  aaronriekenberg/go-ws-proxy:latest \
  -tcpHostAndPort your-target-host:31415
```

#### Environment and Custom Parameters

You can override the listening address and TCP target by passing command-line arguments:

```bash
docker run -d \
  -p 9090:9090 \
  --name ws-proxy \
  aaronriekenberg/go-ws-proxy:latest \
  -tcpHostAndPort remote-server.example.com:5000
```

#### Available Tags

Images are published on each release:
- `latest` - Latest release
- `vX.Y.Z` - Specific version (e.g., `v1.0.0`)
- `X.Y` - Minor version (e.g., `1.0`)

#### Docker Compose

Create a `docker-compose.yml` file:

```yaml
version: '3.8'

services:
  ws-proxy:
    image: aaronriekenberg/go-ws-proxy:latest
    container_name: ws-proxy
    ports:
      - "80:80"
    command:
      - -tcpHostAndPort
      - "localhost:31415"
    restart: unless-stopped
```

With other services in the same network:

```yaml
version: '3.8'

services:
  ws-proxy:
    image: aaronriekenberg/go-ws-proxy:latest
    container_name: ws-proxy
    ports:
      - "80:80"
    command:
      - -tcpHostAndPort
      - "backend-service:5000"
    depends_on:
      - backend-service
    restart: unless-stopped
    networks:
      - app-network

  backend-service:
    image: your-backend-image:latest
    container_name: backend-service
    networks:
      - app-network

networks:
  app-network:
    driver: bridge
```

Run with docker-compose:

```bash
docker-compose up -d
```

#### Connecting to Services on the Host

If you need to connect to a service running on the host machine (not in Docker), use `host.docker.internal` (Docker Desktop) or `host-gateway` (Docker on Linux):

```bash
# Docker Desktop (macOS, Windows)
docker run -d \
  -p 8080:8080 \
  aaronriekenberg/go-ws-proxy:latest \
  -tcpHostAndPort host.docker.internal:31415

# Docker on Linux
docker run -d \
  -p 8080:8080 \
  --add-host=host-gateway:host-gateway \
  aaronriekenberg/go-ws-proxy:latest \
  -tcpHostAndPort host-gateway:31415
```
