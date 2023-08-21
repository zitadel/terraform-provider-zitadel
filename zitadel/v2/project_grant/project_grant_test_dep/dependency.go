package project_grant_test_dep

import (
	"testing"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func Create(t *testing.T, frame *test_utils.OrgTestFrame, projectID, grantedOrgID string) string {
	dep, err := frame.AddProjectGrant(frame, &management.AddProjectGrantRequest{
		ProjectId:    projectID,
		GrantedOrgId: grantedOrgID,
	})
	if err != nil {
		t.Errorf("failed to create a project grant: %v", err)
	}
	return dep.GetGrantId()
}
