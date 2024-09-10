package project_test_dep

import (
	"testing"

	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project"
)

func Create(t *testing.T, frame *test_utils.OrgTestFrame, name string) (string, string) {
	return test_utils.CreateDefaultDependency(t,
		"zitadel_project",
		project.ProjectIDVar,
		func() (string, error) {
			p, err := frame.AddProject(frame, &management.AddProjectRequest{Name: name})
			return p.GetId(), err
		})
}
