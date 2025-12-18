#!/bin/bash

# Helper script to run the PlayPI Dashboard Docker container with all necessary ports

IMAGE_NAME="taqelah/playpi_dashboard:latest"

echo "Starting PlayPI Dashboard container..."
echo ""
echo "Exposed ports:"
echo "  - 8000: Web Dashboard"
echo "  - 8080: RESTful Inventory Manager"
echo "  - 8081: GraphQL Inventory Manager"
echo "  - 8082: gRPC Inventory Manager"
echo "  - 8084: gRPC User Registration"
echo "  - 8085: RESTful Task Manager"
echo "  - 8086: WebSocket Live Chat"
echo ""

docker run --rm \
  -p 8000:8000 \
  -p 8080:8080 \
  -p 8081:8081 \
  -p 8082:8082 \
  -p 8084:8084 \
  -p 8085:8085 \
  -p 8086:8086 \
  --name playpi-dashboard \
  $IMAGE_NAME

echo ""
echo "Dashboard stopped."
