package machine_user_test_dep

import (
	"testing"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/machine_user"
)

func Create(t *testing.T, frame *test_utils.OrgTestFrame) (string, string) {
	return test_utils.CreateDefaultDependency(t, "zitadel_machine_user", machine_user.UserIDVar, func() (string, error) {
		user, err := frame.AddMachineUser(frame, &management.AddMachineUserRequest{
			UserName: frame.UniqueResourcesID,
			Name:     "Don't care",
		})
		userID := user.GetUserId()
		return userID, err
	})
}
