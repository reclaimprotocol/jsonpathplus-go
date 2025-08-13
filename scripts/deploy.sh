#!/bin/bash

# Deployment script for JSONPath Plus Go to GCR and Cloud Run
set -e

PROJECT_ID=${PROJECT_ID:-"your-gcp-project-id"}
IMAGE_NAME="jsonpathplus-go"
VERSION=${VERSION:-"latest"}
REGISTRY="gcr.io"
SERVICE_NAME=${SERVICE_NAME:-"jsonpathplus-api"}
REGION=${REGION:-"us-central1"}

GCR_IMAGE="${REGISTRY}/${PROJECT_ID}/${IMAGE_NAME}:${VERSION}"

echo "🚀 Deploying JSONPath Plus Go to Google Cloud..."

# Check if the image exists
if ! docker image inspect ${GCR_IMAGE} &> /dev/null; then
    echo "❌ Image ${GCR_IMAGE} not found locally. Run build.sh first."
    exit 1
fi

# Push to GCR
echo "📤 Pushing to Google Container Registry..."
docker push ${GCR_IMAGE}

# Deploy to Cloud Run
echo "☁️  Deploying to Cloud Run..."
gcloud run deploy ${SERVICE_NAME} \
    --image ${GCR_IMAGE} \
    --platform managed \
    --region ${REGION} \
    --port 8080 \
    --memory 512Mi \
    --cpu 1 \
    --min-instances 0 \
    --max-instances 10 \
    --allow-unauthenticated \
    --timeout 300 \
    --concurrency 80 \
    --set-env-vars "GIN_MODE=release,LOG_LEVEL=info" \
    --project ${PROJECT_ID}

# Get the service URL
SERVICE_URL=$(gcloud run services describe ${SERVICE_NAME} --region ${REGION} --format="value(status.url)" --project ${PROJECT_ID})

echo ""
echo "✅ Deployment completed successfully!"
echo "🌐 Service URL: ${SERVICE_URL}"
echo "📊 Health check: ${SERVICE_URL}/health"
echo "🔍 Metrics: ${SERVICE_URL}/metrics"
echo ""
echo "Test the service:"
echo "curl '${SERVICE_URL}/query?path=\$.users[*].name'"