resource "zitadel_action_target" "default" {
	name               = "request-interceptor"
	endpoint           = "https://example.com/security/validate"
	target_type        = "REST_WEBHOOK"
	timeout            = "5s"
	interrupt_on_error = true
}

resource "zitadel_action_execution_request" "default" {
	method = "/zitadel.session.v2.SessionService/ListSessions"
	target_ids = [zitadel_action_target.default.id]
}
