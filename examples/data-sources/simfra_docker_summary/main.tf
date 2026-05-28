data "simfra_docker_summary" "info" {}

output "containers_running" {
  value = data.simfra_docker_summary.info.containers_running
}

output "images" {
  value = data.simfra_docker_summary.info.images
}
