---
page_title: "Provider: Simfra"
description: |-
  The Simfra provider manages accounts and AWS resources in a Simfra cloud emulator.
---

# Simfra Provider

The Simfra provider is used to manage accounts and pre-seed AWS resources in a [Simfra](https://simfra.dev) cloud emulator instance. Simfra is a single-binary AWS cloud emulator that provides a local, high-fidelity simulation of AWS services for development and testing.

The provider offers two categories of functionality:

- **Admin API** -- Create and manage Simfra accounts, query server health, inspect Docker containers and networks, manage port forwards, and retrieve infrastructure metadata. These operations use the Simfra admin API and require only the `endpoint` (and optionally `admin_token`).

- **AWS resource seeding** -- Create AWS resources with specific, deterministic IDs (e.g., Route53 hosted zones with a known zone ID, Organizations with a known org ID). These operations use the standard AWS SDK against the Simfra endpoint and require `access_key`, `secret_key`, and `region` in addition to the `endpoint`.

## Example Usage

### Admin API only

```terraform
provider "simfra" {
  endpoint    = "http://localhost:4599"
  admin_token = "my-admin-token"
}

resource "simfra_account" "dev" {
  account_id       = "123456789012"
  bootstrap        = "standard"
  bootstrap_region = "us-east-1"
}
```

### Admin API + AWS resource seeding

```terraform
provider "simfra" {
  endpoint       = "http://localhost:4599"
  admin_token    = "my-admin-token"
  skip_tls_verify = true

  access_key = "AKIAIOSFODNN7EXAMPLE"
  secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  region     = "us-east-1"
}

resource "simfra_account" "dev" {
  account_id       = "123456789012"
  bootstrap        = "standard"
  bootstrap_region = "us-east-1"
}

resource "simfra_route53_zone" "example" {
  id      = "Z0123456789ABCDEFGHIJ"
  name    = "example.com"
  comment = "Production domain"
}
```

## Authentication

The provider supports two authentication mechanisms, which can be used independently or together:

### Admin API authentication

The admin API is used for account management and server introspection. Configure it with `endpoint` and optionally `admin_token`:

- `endpoint` -- The Simfra server URL. Can also be set via the `SIMFRA_ENDPOINT` environment variable.
- `admin_token` -- A bearer token for the admin API. Can also be set via the `SIMFRA_ADMIN_TOKEN` environment variable. Optional if the server does not require authentication.

### AWS API authentication

AWS API authentication is required for resources that create AWS objects (Route53 zones, Organizations). Configure it with `access_key`, `secret_key`, and `region`:

- `access_key` -- An AWS access key ID for the Simfra account. Can also be set via the `AWS_ACCESS_KEY_ID` environment variable.
- `secret_key` -- The corresponding secret access key. Can also be set via the `AWS_SECRET_ACCESS_KEY` environment variable.
- `region` -- The AWS region. Can also be set via `AWS_REGION` or `AWS_DEFAULT_REGION`. Defaults to `us-east-1` when credentials are provided but region is not set.

## Schema

### Optional

- `endpoint` (String) -- Simfra server endpoint URL. Can also be set with the `SIMFRA_ENDPOINT` environment variable.
- `admin_token` (String, Sensitive) -- Admin API bearer token. Can also be set with the `SIMFRA_ADMIN_TOKEN` environment variable.
- `skip_tls_verify` (Boolean) -- Skip TLS certificate verification. Can also be set with the `SIMFRA_SKIP_TLS_VERIFY` environment variable.
- `access_key` (String) -- AWS access key ID for creating AWS resources. Can also be set with the `AWS_ACCESS_KEY_ID` environment variable.
- `secret_key` (String, Sensitive) -- AWS secret access key for creating AWS resources. Can also be set with the `AWS_SECRET_ACCESS_KEY` environment variable.
- `region` (String) -- AWS region for creating AWS resources. Can also be set with the `AWS_REGION` environment variable.
