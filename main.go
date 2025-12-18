/*
Copyright © 2025 Abhijeet Vaikar
*/
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/abhivaikar/playpi/pkg/servicemanager"
	"github.com/abhivaikar/playpi/services/dashboard"
	graphqlInventory "github.com/abhivaikar/playpi/services/graphql/inventory_management"
	grpcInventory "github.com/abhivaikar/playpi/services/grpc/inventory_management"
	grpcUserRegistration "github.com/abhivaikar/playpi/services/grpc/user_registration"
	restfulInventory "github.com/abhivaikar/playpi/services/restful/inventory_management"
	restfulTaskManagement "github.com/abhivaikar/playpi/services/restful/task_management"
	websocketLiveChat "github.com/abhivaikar/playpi/services/websocket/live_chat"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "playpi",
		Short: "PlayPI - A simple API Playground to practice API testing",
		Long:  `PlayPI is your local playground with restful, grpc, graphql and websocket APIs to learn and practice API testing!`,
	}

	// startCmd represents the start command
	var startCmd = &cobra.Command{
		Use:   "start [api-type]",
		Short: "Start a specific PlayPI API playground",
		Long: `Start a specific PlayPI API playground by specifying the API type:
Available options:
- dashboard (Web-based management for all services)
- restful-inventory-manager
- graphql-inventory-manager
- grpc-inventory-manager
- restful-task-manager
- grpc-user-registration
- websocket-live-chat`,
		Args: cobra.ExactArgs(1), // Ensure exactly one argument is provided
		Run: func(cmd *cobra.Command, args []string) {
			apiType := args[0]

			switch apiType {
			case "dashboard":
				fmt.Println("Starting PlayPI Dashboard...")
				startDashboard()
			case "restful-inventory-manager":
				fmt.Println("Starting RESTful API Playground for an inventory management system...")
				restfulInventory.StartServer()
			case "graphql-inventory-manager":
				fmt.Println("Starting GraphQL API Playground for an inventory management system...")
				graphqlInventory.StartServer()
			case "grpc-inventory-manager":
				fmt.Println("Starting gRPC API Playground for an inventory management system...")
				grpcInventory.StartServer()
			case "restful-task-manager":
				fmt.Println("Starting RESTful API Playground for a task management system...")
				restfulTaskManagement.StartServer()
			case "grpc-user-registration":
				fmt.Println("Starting gRPC API Playground for user registration and sign in...")
				grpcUserRegistration.StartServer()
			case "websocket-live-chat":
				fmt.Println("Starting WebSocket Playground for live chat...")
				wsLiveChatServer := websocketLiveChat.NewWebSocketServer()
				wsLiveChatServer.StartServer()
			default:
				fmt.Printf("Invalid API type: %s\n", apiType)
				fmt.Println("Available options: dashboard, restful-inventory-manager, graphql-inventory-manager, grpc-inventory-manager, restful-task-manager, grpc-user-registration, websocket-live-chat")
			}
		},
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(startCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// startDashboard initializes and starts the web dashboard
func startDashboard() {
	// Initialize service manager
	sm := servicemanager.NewServiceManager()

	// Register all services
	sm.Register(restfulInventory.NewInventoryService())
	sm.Register(restfulTaskManagement.NewTaskService())
	sm.Register(graphqlInventory.NewGraphQLService())
	sm.Register(grpcInventory.NewGRPCInventoryService())
	sm.Register(grpcUserRegistration.NewGRPCUserService())
	sm.Register(websocketLiveChat.NewWebSocketChatService())

	// Start dashboard
	dashboardServer := dashboard.NewDashboardServer(sm)
	fmt.Println("Dashboard running on http://localhost:8000")
	fmt.Println("Open your browser to manage all services")

	if err := dashboardServer.Start(context.Background()); err != nil {
		log.Fatal(err)
	}
}
