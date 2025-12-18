package task_management

import (
	"context"
	"net/http"

	"github.com/abhivaikar/playpi/pkg/servicemanager"
)

// TaskService implements the ManagedService interface
type TaskService struct {
	server *http.Server
}

// NewTaskService creates a new TaskService instance
func NewTaskService() *TaskService {
	return &TaskService{}
}

// Start starts the RESTful Task Management service
func (s *TaskService) Start(ctx context.Context) error {
	r := setupRouter()

	s.server = &http.Server{
		Addr:    ":8085",
		Handler: r,
	}

	// This blocks until server is stopped
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop gracefully stops the service
func (s *TaskService) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}

// GetInfo returns service information
func (s *TaskService) GetInfo() servicemanager.ServiceInfo {
	return servicemanager.ServiceInfo{
		ID:          "restful-task-manager",
		Name:        "RESTful Task Manager",
		Description: "A RESTful API for task management with priority and status tracking",
		Port:        8085,
		Type:        "restful",
	}
}
