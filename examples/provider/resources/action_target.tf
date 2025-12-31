resource "zitadel_action_target" "default" {
  name               = "test-target-name"
  endpoint           = "https://example.com/endpoint"
  target_type        = "REST_WEBHOOK"
  timeout            = "15s"
  interrupt_on_error = true
}
