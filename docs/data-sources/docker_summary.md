---
page_title: "simfra_docker_summary Data Source - terraform-provider-simfra"
subcategory: ""
description: |-
  Retrieves Docker resource summary from Simfra.
---

# simfra_docker_summary (Data Source)

Retrieves a summary of Docker resources managed by the Simfra instance, including counts of containers, images, networks, and volumes.

## Example Usage

```terraform
data "simfra_docker_summary" "info" {}

output "containers_running" {
  value = data.simfra_docker_summary.info.containers_running
}

output "images" {
  value = data.simfra_docker_summary.info.images
}
```

## Schema

### Read-Only

- `containers_total` (Number) -- Total number of Docker containers.
- `containers_running` (Number) -- Number of running Docker containers.
- `containers_stopped` (Number) -- Number of stopped Docker containers.
- `images` (Number) -- Number of Docker images.
- `networks` (Number) -- Number of Docker networks.
- `volumes` (Number) -- Number of Docker volumes.
