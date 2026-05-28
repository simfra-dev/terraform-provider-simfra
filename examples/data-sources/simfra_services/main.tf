data "simfra_services" "all" {}

output "service_names" {
  value = [for s in data.simfra_services.all.services : s.name]
}
