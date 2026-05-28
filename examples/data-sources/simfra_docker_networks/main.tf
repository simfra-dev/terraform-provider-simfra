data "simfra_docker_networks" "all" {}

output "network_names" {
  value = [for n in data.simfra_docker_networks.all.networks : n.name]
}
