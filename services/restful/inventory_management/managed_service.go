package restful

import (
	"context"
	"net/http"

	"github.com/abhivaikar/playpi/pkg/servicemanager"
)

// InventoryService implements the ManagedService interface
type InventoryService struct {
	server *http.Server
}

// NewInventoryService creates a new InventoryService instance
func NewInventoryService() *InventoryService {
	return &InventoryService{}
}

// Start starts the RESTful Inventory service
func (s *InventoryService) Start(ctx context.Context) error {
	inventory = GetMockInventory()
	nextID = len(inventory) + 1
	r := setupRouter()

	s.server = &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// This blocks until server is stopped
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop gracefully stops the service
func (s *InventoryService) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}

// GetInfo returns service information
func (s *InventoryService) GetInfo() servicemanager.ServiceInfo {
	return servicemanager.ServiceInfo{
		ID:          "restful-inventory-manager",
		Name:        "RESTful Inventory Manager",
		Description: "A RESTful API for inventory management with full CRUD operations",
		Port:        8080,
		Type:        "restful",
	}
}
