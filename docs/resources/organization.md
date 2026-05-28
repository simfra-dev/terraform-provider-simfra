---
page_title: "simfra_organization Resource - terraform-provider-simfra"
subcategory: ""
description: |-
  Creates an AWS Organization with optional specific org and root IDs.
---

# simfra_organization (Resource)

Creates an AWS Organization in the Simfra account. This resource allows you to pre-seed Organizations with deterministic organization and root IDs, which is useful for matching production infrastructure IDs in local development and testing.

This resource requires AWS credentials (`access_key` and `secret_key`) to be configured on the provider.

## Example Usage

```terraform
# Create an AWS Organization with specific IDs
resource "simfra_organization" "example" {
  id          = "o-abc1234567"
  root_id     = "r-ab12"
  feature_set = "ALL"
}

# Create an AWS Organization with auto-generated IDs
resource "simfra_organization" "auto" {
  feature_set = "ALL"
}

output "organization_arn" {
  value = simfra_organization.example.arn
}

output "master_account_id" {
  value = simfra_organization.example.master_account_id
}
```

## Schema

### Optional

- `id` (String) -- Organization ID (e.g. `o-abc1234567`). If specified, the organization is created with this exact ID. If omitted, Simfra generates one. Changing this forces a new resource.
- `root_id` (String) -- Organization root ID (e.g. `r-ab12`). If specified, the root is created with this exact ID. If omitted, Simfra generates one. Changing this forces a new resource.
- `feature_set` (String) -- Feature set: `ALL` or `CONSOLIDATED_BILLING`. Defaults to `ALL`. Changing this forces a new resource.

### Read-Only

- `arn` (String) -- Organization ARN.
- `master_account_id` (String) -- Management account ID.
- `master_account_arn` (String) -- Management account ARN.
