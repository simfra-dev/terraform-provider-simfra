data "simfra_docker_containers" "all" {}

output "container_names" {
  value = [for c in data.simfra_docker_containers.all.containers : c.name]
}
