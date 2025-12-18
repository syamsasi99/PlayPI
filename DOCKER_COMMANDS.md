# PlayPI Dashboard - Docker Commands Quick Reference

## Pull the Image

```bash
docker pull taqelah/playpi_dashboard:latest
```

## Run the Dashboard

### Full Command (Copy-Paste Ready)

```bash
docker run -p 8000:8000 -p 8080:8080 -p 8081:8081 -p 8082:8082 -p 8084:8084 -p 8085:8085 -p 8086:8086 taqelah/playpi_dashboard:latest
```

### Run in Detached Mode (Background)

```bash
docker run -d --name playpi-dashboard -p 8000:8000 -p 8080:8080 -p 8081:8081 -p 8082:8082 -p 8084:8084 -p 8085:8085 -p 8086:8086 taqelah/playpi_dashboard:latest
```

### Run with Auto-Remove (Container deleted on stop)

```bash
docker run --rm -p 8000:8000 -p 8080:8080 -p 8081:8081 -p 8082:8082 -p 8084:8084 -p 8085:8085 -p 8086:8086 taqelah/playpi_dashboard:latest
```

## Access the Dashboard

Open your browser to: **http://localhost:8000**

## Port Mapping Explained

| Host Port | Container Port | Service |
|-----------|---------------|---------|
| 8000 | 8000 | Web Dashboard (UI) |
| 8080 | 8080 | RESTful Inventory Manager |
| 8081 | 8081 | GraphQL Inventory Manager |
| 8082 | 8082 | gRPC Inventory Manager |
| 8084 | 8084 | gRPC User Registration |
| 8085 | 8085 | RESTful Task Manager |
| 8086 | 8086 | WebSocket Live Chat |

**⚠️ IMPORTANT:** All ports (8000-8086) MUST be mapped for the dashboard "Try It" feature to work. The browser needs to access the API services directly.

## Container Management

### Check Running Containers

```bash
docker ps | grep playpi
```

### View Container Logs

```bash
docker logs playpi-dashboard
```

### Follow Container Logs (Real-time)

```bash
docker logs -f playpi-dashboard
```

### Stop the Container

```bash
docker stop playpi-dashboard
```

### Remove the Container

```bash
docker rm playpi-dashboard
```

### Stop and Remove in One Command

```bash
docker stop playpi-dashboard && docker rm playpi-dashboard
```

## Helper Script

For convenience, use the provided helper script:

```bash
./docker_run.sh
```

This automatically runs the container with all required ports exposed.

## Building Locally

If you want to build the image yourself:

```bash
docker build -t taqelah/playpi_dashboard:latest .
```

## Pushing to Docker Hub

```bash
./docker_push.sh
```

Or manually:

```bash
docker login
docker push taqelah/playpi_dashboard:latest
```

## Troubleshooting

### "Failed to fetch" Error in Dashboard

This means the API service ports aren't accessible from your browser. Make sure you ran the container with ALL ports mapped:

```bash
# ❌ WRONG - Only dashboard port
docker run -p 8000:8000 taqelah/playpi_dashboard:latest

# ✅ CORRECT - All ports
docker run -p 8000:8000 -p 8080:8080 -p 8081:8081 -p 8082:8082 -p 8084:8084 -p 8085:8085 -p 8086:8086 taqelah/playpi_dashboard:latest
```

### Port Already in Use

If you see "port is already allocated":

```bash
# Find what's using the port (example for port 8000)
lsof -i :8000

# Kill the process or stop the conflicting container
docker stop $(docker ps -q --filter "publish=8000")
```

### Container Won't Start

Check the logs:

```bash
docker logs playpi-dashboard
```

### Remove Old Images

Clean up old/unused images:

```bash
docker image prune -a
```
