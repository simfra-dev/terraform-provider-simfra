data "simfra_sso_sessions" "all" {}

output "active_sessions" {
  value = [for s in data.simfra_sso_sessions.all.sessions : {
    user_name = s.user_name
    expired   = s.expired
  }]
}
