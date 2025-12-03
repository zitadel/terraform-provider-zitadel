data "zitadel_instance" "default" {
	# instance_id is optional - uses current context if not provided
}

output "instance" {
	value = data.zitadel_instance.default
}

output "primary_domain" {
	value = data.zitadel_instance.default.primary_domain
}
