terraform {
  required_providers {
    simfra = {
      source = "simfra-dev/simfra"
    }
  }
}

provider "simfra" {
  endpoint = "http://localhost:4599"

  # Optional: admin API token
  # admin_token = "my-admin-token"

  # Required for simfra_route53_zone and simfra_organization resources.
  # Can also be set via AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_REGION.
  access_key = "AKIAIOSFODNN7EXAMPLE"
  secret_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  region     = "us-east-1"
}
