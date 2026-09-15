variable "location" {
  description = "Azure region. Free tiers are available in all public regions."
  type        = string
  default     = "centralindia"
}

variable "resource_group_name" {
  description = "Resource group for all Lumina-Plane Azure resources."
  type        = string
  default     = "lumina-plane-rg"
}

variable "mongo_database_name" {
  description = "MongoDB database the app uses (shared free-tier throughput)."
  type        = string
  default     = "lumina_plane"
}

variable "groq_api_key" {
  description = "Groq API key injected into the app as GROQ_API_KEY."
  type        = string
  sensitive   = true
}

variable "container_image" {
  description = "Container image (name:tag) to run, published by GitHub Actions."
  type        = string
  default     = "ghcr.io/sayskar500/lumina-plane:latest"
}

variable "container_registry_url" {
  description = "Registry URL the image is pulled from."
  type        = string
  default     = "https://ghcr.io"
}

variable "container_registry_username" {
  description = "Registry username; leave empty for public images."
  type        = string
  default     = ""
  sensitive   = true
}

variable "container_registry_password" {
  description = "Registry password/token; leave empty for public images."
  type        = string
  default     = ""
  sensitive   = true
}
