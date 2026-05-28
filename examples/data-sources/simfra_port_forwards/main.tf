data "simfra_port_forwards" "all" {}

output "active_forwards" {
  value = [for pf in data.simfra_port_forwards.all.port_forwards : {
    id         = pf.id
    local_port = pf.local_port
    service    = pf.service
  }]
}
