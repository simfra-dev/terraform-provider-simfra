data "simfra_ca_info" "info" {}

output "root_subject" {
  value = data.simfra_ca_info.info.root_subject
}

output "root_expires" {
  value = data.simfra_ca_info.info.root_not_after
}
