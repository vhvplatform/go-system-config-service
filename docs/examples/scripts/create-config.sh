#!/bin/bash
# Example script to create a configuration with versioning

set -e

API_URL="${API_URL:-http://localhost:8085}"
AUTH_TOKEN="${AUTH_TOKEN}"

if [ -z "$AUTH_TOKEN" ]; then
  echo "Error: AUTH_TOKEN environment variable is required"
  exit 1
fi

echo "Creating new configuration..."

# Create configuration
response=$(curl -s -X POST "$API_URL/api/v1/configs" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "key": "api.rate_limit",
    "value": {
      "requests_per_second": 100,
      "burst_size": 200,
      "enabled": true
    },
    "environment": "production",
    "tenant_id": "global",
    "description": "API rate limiting configuration",
    "tags": ["api", "rate-limit", "performance"]
  }')

echo "Response: $response"

# Extract version from response
version=$(echo "$response" | jq -r '.version')
echo "Created version: $version"

# Wait a moment
sleep 2

# Activate the configuration
echo "Activating configuration..."
curl -s -X POST "$API_URL/api/v1/configs/api.rate_limit/activate" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"version\": \"$version\"}"

echo "Configuration activated successfully!"

# View version history
echo "Fetching version history..."
curl -s -X GET "$API_URL/api/v1/configs/api.rate_limit/history" \
  -H "Authorization: Bearer $AUTH_TOKEN" | jq '.'
