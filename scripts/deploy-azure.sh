#!/bin/bash
# Deploy Lumina-Plane to Azure free tiers (App Service F1 + Cosmos DB free).
#
# Prerequisites (see docs/AZURE.md):
#   1. An Azure subscription WITHOUT a credit card: "Azure for Students"
#      (https://azure.microsoft.com/free/students) or via the GitHub Student
#      Developer Pack.
#   2. Azure CLI installed and logged in:  az login
#   3. Groq API key in terraform/terraform.tfvars (shared with the Render
#      deployment) or exported as TF_VAR_groq_api_key.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
TF_DIR="$ROOT_DIR/terraform/azure"

for tool in terraform jq curl az; do
  command -v "$tool" >/dev/null 2>&1 || {
    echo "❌ Required tool '$tool' not found."
    [ "$tool" = "az" ] && echo "   Install the Azure CLI: brew install azure-cli"
    exit 1
  }
done

# 1. Azure authentication
echo "🔍 Checking Azure login..."
if ! az account show >/dev/null 2>&1; then
  echo "❌ Not logged in to Azure. Run:  az login"
  exit 1
fi
SUB_ID=$(az account show --query id -o tsv)
SUB_NAME=$(az account show --query name -o tsv)
echo "✅ Subscription: $SUB_NAME ($SUB_ID)"
export ARM_SUBSCRIPTION_ID="$SUB_ID"

# 2. Groq key: reuse the existing Render tfvars unless already exported
if [ -z "${TF_VAR_groq_api_key:-}" ] && [ -f "$ROOT_DIR/terraform/terraform.tfvars" ]; then
  TF_VAR_groq_api_key=$(awk -F'"' '/^groq_api_key/ {print $2}' "$ROOT_DIR/terraform/terraform.tfvars")
  export TF_VAR_groq_api_key
fi
if [ -z "${TF_VAR_groq_api_key:-}" ]; then
  echo "❌ No Groq API key found. Add groq_api_key to terraform/terraform.tfvars"
  echo "   or export TF_VAR_groq_api_key."
  exit 1
fi

# 3. Terraform
echo "🚀 Deploying Lumina-Plane to Azure (App Service F1 + Cosmos DB free tier)..."
terraform -chdir="$TF_DIR" init -input=false >/dev/null

APPLY_LOG="$TF_DIR/.apply.log"
if ! terraform -chdir="$TF_DIR" apply -auto-approve -input=false 2>&1 | tee "$APPLY_LOG"; then
  echo ""
  echo "❌ Terraform apply failed."
  if grep -q "FreeTrialsAndStudentSubscriptionsNotEligible\|not eligible\|subscription" "$APPLY_LOG"; then
    echo "   → Your subscription may not support one of the free tiers. Check the"
    echo "     error above; Cosmos DB free tier is limited to one account per subscription."
  fi
  exit 1
fi

# 4. Output
URL=$(terraform -chdir="$TF_DIR" output -raw web_app_url)
APP=$(terraform -chdir="$TF_DIR" output -raw web_app_name)
RG=$(terraform -chdir="$TF_DIR" output -raw resource_group_name)

echo ""
echo "✅ Azure Deployment Complete!"
echo "--------------------------------------------------"
echo "🌐 Live URL: $URL"
echo "📦 App:      $APP (resource group: $RG)"
echo "--------------------------------------------------"
echo "Test it (first request may take ~1 min — F1 cold start):"
echo "  curl $URL/health"
echo "  LUMINA_SERVER_URL=\"$URL\" ./bin/lumina ask \"Hello Azure!\""
echo ""
echo "Stream logs:  az webapp log tail --name $APP --resource-group $RG"
