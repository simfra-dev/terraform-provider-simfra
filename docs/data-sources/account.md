---
page_title: "simfra_account Data Source - terraform-provider-simfra"
subcategory: ""
description: |-
  Retrieves details of a Simfra account.
---

# simfra_account (Data Source)

Retrieves details of a Simfra account by account ID, including the root user credentials.

## Example Usage

```terraform
data "simfra_account" "dev" {
  account_id = "123456789012"
}

output "root_access_key_id" {
  value = data.simfra_account.dev.root_access_key_id
}
```

## Schema

### Required

- `account_id` (String) -- The 12-digit AWS account ID.

### Read-Only

- `root_access_key_id` (String, Sensitive) -- Root user access key ID.
- `root_secret_access_key` (String, Sensitive) -- Root user secret access key.
- `created_at` (String) -- Account creation timestamp.
