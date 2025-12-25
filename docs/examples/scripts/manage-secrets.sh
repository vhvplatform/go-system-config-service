#!/bin/bash
# Example script to create and rotate secrets

set -e

API_URL="${API_URL:-http://localhost:8085}"
AUTH_TOKEN="${AUTH_TOKEN}"

if [ -z "$AUTH_TOKEN" ]; then
  echo "Error: AUTH_TOKEN environment variable is required"
  exit 1
fi

# Create a secret
echo "Creating secret..."
curl -s -X POST "$API_URL/api/v1/secrets" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "key": "db_password",
    "value": "super-secret-password-123",
    "environment": "production",
    "tenant_id": "tenant-xyz",
    "metadata": {
      "description": "Production database password",
      "rotation_days": 90,
      "created_by": "admin"
    }
  }' | jq '.'

echo "Secret created successfully!"

# Retrieve secret (value will be masked in logs)
echo "Retrieving secret..."
secret=$(curl -s -X GET "$API_URL/api/v1/secrets/db_password" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "X-Environment: production" \
  -H "X-Tenant-ID: tenant-xyz")

echo "Secret retrieved (value not shown for security)"

# Rotate secret
echo "Rotating secret..."
curl -s -X POST "$API_URL/api/v1/secrets/db_password/rotate" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "environment": "production",
    "tenant_id": "tenant-xyz"
  }' | jq '.'

echo "Secret rotated successfully!"

# View audit log
echo "Fetching secret audit log..."
curl -s -X GET "$API_URL/api/v1/secrets/db_password/audit" \
  -H "Authorization: Bearer $AUTH_TOKEN" \
  -H "X-Environment: production" \
  -H "X-Tenant-ID: tenant-xyz" | jq '.'
