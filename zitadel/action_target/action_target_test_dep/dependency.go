package action_target_test_dep

import (
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_target"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	actionv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"
)

func Create(t *testing.T, frame *test_utils.OrgTestFrame) (string, string) {
	return test_utils.CreateDefaultDependency(t, "zitadel_action_target", action_target.TargetIDVar, func() (string, error) {
		client, err := helper.GetActionClient(frame.Context, frame.ClientInfo)
		a, err := client.CreateTarget(frame.Context, &actionv2.CreateTargetRequest{
			Name:     frame.UniqueResourcesID,
			Endpoint: "https://example.com",
			Timeout:  durationpb.New(10 * time.Second),
			TargetType: &actionv2.CreateTargetRequest_RestWebhook{
				RestWebhook: &actionv2.RESTWebhook{InterruptOnError: false},
			},
		})
		return a.GetId(), err
	})
}
