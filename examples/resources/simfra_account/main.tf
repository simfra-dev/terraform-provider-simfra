resource "simfra_account" "dev" {
  account_id = "123456789012"

  bootstrap        = "standard"
  bootstrap_region = "us-east-1"
}

output "root_access_key_id" {
  value = simfra_account.dev.root_access_key_id
}

output "root_secret_access_key" {
  value     = simfra_account.dev.root_secret_access_key
  sensitive = true
}
