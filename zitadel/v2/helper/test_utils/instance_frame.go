package test_utils

import (
	"context"

	"github.com/zitadel/terraform-provider-zitadel/acceptance"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/zitadel-go/v2/pkg/client/admin"
)

type InstanceTestFrame struct {
	BaseTestFrame
	*admin.Client
}

func NewInstanceTestFrame(resourceType string) (*InstanceTestFrame, error) {
	ctx := context.Background()
	cfg := acceptance.GetConfig().InstanceLevel
	baseFrame, err := NewBaseTestFrame(ctx, resourceType, cfg.Domain, cfg.AdminSAJSON)
	if err != nil {
		return nil, err
	}
	adminClient, err := helper.GetAdminClient(baseFrame.ClientInfo)
	if err != nil {
		return nil, err
	}
	return &InstanceTestFrame{
		BaseTestFrame: *baseFrame,
		Client:        adminClient,
	}, nil
}
