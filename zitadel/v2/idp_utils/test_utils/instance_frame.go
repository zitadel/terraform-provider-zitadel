package test_utils

import (
	"github.com/zitadel/zitadel-go/v2/pkg/client/admin"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

type InstanceTestFrame struct {
	BaseTestFrame
	*admin.Client
}

func NewInstanceTestFrame(resourceType string) (*InstanceTestFrame, error) {
	baseFrame, err := NewBaseTestFrame(resourceType)
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
	}, err
}
