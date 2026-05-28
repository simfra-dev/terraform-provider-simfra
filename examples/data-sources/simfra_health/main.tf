data "simfra_health" "check" {}

output "simfra_ready" {
  value = data.simfra_health.check.ready
}
