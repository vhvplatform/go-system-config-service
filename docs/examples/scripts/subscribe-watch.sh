#!/bin/bash
# Example script to subscribe to configuration changes

set -e

API_URL="${API_URL:-http://localhost:8085}"
AUTH_TOKEN="${AUTH_TOKEN}"
CALLBACK_URL="${CALLBACK_URL}"

if [ -z "$AUTH_TOKEN" ]; then
  echo "Error: AUTH_TOKEN environment variable is required"
  exit 1
fi

if [ -z "$CALLBACK_URL" ]; then
  echo "Error: CALLBACK_URL environment variable is required"
  exit 1
fi

echo "Subscribing to configuration changes..."

# Subscribe to configuration changes
response=$(curl -s -X POST "$API_URL/api/v1/watch/subscribe" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"subscriber_id\": \"service-$(uuidgen)\",
    \"service_name\": \"example-service\",
    \"tenant_id\": \"tenant-xyz\",
    \"watch_patterns\": [
      \"db.*\",
      \"api.*.timeout\",
      \"feature_flags.*\"
    ],
    \"environments\": [\"production\", \"staging\"],
    \"callback_url\": \"$CALLBACK_URL\",
    \"retry_policy\": {
      \"max_attempts\": 3,
      \"backoff\": \"exponential\"
    }
  }")

echo "Response: $response"

# Extract subscription ID
subscription_id=$(echo "$response" | jq -r '.subscription_id')
echo "Subscription ID: $subscription_id"
echo "Save this ID to unsubscribe later"

# List active subscriptions
echo "Listing active subscriptions..."
curl -s -X GET "$API_URL/api/v1/watch/subscriptions" \
  -H "Authorization: Bearer $AUTH_TOKEN" | jq '.'
