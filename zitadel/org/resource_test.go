package org_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	adminpb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	orgpb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/org"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org"
)

func TestAccOrg(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, org.NameVar, exampleAttributes).AsString()
	initialProperty := "initialorgname_" + frame.UniqueResourcesID
	updatedProperty := "updatedorgname_" + frame.UniqueResourcesID
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		initialProperty, updatedProperty,
		"", "", "",
		false,
		checkRemoteProperty(frame, idFromState(frame)),
		helper.ZitadelGeneratedIdOnlyRegex,
		orgGone(frame, idFromState(frame)),
		test_utils.ImportResourceId(frame.BaseTestFrame),
	)
}

func idFromState(frame *test_utils.OrgTestFrame) func(*terraform.State) string {
	return func(state *terraform.State) string {
		return frame.State(state).ID
	}
}

// orgGone verifies the org is effectively gone after destroy by looking it up by ID.
// Treat NOT_FOUND, PERMISSION_DENIED, and UNAUTHENTICATED as successful "gone" states.
func orgGone(
	frame *test_utils.OrgTestFrame,
	idFromState func(*terraform.State) string,
) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		orgID := idFromState(state)

		ctx := helper.CtxSetOrgID(context.Background(), orgID)

		adm, err := helper.GetAdminClient(ctx, frame.ClientInfo)
		if err != nil {
			return fmt.Errorf("get admin client: %w", err)
		}

		resp, err := adm.GetOrgByID(ctx, &adminpb.GetOrgByIDRequest{Id: orgID})
		if err != nil {
			st, _ := status.FromError(err)
			switch st.Code() {
			case codes.NotFound, codes.PermissionDenied, codes.Unauthenticated:
				return nil
			}
			// Tolerate gateway/textual variants
			msg := strings.ToLower(err.Error())
			if strings.Contains(msg, "not found") ||
				strings.Contains(msg, "permission denied") ||
				strings.Contains(msg, "unauthenticated") {
				return nil
			}
			return fmt.Errorf("unexpected error after destroy (id=%s): %v", orgID, err)
		}

		// Call succeeded: consider it "gone" if it's not ACTIVE (e.g., soft-deleted).
		if resp != nil && resp.Org != nil && resp.Org.State != orgpb.OrgState_ORG_STATE_ACTIVE {
			return nil
		}

		return fmt.Errorf("expected org to be gone, but it's still ACTIVE (id=%s)", orgID)
	}
}
