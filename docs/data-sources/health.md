---
page_title: "simfra_health Data Source - terraform-provider-simfra"
subcategory: ""
description: |-
  Checks Simfra server health status.
---

# simfra_health (Data Source)

Checks whether the Simfra server is healthy and ready to accept requests. Useful as a dependency to ensure the server is running before creating resources.

## Example Usage

```terraform
data "simfra_health" "check" {}

output "simfra_ready" {
  value = data.simfra_health.check.ready
}
```

## Schema

### Read-Only

- `ready` (Boolean) -- Whether the Simfra server is ready.
