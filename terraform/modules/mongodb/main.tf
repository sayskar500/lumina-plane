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
variable "db_username" { default = "lumina_app" }
variable "db_password" {
  type      = string
  sensitive = true
}

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

# App database user with least-privilege access (readWrite on lumina_plane).
# Note: user propagation on Atlas can take a minute or two after creation.
resource "mongodbatlas_database_user" "app" {
  project_id         = var.project_id
  username           = var.db_username
  password           = var.db_password
  auth_database_name = "admin"

  roles {
    role_name     = "readWrite"
    database_name = "lumina_plane"
  }
}

# Render (especially the free tier) egresses from shared, non-static IPs, so
# the project must allow all inbound addresses or the app cannot reach Atlas.
resource "mongodbatlas_project_ip_access_list" "render_egress" {
  project_id = var.project_id
  cidr_block = "0.0.0.0/0"
  comment    = "Render shared egress (see docs/CLOUD.md)"
}

# connection_strings is a list of objects; standard_srv provides the SRV host.
# We embed the app user's credentials so the service can actually authenticate,
# and pin the database name with standard Atlas driver options.
output "connection_string" {
  value = "mongodb+srv://${var.db_username}:${var.db_password}@${replace(mongodbatlas_cluster.cluster.connection_strings[0].standard_srv, "mongodb+srv://", "")}/lumina_plane?retryWrites=true&w=majority"
}
