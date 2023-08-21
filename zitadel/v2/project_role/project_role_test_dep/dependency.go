package project_role_test_dep

import (
	"strings"
	"testing"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/project_role"
)

func Create(t *testing.T, frame *test_utils.OrgTestFrame, projectID string, key ...string) string {
	deps := make([]string, len(key))
	for i, k := range key {
		_, dep := test_utils.CreateDefaultDependency(t, "zitadel_project_role", project_role.KeyVar, func() (string, error) {
			_, err := frame.AddProjectRole(frame, &management.AddProjectRoleRequest{
				ProjectId:   projectID,
				RoleKey:     k,
				DisplayName: k,
			})
			return k, err
		})
		deps[i] = dep
	}
	return strings.Join(deps, "\n")
}
