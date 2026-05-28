---
page_title: "simfra_services Data Source - terraform-provider-simfra"
subcategory: ""
description: |-
  Lists all available Simfra services.
---

# simfra_services (Data Source)

Lists all AWS services available in the Simfra instance, including their supported protocols and operations.

## Example Usage

```terraform
data "simfra_services" "all" {}

output "service_names" {
  value = [for s in data.simfra_services.all.services : s.name]
}
```

## Schema

### Read-Only

- `services` (List of Object) -- List of services. Each element has the following attributes:
  - `name` (String) -- Service name (e.g. `sqs`, `s3`, `ec2`).
  - `description` (String) -- Human-readable description of the service.
  - `protocols` (List of String) -- Supported wire protocols (e.g. `query`, `json`, `rest-xml`).
  - `operations` (List of String) -- Supported API operations.
