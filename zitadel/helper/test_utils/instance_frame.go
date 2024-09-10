package test_utils

import (
	"context"
	"testing"

	"github.com/zitadel/zitadel-go/v3/pkg/client/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/acceptance"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

type InstanceTestFrame struct {
	BaseTestFrame
	*admin.Client
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
	return &InstanceTestFrame{
		BaseTestFrame: *baseFrame,
		Client:        adminClient,
	}
}
