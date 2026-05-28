---
page_title: "simfra_port_forwards Data Source - terraform-provider-simfra"
subcategory: ""
description: |-
  Lists all active Simfra port forwards.
---

# simfra_port_forwards (Data Source)

Lists all active port forward sessions in the Simfra instance. Port forwards allow access to private VPC resources from the host.

## Example Usage

```terraform
data "simfra_port_forwards" "all" {}

output "active_forwards" {
  value = [for pf in data.simfra_port_forwards.all.port_forwards : {
    id         = pf.id
    local_port = pf.local_port
    service    = pf.service
  }]
}
```

## Schema

### Read-Only

- `port_forwards` (List of Object) -- List of active port forwards. Each element has the following attributes:
  - `id` (String) -- Port forward session ID.
  - `target_arn` (String) -- ARN of the target resource.
  - `local_port` (Number) -- Port on the local host.
  - `local_address` (String) -- Local bind address.
  - `service` (String) -- AWS service type of the target.
  - `status` (String) -- Port forward status.
