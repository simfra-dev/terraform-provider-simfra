---
page_title: "simfra_accounts Data Source - terraform-provider-simfra"
subcategory: ""
description: |-
  Lists all Simfra accounts.
---

# simfra_accounts (Data Source)

Lists all accounts registered in the Simfra instance.

## Example Usage

```terraform
data "simfra_accounts" "all" {}

output "account_ids" {
  value = [for a in data.simfra_accounts.all.accounts : a.account_id]
}
```

## Schema

### Read-Only

- `accounts` (List of Object) -- List of accounts. Each element has the following attributes:
  - `account_id` (String) -- The 12-digit AWS account ID.
  - `alias` (String) -- Account alias, if set.
  - `created_at` (String) -- Account creation timestamp.
