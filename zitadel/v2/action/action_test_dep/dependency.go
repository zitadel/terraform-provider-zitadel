package action_test_dep

import (
	"testing"
	"time"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/action"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func Create(t *testing.T, frame *test_utils.OrgTestFrame) (string, string) {
	return test_utils.CreateDefaultDependency(t, "zitadel_action", action.ActionIDVar, func() (string, error) {
		a, err := frame.CreateAction(frame, &management.CreateActionRequest{
			Name:          frame.UniqueResourcesID,
			Script:        "not a script",
			Timeout:       durationpb.New(10 * time.Second),
			AllowedToFail: true,
		})
		return a.GetId(), err
	})
}
