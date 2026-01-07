resource "zitadel_action_target" "default" {
	name               = "event-webhook"
	endpoint           = "https://example.com/webhooks/events"
	target_type        = "REST_ASYNC"
	timeout            = "10s"
	interrupt_on_error = false
}

resource "zitadel_action_execution_event" "default" {
	event = "user.human.added"
	target_ids = [zitadel_action_target.default.id]
}
