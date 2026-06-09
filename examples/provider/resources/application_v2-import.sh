# The resource can be imported using the ID format `<id[:org_id]>`, e.g.
# project_id is refreshed from the v2 API on the first read and does not
# need to be passed at import time.
terraform import zitadel_application_v2.imported '123456789012345678:123456789012345678'
