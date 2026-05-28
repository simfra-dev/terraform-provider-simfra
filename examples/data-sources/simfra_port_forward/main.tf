data "simfra_port_forward" "example" {
  id = "pf-abc123"
}

output "local_port" {
  value = data.simfra_port_forward.example.local_port
}

output "target_arn" {
  value = data.simfra_port_forward.example.target_arn
}
