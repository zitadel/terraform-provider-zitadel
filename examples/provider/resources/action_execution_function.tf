resource "zitadel_action_target" "default" {
	name               = "token-enricher"
	endpoint           = "https://example.com/oidc/enrich"
	target_type        = "REST_CALL"
	timeout            = "10s"
	interrupt_on_error = true
}

resource "zitadel_action_execution_function" "default" {
	name = "preaccesstoken"
	target_ids = [zitadel_action_target.default.id]
}
