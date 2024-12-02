package human_user_test_dep

import (
	"testing"

	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/human_user"
)

func Create(t *testing.T, frame *test_utils.OrgTestFrame) (string, string) {
	return test_utils.CreateDefaultDependency(t, "zitadel_human_user", human_user.UserIDVar, func() (string, error) {
		user, err := frame.ImportHumanUser(frame, &management.ImportHumanUserRequest{
			UserName: frame.UniqueResourcesID,
			Profile: &management.ImportHumanUserRequest_Profile{
				FirstName: "Don't",
				LastName:  "Care",
			},
			Email: &management.ImportHumanUserRequest_Email{
				Email:           "dont@care.com",
				IsEmailVerified: true,
			},
		})
		return user.GetUserId(), err
	})
}
