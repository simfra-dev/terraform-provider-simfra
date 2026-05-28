---
page_title: "simfra_docker_networks Data Source - terraform-provider-simfra"
subcategory: ""
description: |-
  Lists all Simfra-managed Docker networks.
---

# simfra_docker_networks (Data Source)

Lists all Docker networks managed by the Simfra instance. These networks provide VPC isolation for containers when VPC network isolation is enabled.

## Example Usage

```terraform
data "simfra_docker_networks" "all" {}

output "network_names" {
  value = [for n in data.simfra_docker_networks.all.networks : n.name]
}
```

## Schema

### Read-Only

- `networks` (List of Object) -- List of Docker networks. Each element has the following attributes:
  - `id` (String) -- Docker network ID.
  - `name` (String) -- Network name.
  - `driver` (String) -- Network driver (e.g. `bridge`).
  - `scope` (String) -- Network scope (e.g. `local`).
