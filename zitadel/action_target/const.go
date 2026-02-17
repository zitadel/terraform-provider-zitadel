package action_target

const (
	TargetIDVar         = "target_id"
	NameVar             = "name"
	EndpointVar         = "endpoint"
	TargetTypeVar       = "target_type"
	TimeoutVar          = "timeout"
	InterruptOnErrorVar = "interrupt_on_error"
	SigningKeyVar       = "signing_key"
	PayloadTypeVar      = "payload_type"

	targetTypeRestWebhook = "REST_WEBHOOK"
	targetTypeRestCall    = "REST_CALL"
	targetTypeRestAsync   = "REST_ASYNC"

	payloadTypeJSON = "PAYLOAD_TYPE_JSON"
	payloadTypeJWT  = "PAYLOAD_TYPE_JWT"
	payloadTypeJWE  = "PAYLOAD_TYPE_JWE"
)
