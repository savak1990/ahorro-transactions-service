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

# Validate inputs
if [ -z "$TAG_NAME" ] || [ -z "$REPO_NAME" ]; then
    echo "Error: Usage: $0 <tag_name> <repo_name> [aws_region] [secret_name]"
    echo "Error: Example: $0 build-250629-1234 ahorro-transactions-service"
    exit 1
fi

echo "Creating git tag '$TAG_NAME' for repository 'savak1990/$REPO_NAME'..."

# Extract GitHub token from AWS Secrets Manager
GITHUB_TOKEN=$(aws secretsmanager get-secret-value \
    --secret-id "$SECRET_NAME" \
    --query 'SecretString' \
    --output text \
    --region "$AWS_REGION" 2>/dev/null | \
    jq -r '.github_token' 2>/dev/null)

if [ -z "$GITHUB_TOKEN" ] || [ "$GITHUB_TOKEN" = "null" ]; then
    echo "Error: Failed to retrieve GitHub token from secrets manager '$SECRET_NAME'"
    exit 1
fi

# Get the current commit SHA
if [ -n "$CODEBUILD_RESOLVED_SOURCE_VERSION" ]; then
    COMMIT_SHA="$CODEBUILD_RESOLVED_SOURCE_VERSION"
elif git rev-parse HEAD >/dev/null 2>&1; then
    COMMIT_SHA=$(git rev-parse HEAD)
else
    COMMIT_SHA=$(curl -s -H "Authorization: token $GITHUB_TOKEN" \
        "https://api.github.com/repos/savak1990/$REPO_NAME/commits/main" | \
        jq -r '.sha' 2>/dev/null)
    
    if [ -z "$COMMIT_SHA" ] || [ "$COMMIT_SHA" = "null" ]; then
        echo "Error: Failed to get commit SHA"
        exit 1
    fi
fi

# Step 1: Create annotated tag object
TAG_RESPONSE=$(curl -s -X POST \
    -H "Authorization: token $GITHUB_TOKEN" \
    -H "Accept: application/vnd.github.v3+json" \
    -d "{\"tag\":\"$TAG_NAME\",\"message\":\"Release $TAG_NAME - Automated build and deployment\",\"object\":\"$COMMIT_SHA\",\"type\":\"commit\"}" \
    "https://api.github.com/repos/savak1990/$REPO_NAME/git/tags")

# Check if tag object was created successfully
if echo "$TAG_RESPONSE" | jq -e '.sha' >/dev/null 2>&1; then
    TAG_SHA=$(echo "$TAG_RESPONSE" | jq -r '.sha')
else
    echo "Error: Failed to create tag object:"
    echo "$TAG_RESPONSE" | jq '.' 2>/dev/null || echo "$TAG_RESPONSE"
    exit 1
fi

# Step 2: Create tag reference
REF_RESPONSE=$(curl -s -X POST \
    -H "Authorization: token $GITHUB_TOKEN" \
    -H "Accept: application/vnd.github.v3+json" \
    -d "{\"ref\":\"refs/tags/$TAG_NAME\",\"sha\":\"$TAG_SHA\"}" \
    "https://api.github.com/repos/savak1990/$REPO_NAME/git/refs")

# Check if tag reference was created successfully
if echo "$REF_RESPONSE" | jq -e '.ref' >/dev/null 2>&1; then
    echo "Git tag '$TAG_NAME' created successfully"
else
    echo "Error: Failed to create tag reference:"
    echo "$REF_RESPONSE" | jq '.' 2>/dev/null || echo "$REF_RESPONSE"
    exit 1
fi
