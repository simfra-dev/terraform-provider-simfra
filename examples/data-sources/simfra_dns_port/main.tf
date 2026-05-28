data "simfra_dns_port" "dev" {
  account_id = "123456789012"
}

output "dns_port" {
  value = data.simfra_dns_port.dev.port
}
