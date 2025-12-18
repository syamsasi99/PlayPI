# PlayPI Dashboard

Web-based management interface for all PlayPI API services.

## Features

- **Start/Stop Services**: Control individual services with one click
- **Real-time Status Updates**: WebSocket-based live updates of service states
- **Service Details**: View port, description, and type for each service
- **Offline/Local Only**: No external dependencies, completely localhost-based
- **Beautiful UI**: Modern, dark-themed responsive interface
- **Error Handling**: Clear error messages for port conflicts and failures

## Usage

### Start the Dashboard

```bash
./playpi start dashboard
```

### Access in Browser

Open your browser and navigate to:
```
http://localhost:8000
```

### Managing Services

1. **Start a Service**: Click the green "Start" button on any service card
2. **Stop a Service**: Click the red "Stop" button on a running service
3. **View Status**: Service status updates automatically (running/stopped)
4. **Monitor Errors**: Port conflicts and errors are displayed in the service cards

## Architecture

### Backend

- **Server**: Go HTTP server with REST API and WebSocket support
- **Port**: 8000
- **Framework**: Standard library `net/http` with Gorilla WebSocket
- **Service Management**: Unified service manager with lifecycle control

### Frontend

- **Technology**: Vanilla JavaScript (no build step required)
- **Styling**: Modern CSS with dark theme
- **Real-time Updates**: WebSocket for instant status changes
- **Embedded**: Static files bundled in binary using go:embed

## API Endpoints

### REST API

#### GET /api/services
Returns status of all registered services.

**Response:**
```json
[
  {
    "id": "restful-inventory-manager",
    "name": "RESTful Inventory Manager",
    "description": "A RESTful API for inventory management...",
    "port": 8080,
    "type": "restful",
    "status": "stopped"
  }
]
```

#### POST /api/services/start
Start a specific service.

**Request:**
```json
{
  "service_id": "restful-inventory-manager"
}
```

**Response:**
```json
{
  "status": "started"
}
```

#### POST /api/services/stop
Stop a running service.

**Request:**
```json
{
  "service_id": "restful-inventory-manager"
}
```

**Response:**
```json
{
  "status": "stopped"
}
```

### WebSocket

#### WS /ws
Real-time service status updates.

**Connection:**
```javascript
const ws = new WebSocket('ws://localhost:8000/ws');
```

**Messages:** Server pushes service status array whenever changes occur.

## Service Manager

The dashboard uses a unified service manager that provides:

### Features

- **Lifecycle Management**: Start/stop services with context-based cancellation
- **Status Tracking**: Real-time status monitoring (stopped/starting/running/stopping/error)
- **Port Conflict Detection**: Checks port availability before starting services
- **Graceful Shutdown**: 5-second timeout with fallback for gRPC services
- **Error Handling**: Captures startup errors within 500ms
- **Concurrent Safety**: Thread-safe with sync.RWMutex

### Service Types

All 6 PlayPI services are managed:

1. **RESTful Inventory Manager** (port 8080)
2. **RESTful Task Manager** (port 8085)
3. **GraphQL Inventory Manager** (port 8081)
4. **gRPC Inventory Manager** (port 8082)
5. **gRPC User Registration** (port 8084)
6. **WebSocket Live Chat** (port 8086)

## Technical Details

### Dependencies

- `github.com/gorilla/websocket` - WebSocket support (already in go.mod)
- Standard library - Everything else

### File Structure

```
services/dashboard/
├── server.go          # Main dashboard server
├── websocket.go       # WebSocket handler
└── static/
    ├── index.html     # Dashboard UI
    ├── styles.css     # Styling
    └── app.js         # Frontend logic
```

### Concurrency Model

- Service operations run in goroutines to prevent blocking
- WebSocket broadcast channel for real-time updates
- Separate mutexes for service registry and WebSocket clients
- Non-blocking service startup with error channel

### Error Scenarios

The dashboard handles:

- **Port Already in Use**: Detected before service start, error displayed in UI
- **Service Startup Failures**: Captured within 500ms, service remains stopped
- **WebSocket Disconnections**: Automatic reconnection every 3 seconds
- **Invalid Service IDs**: HTTP 500 with error message
- **Service Already Running/Stopped**: Appropriate error messages

## Security Notes

**This is a localhost-only tool**. Do NOT expose to network without:
- Authentication
- HTTPS/WSS
- Proper CORS configuration
- Rate limiting

## Development

### Building

```bash
go build -o playpi main.go
```

Static files are embedded automatically using `//go:embed`.

### Testing

1. Start dashboard: `./playpi start dashboard`
2. Open browser: `http://localhost:8000`
3. Test start/stop operations
4. Check WebSocket updates in browser dev tools

## Future Enhancements

Potential additions:
- Service logs viewer
- Metrics and monitoring
- Custom port configuration
- Batch operations (start all, stop all)
- API request playground
- Service health checks
- Export/import configurations
