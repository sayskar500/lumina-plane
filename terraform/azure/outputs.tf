output "web_app_url" {
  description = "Public URL of the deployed AI Gateway."
  value       = "https://${azurerm_linux_web_app.app.default_hostname}"
}

output "web_app_name" {
  description = "App Service name (for az CLI commands like log tail)."
  value       = azurerm_linux_web_app.app.name
}

output "resource_group_name" {
  description = "Resource group containing all Lumina-Plane Azure resources."
  value       = azurerm_resource_group.rg.name
}

output "cosmos_account_name" {
  description = "Cosmos DB account name (MongoDB API)."
  value       = azurerm_cosmosdb_account.db.name
}

output "mongo_connection_string" {
  description = "MongoDB connection string (contains the account key — sensitive)."
  value       = azurerm_cosmosdb_account.db.primary_mongodb_connection_string
  sensitive   = true
}
