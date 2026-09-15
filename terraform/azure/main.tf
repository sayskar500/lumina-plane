# Lumina-Plane on Azure — free-tier architecture (no credit card required)
#
#   Azure App Service (F1 Free, Linux)  ← runs the Docker image from GHCR
#   Azure Cosmos DB (MongoDB API)       ← free tier: 1000 RU/s + 25 GB forever
#
# Designed for the "Azure for Students" subscription ($100 credit, no card):
# every resource here sits inside a free tier or a hard-capped quota, so the
# deployment cannot incur charges.

terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

provider "azurerm" {
  # ARM_SUBSCRIPTION_ID (set by scripts/deploy-azure.sh) selects the
  # subscription; required by azurerm v4 when using Azure CLI auth.
  features {}
}

# Globally-unique suffix for the Cosmos account and web app names.
resource "random_string" "suffix" {
  length  = 6
  lower   = true
  upper   = false
  special = false
}

resource "azurerm_resource_group" "rg" {
  name     = var.resource_group_name
  location = var.location
}

# ---------------------------------------------------------------------------
# Database: Azure Cosmos DB for MongoDB (RU-based), free tier
# ---------------------------------------------------------------------------
resource "azurerm_cosmosdb_account" "db" {
  name                = "lumina-cosmos-${random_string.suffix.result}"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  offer_type          = "Standard"
  kind                = "MongoDB"

  # One free-tier Cosmos account per subscription: 1000 RU/s + 25 GB free
  # forever, no credit burn.
  free_tier_enabled = true

  # Spending guardrail (by construction): the only throughput allocation is
  # the shared-throughput database below at exactly 1000 RU/s — the free
  # quota — and collections inherit it instead of allocating their own.
  mongo_server_version = "7.0"
  minimal_tls_version  = "Tls12"

  consistency_policy {
    consistency_level = "Session"
  }

  geo_location {
    location          = azurerm_resource_group.rg.location
    failover_priority = 0
    zone_redundant    = false
  }
}

# Shared-throughput database (draws from the free 1000 RU/s). The app's
# collections (prompts, token_logs, counters) are created lazily inside it.
resource "azurerm_cosmosdb_mongo_database" "lumina" {
  name                = var.mongo_database_name
  resource_group_name = azurerm_resource_group.rg.name
  account_name        = azurerm_cosmosdb_account.db.name
  throughput          = 1000
}

# ---------------------------------------------------------------------------
# Compute: Azure App Service (Linux, F1 Free) running the container
# ---------------------------------------------------------------------------
resource "azurerm_service_plan" "plan" {
  name                = "lumina-asp-${random_string.suffix.result}"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  os_type             = "Linux"
  sku_name            = "F1" # Free tier
}

resource "azurerm_linux_web_app" "app" {
  name                = "lumina-plane-server-${random_string.suffix.result}"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  service_plan_id     = azurerm_service_plan.plan.id

  app_settings = {
    # Our container listens on 8000 (from config/config.yaml baked into the image)
    WEBSITES_PORT                       = "8000"
    WEBSITES_ENABLE_APP_SERVICE_STORAGE = "false"

    # 12-factor overrides consumed by the server (env > YAML)
    MONGO_URI    = azurerm_cosmosdb_account.db.primary_mongodb_connection_string
    GROQ_API_KEY = var.groq_api_key
  }

  identity {
    type = "SystemAssigned"
  }

  site_config {
    health_check_path       = "/health"
    scm_minimum_tls_version = "Tls12"
    http2_enabled           = true

    application_stack {
      docker_image_name   = var.container_image
      docker_registry_url = var.container_registry_url

      # Only needed if the image is in a private registry (e.g. GHCR before
      # the package is made public, or Docker Hub private repos).
      docker_registry_username = var.container_registry_username
      docker_registry_password = var.container_registry_password
    }
  }
}
