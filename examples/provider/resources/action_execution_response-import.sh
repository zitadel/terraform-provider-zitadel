# method:/path for specific method, service:name for service, all for all responses
terraform import zitadel_action_execution_response.imported 'method:/zitadel.user.v2.UserService/GetUser'
terraform import zitadel_action_execution_response.imported 'service:zitadel.user.v2.UserService'
terraform import zitadel_action_execution_response.imported 'all'
