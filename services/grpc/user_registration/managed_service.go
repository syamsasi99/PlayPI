package user_registration

import (
	"context"
	"net"

	"github.com/abhivaikar/playpi/pkg/servicemanager"
	pb "github.com/abhivaikar/playpi/services/grpc/user_registration/pb"
	"google.golang.org/grpc"
)

// GRPCUserService implements the ManagedService interface
type GRPCUserService struct {
	grpcServer *grpc.Server
	listener   net.Listener
}

// NewGRPCUserService creates a new GRPCUserService instance
func NewGRPCUserService() *GRPCUserService {
	return &GRPCUserService{}
}

// Start starts the gRPC User Registration service
func (s *GRPCUserService) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", ":8084")
	if err != nil {
		return err
	}
	s.listener = lis

	srv := NewServer()

	s.grpcServer = grpc.NewServer()
	pb.RegisterUserServiceServer(s.grpcServer, srv)

	// This blocks until server is stopped
	if err := s.grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

// Stop gracefully stops the service
func (s *GRPCUserService) Stop(ctx context.Context) error {
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
func (s *GRPCUserService) GetInfo() servicemanager.ServiceInfo {
	return servicemanager.ServiceInfo{
		ID:          "grpc-user-registration",
		Name:        "gRPC User Registration",
		Description: "A gRPC API for user registration, authentication, and profile management",
		Port:        8084,
		Type:        "grpc",
	}
}
