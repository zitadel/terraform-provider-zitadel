terraform {
  required_providers {
    zitadel = {
      source  = "zitadel/zitadel"
      version = "0.0.0"
    }
  }
}

provider zitadel {
  issuer  = "http://localhost:8080/oauth/v2"
  address = "localhost:8080"
  project = "160549024225689888"
  token   = "/Users/benz/go/src/github.com/zitadel/terraform-provider-zitadel/local-token"
}


data zitadelV1Org zitadelV1Org {
  provider = zitadel
  issuer   = "https://issuer.zitadel.dev"
  address  = "api.zitadel.dev:443"
  project  = "70669147545070419"
  token    = "/Users/benz/go/src/github.com/zitadel/terraform-provider-zitadel/zitadel-dev-token"
}

output fetched_org_id {
  value = data.zitadelV1Org.zitadelV1Org.org
}

output fetched_org_name {
  value = data.zitadelV1Org.zitadelV1Org.name
}

resource org org {
  provider = zitadel
  old_id   = data.zitadelV1Org.zitadelV1Org.org
  name     = data.zitadelV1Org.zitadelV1Org.name
}

resource user userTest {
  depends_on = [data.zitadelV1Org.zitadelV1Org, org.org]
  provider   = zitadel

  for_each = {
  for idx, user in data.zitadelV1Org.zitadelV1Org.users : user.user_name => user
  }
  old_id               = each.value.id
  resource_owner       = org.org.id
  state                = each.value.state
  user_name            = each.value.user_name
  login_names          = each.value.login_names
  preferred_login_name = each.value.preferred_login_name
  type                 = each.value.type
  first_name           = each.value.first_name
  last_name            = each.value.last_name
  nick_name            = each.value.nick_name
  display_name         = each.value.display_name
  preferred_language   = each.value.preferred_language
  gender               = each.value.gender
  phone                = each.value.phone
  is_phone_verified    = each.value.is_phone_verified
  email                = each.value.email
  is_email_verified    = each.value.is_email_verified
  name                 = each.value.name
  description          = each.value.description
}


resource project projectTest {
  depends_on = [data.zitadelV1Org.zitadelV1Org, org.org]
  provider   = zitadel

  for_each = {
  for idx, project in data.zitadelV1Org.zitadelV1Org.projects : project.name => project
  }
  old_id                   = each.value.id
  name                     = each.value.name
  state                    = each.value.state
  resource_owner           = org.org.id
  project_role_assertion   = each.value.project_role_assertion
  project_role_check       = each.value.project_role_check
  has_project_check        = each.value.has_project_check
  private_labeling_setting = each.value.private_labeling_setting
}


resource domain domainTest {
  depends_on = [data.zitadelV1Org.zitadelV1Org, org.org]
  provider   = zitadel

  for_each = {
  for idx, domain in data.zitadelV1Org.zitadelV1Org.domains : domain.name => domain
  }

  name                     = each.value.name
  org_id                    = org.org.id
}

