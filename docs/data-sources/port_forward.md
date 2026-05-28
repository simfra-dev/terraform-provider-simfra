---
page_title: "simfra_port_forward Data Source - terraform-provider-simfra"
subcategory: ""
description: |-
  Retrieves details of a Simfra port forward by ID.
---

# simfra_port_forward (Data Source)

Retrieves details of a specific port forward session by its ID. Port forwards allow access to private VPC resources from the host by forwarding a local port to a container's internal port.

## Example Usage

```terraform
data "simfra_port_forward" "example" {
  id = "pf-abc123"
}

output "local_port" {
  value = data.simfra_port_forward.example.local_port
}

output "target_arn" {
  value = data.simfra_port_forward.example.target_arn
}
```

## Schema

### Required

- `id` (String) -- Port forward session ID.

### Read-Only

- `target_arn` (String) -- ARN of the target resource.
- `target_ip` (String) -- IP address of the target container.
- `target_port` (Number) -- Port on the target container.
- `local_port` (Number) -- Port on the local host.
- `local_address` (String) -- Local bind address.
- `vpc_network` (String) -- Docker VPC network name.
- `container_id` (String) -- Docker container ID of the socat relay.
- `service` (String) -- AWS service type of the target (e.g. `rds`, `elasticache`).
- `resource_id` (String) -- Resource identifier of the target.
- `account_id` (String) -- AWS account ID.
- `region` (String) -- AWS region.
- `created_at` (String) -- Port forward creation timestamp.
- `status` (String) -- Port forward status (e.g. `active`).
