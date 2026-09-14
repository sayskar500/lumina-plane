#!/bin/bash
set -e

# Load environment variables from tfvars for the API key
# We use grep/awk to extract the render_api_key from the .tfvars file
RENDER_API_KEY=$(grep "render_api_key" terraform/terraform.tfvars | awk -F'"' '{print $2}')

if [ -z "$RENDER_API_KEY" ]; then
    echo "❌ Error: render_api_key not found in terraform/terraform.tfvars"
    exit 1
fi

echo "🚀 Starting Lumina-Plane Cloud Deployment..."

# 1. Auto-discover Render Owner ID
echo "🔍 Discovering Render Owner ID..."
OWNER_RESPONSE=$(curl -s -H "Authorization: Bearer $RENDER_API_KEY" https://api.render.com/v1/owners)
OWNER_ID=$(echo "$OWNER_RESPONSE" | grep -o '"id":"[^"]*"' | head -n 1 | cut -d'"' -f4)

if [ -z "$OWNER_ID" ]; then
    echo "❌ Error: Could not discover Render Owner ID. Please check your API Key."
    echo "API Response: $OWNER_RESPONSE"
    exit 1
fi

echo "✅ Discovered Owner ID: $OWNER_ID"

# 2. Export variables for Terraform
# Terraform picks up TF_VAR_name as the value for var.name
export TF_VAR_render_api_key="$RENDER_API_KEY"
export TF_VAR_render_owner_id="$OWNER_ID"

# 3. Provision Infrastructure
echo "📦 Provisioning Cloud Infrastructure..."
cd terraform
terraform init
terraform apply -auto-approve

# 4. Extract Outputs
echo "🔍 Extracting deployment metadata..."
SERVICE_ID=$(terraform output -raw render_service_id)
MONGO_URI=$(terraform output -raw mongodb_connection_string)
SERVICE_URL=$(terraform output -raw render_service_url)

if [ -z "$SERVICE_ID" ] || [ -z "$MONGO_URI" ]; then
    echo "❌ Error: Failed to extract service ID or MongoDB URI."
    exit 1
fi

# 5. Patch the Render Service with the actual MongoDB URI
echo "🔗 Linking MongoDB Atlas to Render Service..."
curl -X PATCH "https://api.render.com/v1/services/$SERVICE_ID" \
     -H "Authorization: Bearer $RENDER_API_KEY" \
     -H "Content-Type: application/json" \
     -d "{
       \"env_vars\": {
         \"MONGO_URI\": { \"value\": \"$MONGO_URI\" }
       }
     }"

echo ""
echo "✅ Deployment Complete!"
echo "--------------------------------------------------"
echo "🌐 Live URL: $SERVICE_URL"
echo "🗄️ MongoDB URI: $MONGO_URI"
echo "--------------------------------------------------"
echo "You can now test the platform using: ./bin/lumina ask \"Hello Cloud!\""
