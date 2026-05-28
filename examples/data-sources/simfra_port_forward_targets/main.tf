data "simfra_port_forward_targets" "all" {}

output "available_targets" {
  value = [for t in data.simfra_port_forward_targets.all.targets : {
    arn     = t.arn
    service = t.service
  }]
}
