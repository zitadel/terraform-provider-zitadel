# The resource can be imported using the ID format `<instance_id>/<domain>` or just `<domain>`, e.g.
terraform import zitadel_instance_trusted_domain.imported '123456789012345678/idp.partner.com'
# Or if using instance context:
terraform import zitadel_instance_trusted_domain.imported 'idp.partner.com'
