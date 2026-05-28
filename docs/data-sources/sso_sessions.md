---
page_title: "simfra_sso_sessions Data Source - terraform-provider-simfra"
subcategory: ""
description: |-
  Lists all active SSO sessions.
---

# simfra_sso_sessions (Data Source)

Lists all active SSO (Single Sign-On) sessions in the Simfra instance.

## Example Usage

```terraform
data "simfra_sso_sessions" "all" {}

output "active_sessions" {
  value = [for s in data.simfra_sso_sessions.all.sessions : {
    user_name = s.user_name
    expired   = s.expired
  }]
}
```

## Schema

### Read-Only

- `sessions` (List of Object) -- List of SSO sessions. Each element has the following attributes:
  - `token` (String, Sensitive) -- Session token.
  - `user_id` (String) -- User ID.
  - `user_name` (String) -- User name.
  - `expires_at` (String) -- Session expiration timestamp.
  - `created_at` (String) -- Session creation timestamp.
  - `expired` (Boolean) -- Whether the session has expired.
