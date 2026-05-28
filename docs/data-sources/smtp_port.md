---
page_title: "simfra_smtp_port Data Source - terraform-provider-simfra"
subcategory: ""
description: |-
  Retrieves the SMTP relay port for a Simfra account.
---

# simfra_smtp_port (Data Source)

Retrieves the SMTP relay port for a specific Simfra account. The SMTP relay captures emails sent via SES and makes them available for inspection. The returned port can be used to configure local SMTP clients.

## Example Usage

```terraform
data "simfra_smtp_port" "dev" {
  account_id = "123456789012"
}

output "smtp_port" {
  value = data.simfra_smtp_port.dev.port
}
```

## Schema

### Required

- `account_id` (String) -- AWS account ID.

### Read-Only

- `port` (Number) -- SMTP relay port.
