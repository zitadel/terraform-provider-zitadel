resource "zitadel_action_target" "default" {
	name               = "response-enricher"
	endpoint           = "https://example.com/api/enrich"
	target_type        = "REST_CALL"
	timeout            = "10s"
	interrupt_on_error = false
}

resource "zitadel_action_execution_response" "default" {
	method = "/zitadel.user.v2.UserService/GetUser"
	target_ids = [zitadel_action_target.default.id]
}
