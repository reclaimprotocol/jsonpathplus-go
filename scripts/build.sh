#!/bin/bash

# Build script for JSONPath Plus Go
set -e

PROJECT_ID=${PROJECT_ID:-"your-gcp-project-id"}
IMAGE_NAME="jsonpathplus-go"
VERSION=${VERSION:-"latest"}
REGISTRY="gcr.io"

echo "🚀 Building JSONPath Plus Go for deployment..."

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    echo "❌ gcloud CLI is not installed. Please install it first."
    exit 1
fi

# Check if docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install it first."
    exit 1
fi

# Authenticate with GCR
echo "🔐 Configuring Docker for GCR..."
gcloud auth configure-docker

# Build the Docker image
echo "🏗️  Building Docker image..."
docker build -t ${IMAGE_NAME}:${VERSION} .

# Tag for GCR
GCR_IMAGE="${REGISTRY}/${PROJECT_ID}/${IMAGE_NAME}:${VERSION}"
echo "🏷️  Tagging image as ${GCR_IMAGE}"
docker tag ${IMAGE_NAME}:${VERSION} ${GCR_IMAGE}

echo "✅ Build completed successfully!"
echo "📦 Image: ${GCR_IMAGE}"
echo ""
echo "Next steps:"
echo "1. Test locally: docker run -p 8080:8080 ${IMAGE_NAME}:${VERSION}"
echo "2. Push to GCR: docker push ${GCR_IMAGE}"
echo "3. Deploy: gcloud run deploy --image ${GCR_IMAGE}"