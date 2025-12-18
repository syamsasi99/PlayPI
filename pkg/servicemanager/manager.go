package servicemanager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// ServiceStatus represents the current state of a service
type ServiceStatus string

const (
	StatusStopped  ServiceStatus = "stopped"
	StatusStarting ServiceStatus = "starting"
	StatusRunning  ServiceStatus = "running"
	StatusStopping ServiceStatus = "stopping"
	StatusError    ServiceStatus = "error"
)

// ServiceInfo contains metadata about a service
type ServiceInfo struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Port        int           `json:"port"`
	Type        string        `json:"type"` // "restful", "graphql", "grpc", "websocket"
	Status      ServiceStatus `json:"status"`
	Error       string        `json:"error,omitempty"`
}

// ManagedService interface that all services must implement
type ManagedService interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	GetInfo() ServiceInfo
}

// ServiceManager manages all API services
type ServiceManager struct {
	services map[string]ManagedService
	mu       sync.RWMutex
	contexts map[string]context.CancelFunc
}

// NewServiceManager creates a new ServiceManager instance
func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		services: make(map[string]ManagedService),
		contexts: make(map[string]context.CancelFunc),
	}
}

// Register adds a new service to the manager
func (sm *ServiceManager) Register(service ManagedService) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	info := service.GetInfo()
	sm.services[info.ID] = service
}

// Start starts a service by ID
func (sm *ServiceManager) Start(serviceID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	service, exists := sm.services[serviceID]
	if !exists {
		return errors.New("service not found")
	}

	// Check if already running
	if _, running := sm.contexts[serviceID]; running {
		return errors.New("service already running")
	}

	info := service.GetInfo()

	// Check port availability
	if !IsPortAvailable(info.Port) {
		return fmt.Errorf("port %d is already in use", info.Port)
	}

	ctx, cancel := context.WithCancel(context.Background())
	sm.contexts[serviceID] = cancel

	// Start in goroutine with error handling
	errChan := make(chan error, 1)
	go func() {
		if err := service.Start(ctx); err != nil && err != http.ErrServerClosed {
			errChan <- err
			sm.mu.Lock()
			delete(sm.contexts, serviceID)
			sm.mu.Unlock()
			log.Printf("Service %s error: %v", serviceID, err)
		}
	}()

	// Wait briefly to catch immediate startup errors
	select {
	case err := <-errChan:
		cancel()
		return fmt.Errorf("service failed to start: %w", err)
	case <-time.After(500 * time.Millisecond):
		// Service started successfully
		log.Printf("Service %s started on port %d", serviceID, info.Port)
		return nil
	}
}

// Stop stops a service by ID
func (sm *ServiceManager) Stop(serviceID string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	cancel, exists := sm.contexts[serviceID]
	if !exists {
		return errors.New("service not running")
	}

	// Cancel context
	cancel()
	delete(sm.contexts, serviceID)

	service := sm.services[serviceID]
	ctx, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()

	log.Printf("Stopping service %s...", serviceID)
	return service.Stop(ctx)
}

// GetStatus returns the status of a specific service
func (sm *ServiceManager) GetStatus(serviceID string) (ServiceInfo, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	service, exists := sm.services[serviceID]
	if !exists {
		return ServiceInfo{}, errors.New("service not found")
	}

	info := service.GetInfo()
	_, running := sm.contexts[serviceID]
	if running {
		info.Status = StatusRunning
	} else {
		info.Status = StatusStopped
	}

	return info, nil
}

// GetAllStatus returns the status of all registered services
func (sm *ServiceManager) GetAllStatus() []ServiceInfo {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	var statuses []ServiceInfo
	for id := range sm.services {
		sm.mu.RUnlock() // Temporarily unlock to call GetStatus
		status, _ := sm.GetStatus(id)
		sm.mu.RLock() // Re-lock
		statuses = append(statuses, status)
	}
	return statuses
}

// StopAll stops all running services
func (sm *ServiceManager) StopAll() error {
	sm.mu.Lock()
	runningServices := make([]string, 0, len(sm.contexts))
	for id := range sm.contexts {
		runningServices = append(runningServices, id)
	}
	sm.mu.Unlock()

	for _, id := range runningServices {
		if err := sm.Stop(id); err != nil {
			log.Printf("Error stopping service %s: %v", id, err)
		}
	}

	return nil
}
