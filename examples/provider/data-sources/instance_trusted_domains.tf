data "zitadel_instance_trusted_domains" "default" {
  # instance_id is optional - uses current context if not provided
}

output "trusted_domains" {
  value = data.zitadel_instance_trusted_domains.default.domains
}
