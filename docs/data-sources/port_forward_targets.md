---
page_title: "simfra_port_forward_targets Data Source - terraform-provider-simfra"
subcategory: ""
description: |-
  Lists all available port forward targets.
---

# simfra_port_forward_targets (Data Source)

Lists all resources that are available as port forward targets. These are private VPC resources (such as RDS instances, ElastiCache clusters) that can be reached from the host via port forwarding.

## Example Usage

```terraform
data "simfra_port_forward_targets" "all" {}

output "available_targets" {
  value = [for t in data.simfra_port_forward_targets.all.targets : {
    arn     = t.arn
    service = t.service
  }]
}
```

## Schema

### Read-Only

- `targets` (List of Object) -- List of available port forward targets. Each element has the following attributes:
  - `arn` (String) -- ARN of the target resource.
  - `service` (String) -- AWS service type (e.g. `rds`, `elasticache`).
  - `resource_id` (String) -- Resource identifier.
  - `account_id` (String) -- AWS account ID.
  - `region` (String) -- AWS region.
  - `default_port` (Number) -- Default port for the service.
  - `vpc_network` (String) -- Docker VPC network name.
