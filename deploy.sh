#!/bin/bash
set -euo pipefail

# Resolve the repo root from the script location so this works from any CWD.
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TF_DIR="$SCRIPT_DIR/terraform"

for tool in terraform jq curl; do
  command -v "$tool" >/dev/null 2>&1 || { echo "❌ Required tool '$tool' not found in PATH."; exit 1; }
done

[ -f "$TF_DIR/terraform.tfvars" ] || { echo "❌ $TF_DIR/terraform.tfvars not found. Copy the template from docs/CLOUD.md and fill in your keys."; exit 1; }

RENDER_API_KEY=$(awk -F'"' '/^render_api_key/ {print $2}' "$TF_DIR/terraform.tfvars")
if [ -z "$RENDER_API_KEY" ]; then
    echo "❌ render_api_key not found in terraform/terraform.tfvars"
    exit 1
fi

echo "🚀 Starting Lumina-Plane Cloud Deployment..."

# 1. Auto-discover the Render workspace (owner) ID
echo "🔍 Discovering Render workspace..."
if ! OWNER_RESPONSE=$(curl -sfS -H "Authorization: Bearer $RENDER_API_KEY" -H "Accept: application/json" "https://api.render.com/v1/owners?limit=1"); then
    echo "❌ Could not reach the Render API. Check your render_api_key and network."
    exit 1
fi
OWNER_ID=$(echo "$OWNER_RESPONSE" | jq -r '.[0].owner.id // empty')
if [ -z "$OWNER_ID" ]; then
    echo "❌ Could not discover Render workspace ID. API response:"
    echo "$OWNER_RESPONSE"
    exit 1
fi
echo "✅ Discovered workspace: $OWNER_ID"

# 2. Export variables for Terraform (TF_VAR_<name> maps to var.<name>)
export TF_VAR_render_api_key="$RENDER_API_KEY"
export TF_VAR_render_owner_id="$OWNER_ID"

# 3. Provision infrastructure.
#    Order is enforced by Terraform: the Atlas cluster is created first, and the
#    Render service is created with the real MONGO_URI (no post-apply patching).
echo "📦 Provisioning cloud infrastructure (Atlas M0 can take 5-10 minutes)..."
APPLY_LOG="$TF_DIR/.apply.log"
terraform -chdir="$TF_DIR" init -input=false >/dev/null
APPLY_OK=1
if ! terraform -chdir="$TF_DIR" apply -auto-approve -input=false 2>&1 | tee "$APPLY_LOG"; then
    APPLY_OK=0
fi

# The render provider cannot update free-tier services in place (it sends a
# maintenance-mode field the API rejects). Recreating the service (same name
# and URL, no data loss) is the supported workaround.
if [ "$APPLY_OK" -eq 0 ] && grep -q "maintenance mode can only be configured for non-free tier" "$APPLY_LOG"; then
    echo "⚠️  Free-tier services cannot be updated in place — recreating the Render service..."
    terraform -chdir="$TF_DIR" apply -replace="module.app_server.render_web_service.api" -auto-approve -input=false 2>&1 | tee "$APPLY_LOG" || APPLY_OK=0
elif [ "$APPLY_OK" -eq 0 ]; then
    echo ""
    echo "❌ Terraform apply failed."
    if grep -q "Payment information" "$APPLY_LOG"; then
        echo "   → Render requires a payment method on file to create services via the"
        echo "     API — even on the free plan. Add one at https://dashboard.render.com/billing"
        echo "     then re-run this script."
    fi
    if grep -qi "unauthorized\|401" "$APPLY_LOG"; then
        echo "   → A cloud API key appears to be invalid or lacks permissions."
    fi
    exit 1
fi

# 4. Extract outputs
echo "🔍 Extracting deployment metadata..."
SERVICE_ID=$(terraform -chdir="$TF_DIR" output -raw render_service_id)
SERVICE_URL=$(terraform -chdir="$TF_DIR" output -raw render_service_url)

if [ -z "$SERVICE_ID" ] || [ -z "$SERVICE_URL" ]; then
    echo "❌ Failed to extract service ID or URL from Terraform outputs."
    exit 1
fi

# 5. Verify the Render service received its environment variables (keys only —
#    values may contain secrets). The service's first build/deploy starts
#    automatically on creation.
echo "🔑 Verifying service environment variables..."
ENV_KEYS=$(curl -sfS -H "Authorization: Bearer $RENDER_API_KEY" -H "Accept: application/json" \
    "https://api.render.com/v1/services/$SERVICE_ID/env-vars" | jq -r 'if type=="array" then [.[].envVar.key] | join(", ") else . end') || ENV_KEYS=""
if echo "$ENV_KEYS" | grep -q "MONGO_URI"; then
    echo "✅ MONGO_URI is set on the service (along with: ${ENV_KEYS})"
else
    echo "⚠️  Could not confirm MONGO_URI on the service. Keys found: ${ENV_KEYS:-none}"
fi

echo ""
echo "✅ Deployment Complete!"
echo "--------------------------------------------------"
echo "🌐 Live URL: $SERVICE_URL"
echo "🗄️ MongoDB: managed via Terraform (mongodb_connection_string output; not printed — may contain credentials)"
echo "--------------------------------------------------"
echo "Try it:"
echo "  make build"
echo "  LUMINA_SERVER_URL=\"$SERVICE_URL\" ./bin/lumina health"
echo "  LUMINA_SERVER_URL=\"$SERVICE_URL\" ./bin/lumina ask \"Hello Cloud!\""
