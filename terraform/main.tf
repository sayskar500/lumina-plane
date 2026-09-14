# Terraform configuration for Lumina-Plane
# This showcases Infrastructure as Code (IaC) skills

terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.0"
    }
    render = {
      source = "render-oss/render"
      version = "~> 1.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

provider "mongodbatlas" {
  public_key  = var.atlas_public_key
  private_key = var.atlas_private_key
}
provider "render" {
  api_key  = var.render_api_key
  owner_id = var.render_owner_id
}

variable "render_api_key" {
  type      = string
  sensitive = true
}

variable "render_owner_id" {
  type      = string
}

variable "atlas_public_key" {
  type      = string
  sensitive = true
}

variable "atlas_private_key" {
  type      = string
  sensitive = true
}

variable "atlas_project_id" {
  type      = string
}

# Generated app-credential for the Atlas database user; stored in state
# (gitignored) and injected into the Render service env.
resource "random_password" "atlas_db" {
  length  = 24
  special = false # keeps the connection string URL-safe
}

locals {
  # Using a local variable to "flatten" the computed output from the module
  mongo_uri = module.mongodb_cluster.connection_string
}

module "mongodb_cluster" {
  source      = "./modules/mongodb"
  project_id  = var.atlas_project_id
  db_password = random_password.atlas_db.result
}

module "app_server" {
  source = "./modules/render"
  app_name = "lumina-plane-server"
  repo_url = "https://github.com/sayskar500/lumina-plane"
  mongo_uri = local.mongo_uri
  env_vars = {
    GROQ_API_KEY         = var.groq_api_key
    INFISICAL_CLIENT_ID   = var.infisical_client_id
    INFISICAL_CLIENT_SECRET = var.infisical_client_secret
  }
}

output "render_service_id" {
  value = module.app_server.service_id
}

output "render_service_url" {
  value = module.app_server.service_url
}

output "mongodb_connection_string" {
  value     = module.mongodb_cluster.connection_string
  sensitive = true
}

variable "groq_api_key" {
  type      = string
  sensitive = true
}

variable "infisical_client_id" {
  type      = string
  sensitive = true
}

variable "infisical_client_secret" {
  type      = string
  sensitive = true
}
