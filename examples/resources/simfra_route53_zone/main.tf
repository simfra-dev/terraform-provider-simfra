# Create a Route53 hosted zone with a specific zone ID
resource "simfra_route53_zone" "example" {
  id      = "Z0123456789ABCDEFGHIJ"
  name    = "example.com"
  comment = "Production domain"
}

# Create a Route53 hosted zone with an auto-generated ID
resource "simfra_route53_zone" "internal" {
  name    = "internal.example.com"
  comment = "Internal services"
}

output "zone_id" {
  value = simfra_route53_zone.example.id
}

output "name_servers" {
  value = simfra_route53_zone.example.name_servers
}
