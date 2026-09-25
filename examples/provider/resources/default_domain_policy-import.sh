# Terraform rejects empty import IDs, so the resource is imported with the placeholder id
# `default`, which ZITADEL ignores, e.g.
terraform import zitadel_default_domain_policy.imported 'default'
