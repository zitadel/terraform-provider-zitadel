# The resource can be imported using the ID format `<id[:org_id[:client_secret]]>`, e.g.
# project_id is refreshed from the v2 API on the first read and does not need to
# be passed at import time.
terraform import zitadel_application_v2.imported '123456789012345678:123456789012345678'

# For OIDC/API applications with a generated client secret, pass the secret as
# the optional final segment so it is preserved in state (it is never returned
# by the read API). See the "Migrating to the v2 resources" guide for details.
terraform import zitadel_application_v2.imported '123456789012345678:123456789012345678:THE_CLIENT_SECRET'
