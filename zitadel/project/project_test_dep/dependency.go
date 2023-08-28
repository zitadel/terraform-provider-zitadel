package project_test_dep

import (
	"testing"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/project"
)

func Create(t *testing.T, frame *test_utils.OrgTestFrame) (string, string) {
	return test_utils.CreateDefaultDependency(t, "zitadel_project", project.ProjectIDVar, func() (string, error) {
		p, err := frame.AddProject(frame, &management.AddProjectRequest{Name: frame.UniqueResourcesID})
		return p.GetId(), err
	})
}
