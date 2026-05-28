---
page_title: "simfra_route53_zone Resource - terraform-provider-simfra"
subcategory: ""
description: |-
  Creates a Route53 hosted zone with an optional specific zone ID.
---

# simfra_route53_zone (Resource)

Creates a Route53 hosted zone in the Simfra account. This resource allows you to pre-seed Route53 zones with deterministic zone IDs, which is useful for matching production infrastructure IDs in local development and testing.

This resource requires AWS credentials (`access_key` and `secret_key`) to be configured on the provider.

## Example Usage

```terraform
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
```

## Schema

### Required

- `name` (String) -- Domain name for the hosted zone (e.g. `example.com`). Changing this forces a new resource.

### Optional

- `id` (String) -- Hosted zone ID. If specified, the zone is created with this exact ID. If omitted, Simfra generates one. Changing this forces a new resource.
- `comment` (String) -- Comment for the hosted zone.

### Read-Only

- `name_servers` (List of String) -- Delegation name servers for the hosted zone.
