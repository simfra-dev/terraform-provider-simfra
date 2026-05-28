---
page_title: "simfra_dns_port Data Source - terraform-provider-simfra"
subcategory: ""
description: |-
  Retrieves the DNS server port for a Simfra account.
---

# simfra_dns_port (Data Source)

Retrieves the DNS server port for a specific Simfra account. Each account has its own DNS container that resolves Route53 hosted zones and service DNS names. The returned port can be used to configure local DNS resolution to point at the Simfra DNS server.

## Example Usage

```terraform
data "simfra_dns_port" "dev" {
  account_id = "123456789012"
}

output "dns_port" {
  value = data.simfra_dns_port.dev.port
}
```

## Schema

### Required

- `account_id` (String) -- AWS account ID.

### Read-Only

- `port` (Number) -- DNS server port.
