data "simfra_storage_summary" "info" {}

output "persistence_enabled" {
  value = data.simfra_storage_summary.info.enabled
}

output "db_size" {
  value = data.simfra_storage_summary.info.db_size
}
