package zitadel

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestBindingResourceSchemas(t *testing.T) {
	p := Provider()
	if err := p.InternalValidate(); err != nil {
		t.Fatal(err)
	}
	r := p.ResourcesMap["zitadel_administrator"]
	for name, config := range map[string]map[string]interface{}{
		"instance with org": {"user_id": "u", "resource_type": "instance", "org_id": "o", "roles": []interface{}{"IAM_OWNER"}},
		"missing project":   {"user_id": "u", "resource_type": "project", "roles": []interface{}{"PROJECT_OWNER"}},
		"missing grantee":   {"user_id": "u", "resource_type": "project_grant", "project_id": "p", "roles": []interface{}{"PROJECT_GRANT_OWNER"}},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := r.Diff(context.Background(), &terraform.InstanceState{}, terraform.NewResourceConfigRaw(config), nil); err == nil {
				t.Fatal("invalid administrator scope passed plan validation")
			}
		})
	}
	for name, config := range map[string]map[string]interface{}{
		"explicit instance": {"user_id": "u", "resource_type": "instance", "roles": []interface{}{"IAM_OWNER"}},
		"org":               {"user_id": "u", "resource_type": "organization", "org_id": "o", "roles": []interface{}{"ORG_USER_MANAGER"}},
		"project":           {"user_id": "u", "resource_type": "project", "project_id": "p", "roles": []interface{}{"PROJECT_OWNER"}},
		"project grant":     {"user_id": "u", "resource_type": "project_grant", "project_id": "p", "org_id": "o", "roles": []interface{}{"PROJECT_GRANT_OWNER"}},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := r.Diff(context.Background(), &terraform.InstanceState{}, terraform.NewResourceConfigRaw(config), nil); err != nil {
				t.Fatal(err)
			}
		})
	}
	for _, resource := range []string{"zitadel_authorization", "zitadel_administrator"} {
		for key, f := range p.ResourcesMap[resource].Schema {
			if key == "role_keys" || key == "roles" || key == "state" {
				continue
			}
			if f.Type != schema.TypeString || !f.ForceNew {
				t.Errorf("%s.%s identity is not immutable", resource, key)
			}
		}
	}
}
