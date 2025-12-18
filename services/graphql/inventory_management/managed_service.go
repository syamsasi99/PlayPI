package graphql

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/abhivaikar/playpi/pkg/servicemanager"
	"github.com/graphql-go/graphql"
)

// GraphQLService implements the ManagedService interface
type GraphQLService struct {
	server *http.Server
}

// NewGraphQLService creates a new GraphQLService instance
func NewGraphQLService() *GraphQLService {
	return &GraphQLService{}
}

// Start starts the GraphQL Inventory service
func (s *GraphQLService) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		var params struct {
			Query string `json:"query"`
		}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		result := graphql.Do(graphql.Params{
			Schema:        Schema,
			RequestString: params.Query,
		})

		json.NewEncoder(w).Encode(result)
	})

	s.server = &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	// This blocks until server is stopped
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Stop gracefully stops the service
func (s *GraphQLService) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}

// GetInfo returns service information
func (s *GraphQLService) GetInfo() servicemanager.ServiceInfo {
	return servicemanager.ServiceInfo{
		ID:          "graphql-inventory-manager",
		Name:        "GraphQL Inventory Manager",
		Description: "A GraphQL API for inventory management with flexible queries and mutations",
		Port:        8081,
		Type:        "graphql",
	}
}
