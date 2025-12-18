package live_chat

import (
	"context"

	"github.com/abhivaikar/playpi/pkg/servicemanager"
)

// WebSocketChatService implements the ManagedService interface
type WebSocketChatService struct {
	server *WebSocketServer
}

// NewWebSocketChatService creates a new WebSocketChatService instance
func NewWebSocketChatService() *WebSocketChatService {
	return &WebSocketChatService{
		server: NewWebSocketServer(),
	}
}

// Start starts the WebSocket Chat service
func (s *WebSocketChatService) Start(ctx context.Context) error {
	// The existing StartServer method already handles the blocking call
	s.server.StartServer()
	return nil
}

// Stop gracefully stops the service
func (s *WebSocketChatService) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	// Use the existing StopServer method
	s.server.StopServer()
	return nil
}

// GetInfo returns service information
func (s *WebSocketChatService) GetInfo() servicemanager.ServiceInfo {
	return servicemanager.ServiceInfo{
		ID:          "websocket-live-chat",
		Name:        "WebSocket Live Chat",
		Description: "A WebSocket API for real-time multi-user chat with public and private messaging",
		Port:        8086,
		Type:        "websocket",
	}
}
