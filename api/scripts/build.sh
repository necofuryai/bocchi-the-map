#!/bin/bash

# Docker build and push script for Bocchi The Map API
# Usage: ./scripts/build.sh [environment] [project-id] [region]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
ENVIRONMENT=${1:-dev}
PROJECT_ID=${2:-}
REGION=${3:-asia-northeast1}
IMAGE_NAME="bocchi-api"
SERVICE_NAME="bocchi-api-${ENVIRONMENT}"

# Print usage
usage() {
    echo "Usage: $0 [environment] [project-id] [region]"
    echo ""
    echo "Arguments:"
    echo "  environment  Environment (dev, staging, prod) - default: dev"
    echo "  project-id   GCP Project ID - required"
    echo "  region       GCP Region - default: asia-northeast1"
    echo ""
    echo "Example:"
    echo "  $0 dev my-gcp-project asia-northeast1"
    exit 1
}

# Validate inputs
if [ -z "$PROJECT_ID" ]; then
    echo -e "${RED}Error: Project ID is required${NC}"
    usage
fi

# Validate environment
if [[ ! "$ENVIRONMENT" =~ ^(dev|staging|prod)$ ]]; then
    echo -e "${RED}Error: Environment must be dev, staging, or prod${NC}"
    usage
fi

# Set image tag
IMAGE_TAG="gcr.io/${PROJECT_ID}/${IMAGE_NAME}:${ENVIRONMENT}-$(date +%Y%m%d%H%M%S)"
LATEST_TAG="gcr.io/${PROJECT_ID}/${IMAGE_NAME}:${ENVIRONMENT}-latest"

echo -e "${GREEN}Building Docker image for Bocchi The Map API${NC}"
echo -e "${YELLOW}Environment:${NC} ${ENVIRONMENT}"
echo -e "${YELLOW}Project ID:${NC} ${PROJECT_ID}"
echo -e "${YELLOW}Region:${NC} ${REGION}"
echo -e "${YELLOW}Image Tag:${NC} ${IMAGE_TAG}"
echo ""

# Check if gcloud is installed and authenticated
if ! command -v gcloud &> /dev/null; then
    echo -e "${RED}Error: gcloud CLI is not installed${NC}"
    exit 1
fi

# Check if Docker is running
if ! docker info &> /dev/null; then
    echo -e "${RED}Error: Docker is not running${NC}"
    exit 1
fi

# Configure Docker for GCR
echo -e "${GREEN}Configuring Docker for Google Container Registry...${NC}"
gcloud auth configure-docker --quiet

# Build the Docker image
echo -e "${GREEN}Building Docker image...${NC}"
docker build -t "${IMAGE_TAG}" -t "${LATEST_TAG}" .

if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Docker build failed${NC}"
    exit 1
fi

echo -e "${GREEN}Docker image built successfully!${NC}"

# Push the image to GCR
echo -e "${GREEN}Pushing image to Google Container Registry...${NC}"
docker push "${IMAGE_TAG}"
docker push "${LATEST_TAG}"

if [ $? -ne 0 ]; then
    echo -e "${RED}Error: Docker push failed${NC}"
    exit 1
fi

echo -e "${GREEN}Image pushed successfully!${NC}"

# Deploy to Cloud Run (optional)
read -p "Do you want to deploy to Cloud Run? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${GREEN}Deploying to Cloud Run...${NC}"
    
    gcloud run deploy "${SERVICE_NAME}" \
        --image="${LATEST_TAG}" \
        --platform=managed \
        --region="${REGION}" \
        --project="${PROJECT_ID}" \
        --allow-unauthenticated \
        --port=8080 \
        --memory=1Gi \
        --cpu=1 \
        --max-instances=10 \
        --min-instances=$([ "$ENVIRONMENT" = "prod" ] && echo 1 || echo 0) \
        --quiet

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}Deployment successful!${NC}"
        
        # Get the service URL
        SERVICE_URL=$(gcloud run services describe "${SERVICE_NAME}" \
            --platform=managed \
            --region="${REGION}" \
            --project="${PROJECT_ID}" \
            --format="value(status.url)")
        
        echo -e "${GREEN}Service URL:${NC} ${SERVICE_URL}"
        echo -e "${GREEN}Health check:${NC} ${SERVICE_URL}/health"
    else
        echo -e "${RED}Deployment failed!${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}Skipping deployment.${NC}"
    echo -e "${YELLOW}To deploy manually, run:${NC}"
    echo "gcloud run deploy ${SERVICE_NAME} --image=${LATEST_TAG} --platform=managed --region=${REGION} --project=${PROJECT_ID}"
fi

echo ""
echo -e "${GREEN}Build and push completed successfully!${NC}"
echo -e "${YELLOW}Image tags:${NC}"
echo "  ${IMAGE_TAG}"
echo "  ${LATEST_TAG}"