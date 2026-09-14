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

  # Required when replication_specs is set
  cluster_type = "REPLICASET"

  # For M0 Free Tier (shared tier). The Atlas API requires the region in
  # providerSettings; with this provider that maps to provider_region_name.
  provider_name               = "TENANT"
  backing_provider_name       = "AWS"
  provider_instance_size_name = "M0"
  provider_region_name        = var.region

  replication_specs {
    num_shards = 1
    regions_config {
      region_name     = var.region
      electable_nodes = 3
      priority        = 7
    }
  }
}

# Render (especially the free tier) egresses from shared, non-static IPs, so
# the project must allow all inbound addresses or the app cannot reach Atlas.
resource "mongodbatlas_project_ip_access_list" "render_egress" {
  project_id = var.project_id
  cidr_block = "0.0.0.0/0"
  comment    = "Render shared egress (see docs/CLOUD.md)"
}

# connection_strings is a list of objects; standard_srv is the SRV string the
# app needs (mongodb+srv://...). Note: it carries no credentials — see the DB
# user note in docs/CLOUD.md.
output "connection_string" {
  value = mongodbatlas_cluster.cluster.connection_strings[0].standard_srv
}
