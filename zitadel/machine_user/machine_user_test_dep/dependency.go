package machine_user_test_dep

import (
	"testing"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/machine_user"
)

func Create(t *testing.T, frame *test_utils.OrgTestFrame, username string) (string, string) {
	return test_utils.CreateOrgDefaultDependency(t, "zitadel_machine_user", frame.OrgID, machine_user.UserIDVar, func() (string, error) {
		user, err := frame.AddMachineUser(frame, &management.AddMachineUserRequest{
			UserName: username,
			Name:     "Don't care",
		})
		userID := user.GetUserId()
		return userID, err
	})
}
