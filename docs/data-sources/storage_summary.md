---
page_title: "simfra_storage_summary Data Source - terraform-provider-simfra"
subcategory: ""
description: |-
  Retrieves Simfra storage and persistence summary.
---

# simfra_storage_summary (Data Source)

Retrieves information about the Simfra server's storage and persistence configuration, including whether persistence is enabled, the database path and size, and per-service storage usage.

## Example Usage

```terraform
data "simfra_storage_summary" "info" {}

output "persistence_enabled" {
  value = data.simfra_storage_summary.info.enabled
}

output "db_size" {
  value = data.simfra_storage_summary.info.db_size
}
```

## Schema

### Read-Only

- `enabled` (Boolean) -- Whether persistence is enabled.
- `data_dir` (String) -- Path to the data directory.
- `db_path` (String) -- Path to the SQLite database file.
- `db_size` (Number) -- Size of the database file in bytes.
- `total` (Number) -- Total storage usage in bytes.
- `services` (List of Object) -- Per-service storage usage. Each element has the following attributes:
  - `service` (String) -- Service name.
  - `size` (Number) -- Storage usage in bytes.
