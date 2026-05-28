---
page_title: "simfra_ca_info Data Source - terraform-provider-simfra"
subcategory: ""
description: |-
  Retrieves Simfra certificate authority information.
---

# simfra_ca_info (Data Source)

Retrieves information about the Simfra certificate authority (CA), including the root and intermediate CA certificate details. The root CA is used to sign TLS certificates for Simfra services, and clients trust it to verify connections.

## Example Usage

```terraform
data "simfra_ca_info" "info" {}

output "root_subject" {
  value = data.simfra_ca_info.info.root_subject
}

output "root_expires" {
  value = data.simfra_ca_info.info.root_not_after
}
```

## Schema

### Read-Only

- `persistent` (Boolean) -- Whether the CA is persisted to disk.
- `directory` (String) -- Directory where CA files are stored.
- `root_subject` (String) -- Root CA certificate subject.
- `root_serial_number` (String) -- Root CA certificate serial number.
- `root_not_after` (String) -- Root CA certificate expiration date.
- `intermediate_subject` (String) -- Intermediate CA certificate subject.
- `intermediate_serial_number` (String) -- Intermediate CA certificate serial number.
- `intermediate_not_after` (String) -- Intermediate CA certificate expiration date.
