data "simfra_smtp_port" "dev" {
  account_id = "123456789012"
}

output "smtp_port" {
  value = data.simfra_smtp_port.dev.port
}
