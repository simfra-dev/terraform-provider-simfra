data "simfra_accounts" "all" {}

output "account_ids" {
  value = [for a in data.simfra_accounts.all.accounts : a.account_id]
}
