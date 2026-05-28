data "simfra_account" "dev" {
  account_id = "123456789012"
}

output "root_access_key_id" {
  value = data.simfra_account.dev.root_access_key_id
}
