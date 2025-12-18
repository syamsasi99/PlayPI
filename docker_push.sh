#!/bin/bash

set -e  # Exit if any command fails

# Define variables
IMAGE_NAME="taqelah/playpi_dashboard"
TAG="latest"

echo "Building Docker image for PlayPI Dashboard..."
docker build -t $IMAGE_NAME:$TAG .

echo "Tagging image..."
docker tag $IMAGE_NAME:$TAG $IMAGE_NAME:$TAG

echo "Pushing Docker image to Docker Hub..."
docker push $IMAGE_NAME:$TAG

echo "Docker image $IMAGE_NAME:$TAG pushed successfully!"
echo ""
echo "To run the dashboard, use:"
echo "  docker run -p 8000:8000 -p 8080:8080 -p 8081:8081 -p 8082:8082 -p 8084:8084 -p 8085:8085 -p 8086:8086 $IMAGE_NAME:$TAG"
echo ""
echo "Or use the shorter format:"
echo "  docker run -p 8000-8002:8000-8002 -p 8004-8006:8004-8006 $IMAGE_NAME:$TAG"
echo ""
echo "Then open http://localhost:8000 in your browser"
echo ""
echo "NOTE: You MUST expose all ports (8000, 8080-8086) for the dashboard to work properly."
echo "The dashboard (8000) needs to communicate with the API services (8080-8086)."