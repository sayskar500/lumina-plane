# Render App Module
terraform {
  required_providers {
    render = {
      source  = "render-oss/render"
      version = "~> 1.0"
    }
  }
}

variable "app_name" { type = string }
variable "repo_url" { type = string }

# Real Atlas connection string. Because the root module passes
# module.mongodb_cluster.connection_string here, Terraform orders the Atlas
# cluster creation first, so the value is known by the time the service is
# created — no post-apply API patching needed.
variable "mongo_uri" { type = string }
variable "env_vars" { type = any }

resource "render_web_service" "api" {
  name   = var.app_name
  region = "oregon"
  # NOTE: the render-oss provider documents only paid plans (starter, standard,
  # pro, ...). The Render API also accepts "free" for web services, but creating
  # ANY service via the API requires a payment method on file (HTTP 402).
  plan = "free"

  runtime_source = {
    docker = {
      repo_url        = var.repo_url
      branch          = "main"
      dockerfile_path = "Dockerfile.server" # repo has no root Dockerfile
    }
  }

  env_vars = merge(
    {
      MONGO_URI = { value = var.mongo_uri }
    },
    {
      for k, v in var.env_vars : k => {
        value = v
      }
    }
  )
}

output "service_id" {
  value = render_web_service.api.id
}

output "service_url" {
  value = render_web_service.api.url
}
