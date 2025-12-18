package dashboard

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for localhost
	},
}

// handleWebSocket upgrades HTTP connection to WebSocket and manages client
func (d *DashboardServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// Register client
	d.clientsMu.Lock()
	d.wsClients[conn] = true
	d.clientsMu.Unlock()

	log.Println("New WebSocket client connected")

	// Send initial status
	services := d.serviceManager.GetAllStatus()
	if err := conn.WriteJSON(services); err != nil {
		log.Printf("Error sending initial status: %v", err)
		d.clientsMu.Lock()
		delete(d.wsClients, conn)
		d.clientsMu.Unlock()
		return
	}

	// Keep connection alive and listen for messages
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			d.clientsMu.Lock()
			delete(d.wsClients, conn)
			d.clientsMu.Unlock()
			break
		}

		// Echo back or handle ping/pong
		if messageType == websocket.PingMessage {
			if err := conn.WriteMessage(websocket.PongMessage, message); err != nil {
				log.Printf("WebSocket pong error: %v", err)
				break
			}
		}

		// Handle any client messages (optional - currently just keeps connection alive)
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err == nil {
			if msgType, ok := msg["type"].(string); ok && msgType == "ping" {
				// Respond to ping with current status
				services := d.serviceManager.GetAllStatus()
				if err := conn.WriteJSON(services); err != nil {
					log.Printf("Error sending ping response: %v", err)
					break
				}
			}
		}
	}

	log.Println("WebSocket client disconnected")
}
