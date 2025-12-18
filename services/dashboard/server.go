package dashboard

import (
	"context"
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"sync"

	"github.com/abhivaikar/playpi/pkg/servicemanager"
	"github.com/gorilla/websocket"
)

//go:embed static/*
var staticFiles embed.FS

// DashboardServer manages the web dashboard
type DashboardServer struct {
	serviceManager *servicemanager.ServiceManager
	httpServer     *http.Server
	wsClients      map[*websocket.Conn]bool
	clientsMu      sync.RWMutex
	broadcast      chan []byte
}

// NewDashboardServer creates a new DashboardServer instance
func NewDashboardServer(sm *servicemanager.ServiceManager) *DashboardServer {
	return &DashboardServer{
		serviceManager: sm,
		wsClients:      make(map[*websocket.Conn]bool),
		broadcast:      make(chan []byte, 100),
	}
}

// Start starts the dashboard HTTP server
func (d *DashboardServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	// Serve static files
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return err
	}
	mux.Handle("/", http.FileServer(http.FS(staticFS)))

	// API endpoints
	mux.HandleFunc("/api/services", d.handleGetServices)
	mux.HandleFunc("/api/services/start", d.handleStartService)
	mux.HandleFunc("/api/services/stop", d.handleStopService)
	mux.HandleFunc("/ws", d.handleWebSocket)

	d.httpServer = &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	// Start broadcast goroutine
	go d.broadcastHandler()

	log.Println("Dashboard server running on http://localhost:8000")

	// This blocks until server is stopped
	if err := d.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop gracefully stops the dashboard server
func (d *DashboardServer) Stop(ctx context.Context) error {
	if d.httpServer == nil {
		return nil
	}
	return d.httpServer.Shutdown(ctx)
}

// handleGetServices returns all service statuses
func (d *DashboardServer) handleGetServices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	services := d.serviceManager.GetAllStatus()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(services)
}

// handleStartService starts a service
func (d *DashboardServer) handleStartService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ServiceID string `json:"service_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := d.serviceManager.Start(req.ServiceID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Broadcast status update
	d.broadcastStatus()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "started"})
}

// handleStopService stops a service
func (d *DashboardServer) handleStopService(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ServiceID string `json:"service_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := d.serviceManager.Stop(req.ServiceID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Broadcast status update
	d.broadcastStatus()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "stopped"})
}

// broadcastStatus sends current service statuses to all WebSocket clients
func (d *DashboardServer) broadcastStatus() {
	services := d.serviceManager.GetAllStatus()
	data, err := json.Marshal(services)
	if err != nil {
		log.Printf("Error marshaling service status: %v", err)
		return
	}

	select {
	case d.broadcast <- data:
	default:
		log.Println("Broadcast channel full, skipping update")
	}
}

// broadcastHandler processes broadcast messages
func (d *DashboardServer) broadcastHandler() {
	for {
		message := <-d.broadcast
		d.clientsMu.RLock()
		for client := range d.wsClients {
			if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				client.Close()
				d.clientsMu.RUnlock()
				d.clientsMu.Lock()
				delete(d.wsClients, client)
				d.clientsMu.Unlock()
				d.clientsMu.RLock()
			}
		}
		d.clientsMu.RUnlock()
	}
}
