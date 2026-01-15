# method:/path for specific method, service:name for service, all for all requests
terraform import zitadel_action_execution_request.imported 'method:/zitadel.session.v2.SessionService/ListSessions'
terraform import zitadel_action_execution_request.imported 'service:zitadel.session.v2.SessionService'
terraform import zitadel_action_execution_request.imported 'all'