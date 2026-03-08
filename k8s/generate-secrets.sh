#!/bin/bash

# InSavein Kubernetes Secrets Generator
# This script generates secure random secrets for the InSavein platform

set -e

echo "=================================================="
echo "InSavein Kubernetes Secrets Generator"
echo "=================================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to generate a random base64 string
generate_secret() {
    local bytes=$1
    openssl rand -base64 "$bytes" | tr -d '\n'
}

# Function to generate a random password
generate_password() {
    openssl rand -base64 32 | tr -d '\n'
}

echo "Generating secure random secrets..."
echo ""

# Generate all secrets
JWT_SECRET=$(generate_secret 64)
DATA_ENCRYPTION_KEY=$(generate_secret 32)
SESSION_ENCRYPTION_KEY=$(generate_secret 32)
DB_PASSWORD=$(generate_password)
DB_REPLICA_PASSWORD=$(generate_password)
POSTGRES_PASSWORD=$(generate_password)
REPLICATION_PASSWORD=$(generate_password)

echo "=== Generated Secrets ==="
echo ""
echo "JWT_SECRET_KEY:"
echo "$JWT_SECRET"
echo ""
echo "DATA_ENCRYPTION_KEY:"
echo "$DATA_ENCRYPTION_KEY"
echo ""
echo "SESSION_ENCRYPTION_KEY:"
echo "$SESSION_ENCRYPTION_KEY"
echo ""
echo "DB_PASSWORD:"
echo "$DB_PASSWORD"
echo ""
echo "DB_REPLICA_PASSWORD:"
echo "$DB_REPLICA_PASSWORD"
echo ""
echo "POSTGRES_PASSWORD:"
echo "$POSTGRES_PASSWORD"
echo ""
echo "REPLICATION_PASSWORD:"
echo "$REPLICATION_PASSWORD"
echo ""

# Ask if user wants to create a secrets file
echo ""
echo -e "${YELLOW}Would you like to create a secrets-generated.yaml file with these values? [y/N]${NC}"
read -r response

if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    OUTPUT_FILE="secrets-generated.yaml"
    
    cat > "$OUTPUT_FILE" <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: insavein-secrets
  namespace: insavein
  labels:
    app: insavein-platform
type: Opaque
stringData:
  # Database Credentials
  DB_USER: "insavein_user"
  DB_PASSWORD: "$DB_PASSWORD"
  DB_REPLICA_USER: "insavein_replica_user"
  DB_REPLICA_PASSWORD: "$DB_REPLICA_PASSWORD"
  
  # JWT Secret Key
  JWT_SECRET_KEY: "$JWT_SECRET"
  
  # Email Service API Keys
  # TODO: Replace with actual email service credentials
  EMAIL_API_KEY: "CHANGE_ME_EMAIL_SERVICE_API_KEY"
  EMAIL_API_SECRET: "CHANGE_ME_EMAIL_SERVICE_API_SECRET"
  
  # Push Notification Service
  # TODO: Replace with actual FCM credentials
  FCM_SERVER_KEY: "CHANGE_ME_FCM_SERVER_KEY"
  FCM_PROJECT_ID: "CHANGE_ME_FCM_PROJECT_ID"
  
  # External API Keys
  # TODO: Replace with actual OpenAI API key
  OPENAI_API_KEY: "CHANGE_ME_OPENAI_API_KEY"
  
  # Encryption Keys
  DATA_ENCRYPTION_KEY: "$DATA_ENCRYPTION_KEY"
  SESSION_ENCRYPTION_KEY: "$SESSION_ENCRYPTION_KEY"

---
apiVersion: v1
kind: Secret
metadata:
  name: postgres-credentials
  namespace: insavein
  labels:
    app: postgres
type: Opaque
stringData:
  POSTGRES_USER: "postgres"
  POSTGRES_PASSWORD: "$POSTGRES_PASSWORD"
  REPLICATION_USER: "replicator"
  REPLICATION_PASSWORD: "$REPLICATION_PASSWORD"
EOF

    echo ""
    echo -e "${GREEN}✓ Secrets file created: $OUTPUT_FILE${NC}"
    echo ""
    echo -e "${YELLOW}IMPORTANT:${NC}"
    echo "1. Review the file and replace remaining CHANGE_ME placeholders"
    echo "2. Add email service credentials (SendGrid/AWS SES)"
    echo "3. Add Firebase Cloud Messaging credentials"
    echo "4. Add OpenAI API key (if using AI recommendations)"
    echo "5. Store this file securely (DO NOT commit to version control)"
    echo "6. Apply with: kubectl apply -f $OUTPUT_FILE"
    echo ""
    echo -e "${RED}⚠️  SECURITY WARNING:${NC}"
    echo "This file contains sensitive credentials. Keep it secure!"
    echo ""
else
    echo ""
    echo "Secrets not saved to file. Copy the values above manually."
    echo ""
fi

echo "=================================================="
echo "Secret generation complete!"
echo "=================================================="
