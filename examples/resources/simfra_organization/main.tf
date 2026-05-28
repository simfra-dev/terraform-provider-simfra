# Create an AWS Organization with specific IDs
resource "simfra_organization" "example" {
  id          = "o-abc1234567"
  root_id     = "r-ab12"
  feature_set = "ALL"
}

# Create an AWS Organization with auto-generated IDs
resource "simfra_organization" "auto" {
  feature_set = "ALL"
}

output "organization_arn" {
  value = simfra_organization.example.arn
}

output "master_account_id" {
  value = simfra_organization.example.master_account_id
}
