#!/usr/bin/env bash
# scripts/github-meta.sh
# Usage: github-meta.sh <owner> <repo> <branch> <github_token>
# Prints GIT_BRANCH, GIT_COMMIT, GIT_SHORT as shell exports

set -euo pipefail
OWNER="$1"
REPO="$2"
BRANCH="$3"
GITHUB_TOKEN="$4"

API_URL="https://api.github.com/repos/$OWNER/$REPO/commits/$BRANCH"

COMMIT_JSON=$(curl -s -H "Authorization: token $GITHUB_TOKEN" "$API_URL")
GIT_COMMIT=$(echo "$COMMIT_JSON" | jq -r .sha)
GIT_SHORT=$(echo "$GIT_COMMIT" | cut -c1-7)
GIT_BRANCH="$BRANCH"

# Output as shell exports
echo "export GIT_BRANCH='$GIT_BRANCH'"
echo "export GIT_COMMIT='$GIT_COMMIT'"
echo "export GIT_SHORT='$GIT_SHORT'"
