# MongoDB Atlas Module
terraform {
  required_providers {
    mongodbatlas = {
      source  = "mongodb/mongodbatlas"
      version = "~> 1.0"
    }
  }
}

variable "project_id" { type = string }
variable "region" { default = "US_EAST_1" }

resource "mongodbatlas_cluster" "cluster" {
  name       = "lumina-cluster"
  project_id = var.project_id

  # Correct configuration for M0 (Free Tier)
  provider_name               = "TENANT"
  backing_provider_name       = "AWS"
  provider_instance_size_name = "M0"

  # Required by the API for M0 clusters
  region_name = var.region

  replication_specs {
    num_shards = 1
    regions_config {
      region_name     = var.region
      electable_nodes = 3
      priority        = 7
    }
  }
}

output "connection_string" {
  value = mongodbatlas_cluster.cluster.connection_strings[0]
}
