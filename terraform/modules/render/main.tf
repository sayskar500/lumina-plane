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
variable "mongo_uri" { type = any }
variable "env_vars" { type = any }

resource "render_web_service" "api" {
  name      = var.app_name
  region    = "oregon"
  plan      = "starter"

  runtime_source = {
    docker = {
      repo_url = var.repo_url
      branch   = "main"
    }
  }

  # We set MONGO_URI to a placeholder to avoid the 'computed object' type error.
  # We then merge the other environment variables (Groq, Infisical) which are plain strings.
  env_vars = merge(
    {
      MONGO_URI = { value = "PENDING_ATLAS_URI" }
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
