#!/usr/bin/env bash
set -euo pipefail

# Read current version from version file
CURRENT_VERSION=$(cat version)
echo "Current version: ${CURRENT_VERSION}"

# Increment version by 0.0.1
IFS='.' read -r major minor patch <<< "$CURRENT_VERSION"
NEW_PATCH=$((patch + 1))
NEW_VERSION="${major}.${minor}.${NEW_PATCH}"

# Set TAG variable
TAG="${NEW_VERSION}"

IMAGE_NAME="pmz-budget-api"
AWS_REGION="eu-west-1"
ECR_REG="116411769996.dkr.ecr.eu-west-1.amazonaws.com/pmz-budget-api"
AWS_PROFILE="admin"

# Raspberry Pi 64-bit OS → linux/arm64
# Raspberry Pi 32-bit OS → linux/arm/v7
PLATFORM="linux/arm/v7"



echo "==> Logging in to ECR: ${ECR_REG}"
aws ecr get-login-password --region "$AWS_REGION" --profile "$AWS_PROFILE" \
    | docker login --username AWS --password-stdin "$ECR_REG"


echo ">> Building ${IMAGE_NAME}:${TAG} for ${PLATFORM} ..."
docker buildx build \
  --platform "${PLATFORM}" \
  -t "${ECR_REG}:${TAG}" \
  -f "docker/Dockerfile" \
  --push \
  .

echo "✅ Image built: ${IMAGE_NAME}:${TAG}"
echo "Run on Pi with: docker run --rm ${IMAGE_NAME}:${TAG}"
# Write new version back to file
echo "${NEW_VERSION}" > version
echo "Updated version: ${NEW_VERSION}"