---
page_title: "simfra_docker_containers Data Source - terraform-provider-simfra"
subcategory: ""
description: |-
  Lists all Simfra-managed Docker containers.
---

# simfra_docker_containers (Data Source)

Lists all Docker containers managed by the Simfra instance. This includes containers for services that require runtime processes such as DNS, API Gateway, and load balancers.

## Example Usage

```terraform
data "simfra_docker_containers" "all" {}

output "container_names" {
  value = [for c in data.simfra_docker_containers.all.containers : c.name]
}
```

## Schema

### Read-Only

- `containers` (List of Object) -- List of Docker containers. Each element has the following attributes:
  - `id` (String) -- Docker container ID.
  - `name` (String) -- Container name.
  - `image` (String) -- Docker image name.
  - `state` (String) -- Container state (e.g. `running`, `exited`).
  - `status` (String) -- Human-readable container status.
  - `created` (String) -- Container creation timestamp.
