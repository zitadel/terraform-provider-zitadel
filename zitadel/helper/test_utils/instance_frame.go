package test_utils

import (
	"context"
	"testing"

	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	instanceV2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/instance/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/acceptance"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

type InstanceTestFrame struct {
	BaseTestFrame
	AdminClient admin.AdminServiceClient
	InstanceID  string
}

func NewInstanceTestFrame(t *testing.T, resourceType string) *InstanceTestFrame {
	ctx := context.Background()
	cfg := acceptance.GetConfig().InstanceLevel
	baseFrame, err := NewBaseTestFrame(ctx, resourceType, cfg.Domain, cfg.AdminSAJSON)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	adminClient, err := helper.GetAdminClient(baseFrame.Context, baseFrame.ClientInfo)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}

	instanceClient, err := helper.GetInstanceClient(baseFrame.Context, baseFrame.ClientInfo)
	if err != nil {
		t.Fatalf("failed to get instance client: %v", err)
	}

	instanceResp, err := instanceClient.GetInstance(baseFrame.Context, &instanceV2.GetInstanceRequest{})
	if err != nil {
		t.Fatalf("failed to get instance: %v", err)
	}

	return &InstanceTestFrame{
		BaseTestFrame: *baseFrame,
		AdminClient:   adminClient,
		InstanceID:    instanceResp.Instance.Id,
	}
}
