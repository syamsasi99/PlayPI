package grpc

import (
	"context"
	"net"

	"github.com/abhivaikar/playpi/pkg/servicemanager"
	pb "github.com/abhivaikar/playpi/services/grpc/inventory_management/pb"
	"google.golang.org/grpc"
)

// GRPCInventoryService implements the ManagedService interface
type GRPCInventoryService struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

// NewGRPCInventoryService creates a new GRPCInventoryService instance
func NewGRPCInventoryService() *GRPCInventoryService {
	return &GRPCInventoryService{}
}

// Start starts the gRPC Inventory service
func (s *GRPCInventoryService) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		return err
	}
	s.listener = lis

	srv := &server{}
	srv.loadMockData()

	s.grpcServer = grpc.NewServer()
	pb.RegisterInventoryServiceServer(s.grpcServer, srv)

	// This blocks until server is stopped
	if err := s.grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

// Stop gracefully stops the service
func (s *GRPCInventoryService) Stop(ctx context.Context) error {
	if s.grpcServer == nil {
		return nil
	}

	// GracefulStop waits for ongoing RPCs to complete
	stopped := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.grpcServer.Stop() // Force stop if timeout
		return ctx.Err()
	case <-stopped:
		return nil
	}
}

// GetInfo returns service information
func (s *GRPCInventoryService) GetInfo() servicemanager.ServiceInfo {
	return servicemanager.ServiceInfo{
		ID:          "grpc-inventory-manager",
		Name:        "gRPC Inventory Manager",
		Description: "A gRPC API for inventory management with protocol buffer definitions",
		Port:        8082,
		Type:        "grpc",
	}
}
