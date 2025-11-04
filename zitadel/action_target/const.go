package action_target

const (
	TargetIDVar         = "target_id"
	NameVar             = "name"
	EndpointVar         = "endpoint"
	TargetTypeVar       = "target_type"
	TimeoutVar          = "timeout"
	InterruptOnErrorVar = "interrupt_on_error"
	SigningKeyVar       = "signing_key"

	targetTypeRestWebhook = "REST_WEBHOOK"
	targetTypeRestCall    = "REST_CALL"
	targetTypeRestAsync   = "REST_ASYNC"
)
