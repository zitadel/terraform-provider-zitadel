# The resource can be imported using the ID format `<grant_id:user_id[:org_id]>`, e.g.
terraform import zitadel_user_grant.imported '123456789012345678:123456789012345678:123456789012345678'
# The grant ID is not shown in the ZITADEL console. It can be looked up per user with the Management API
# ListUserGrants endpoint (https://zitadel.com/docs/apis/resources/mgmt/management-service-list-user-grants),
# where each result carries its `id`.
