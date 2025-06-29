#!/bin/bash

# Git Tag Release Script
# Creates and pushes a Git tag via GitHub API
# Usage: ./scripts/git-tag.sh <tag_name> <repo_name> [aws_region] [secret_name]

set -e  # Exit on any error

# Parameters
TAG_NAME="${1}"
REPO_NAME="${2}"
AWS_REGION="${3:-eu-west-1}"
SECRET_NAME="${4:-ahorro-app-secrets}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Validate inputs
if [ -z "$TAG_NAME" ] || [ -z "$REPO_NAME" ]; then
    log_error "Usage: $0 <tag_name> <repo_name> [aws_region] [secret_name]"
    log_error "Example: $0 build-250629-1234 ahorro-transactions-service"
    exit 1
fi

log_info "Creating and pushing git tag: $TAG_NAME"
log_info "Repository: savak1990/$REPO_NAME"
log_info "AWS Region: $AWS_REGION"
log_info "Secret Name: $SECRET_NAME"

# Extract GitHub token from AWS Secrets Manager
log_info "Retrieving GitHub token from AWS Secrets Manager..."
GITHUB_TOKEN=$(aws secretsmanager get-secret-value \
    --secret-id "$SECRET_NAME" \
    --query 'SecretString' \
    --output text \
    --region "$AWS_REGION" 2>/dev/null | \
    jq -r '.github_token' 2>/dev/null)

if [ -z "$GITHUB_TOKEN" ] || [ "$GITHUB_TOKEN" = "null" ]; then
    log_error "Failed to retrieve GitHub token from secrets manager"
    log_error "Make sure the secret '$SECRET_NAME' exists and contains 'github_token' field"
    exit 1
fi

log_success "GitHub token retrieved successfully"

# Get the current commit SHA
log_info "Determining commit SHA..."
if [ -n "$CODEBUILD_RESOLVED_SOURCE_VERSION" ]; then
    COMMIT_SHA="$CODEBUILD_RESOLVED_SOURCE_VERSION"
    log_info "Using CodeBuild commit SHA: $COMMIT_SHA"
elif git rev-parse HEAD >/dev/null 2>&1; then
    COMMIT_SHA=$(git rev-parse HEAD)
    log_info "Using local Git commit SHA: $COMMIT_SHA"
else
    log_warning "No local Git repository found, fetching latest commit from GitHub API..."
    COMMIT_SHA=$(curl -s -H "Authorization: token $GITHUB_TOKEN" \
        "https://api.github.com/repos/savak1990/$REPO_NAME/commits/main" | \
        jq -r '.sha' 2>/dev/null)
    
    if [ -z "$COMMIT_SHA" ] || [ "$COMMIT_SHA" = "null" ]; then
        log_error "Failed to get commit SHA from GitHub API"
        exit 1
    fi
    log_info "Using GitHub API commit SHA: $COMMIT_SHA"
fi

# Step 1: Create annotated tag object
log_info "Step 1: Creating annotated tag object..."
TAG_RESPONSE=$(curl -s -X POST \
    -H "Authorization: token $GITHUB_TOKEN" \
    -H "Accept: application/vnd.github.v3+json" \
    -d "{\"tag\":\"$TAG_NAME\",\"message\":\"Release $TAG_NAME - Automated build and deployment\",\"object\":\"$COMMIT_SHA\",\"type\":\"commit\"}" \
    "https://api.github.com/repos/savak1990/$REPO_NAME/git/tags")

# Check if tag object was created successfully
if echo "$TAG_RESPONSE" | jq -e '.sha' >/dev/null 2>&1; then
    TAG_SHA=$(echo "$TAG_RESPONSE" | jq -r '.sha')
    log_success "Tag object created successfully (SHA: $TAG_SHA)"
else
    log_error "Failed to create tag object:"
    echo "$TAG_RESPONSE" | jq '.' 2>/dev/null || echo "$TAG_RESPONSE"
    exit 1
fi

# Step 2: Create tag reference
log_info "Step 2: Creating tag reference..."
REF_RESPONSE=$(curl -s -X POST \
    -H "Authorization: token $GITHUB_TOKEN" \
    -H "Accept: application/vnd.github.v3+json" \
    -d "{\"ref\":\"refs/tags/$TAG_NAME\",\"sha\":\"$TAG_SHA\"}" \
    "https://api.github.com/repos/savak1990/$REPO_NAME/git/refs")

# Check if tag reference was created successfully
if echo "$REF_RESPONSE" | jq -e '.ref' >/dev/null 2>&1; then
    log_success "Tag reference created successfully"
    log_success "🎉 Git tag '$TAG_NAME' created and pushed successfully via GitHub API!"
    log_info "View the tag at: https://github.com/savak1990/$REPO_NAME/releases/tag/$TAG_NAME"
else
    log_error "Failed to create tag reference:"
    echo "$REF_RESPONSE" | jq '.' 2>/dev/null || echo "$REF_RESPONSE"
    exit 1
fi
