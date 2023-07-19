package test_utils

import (
	"context"

	"github.com/zitadel/terraform-provider-zitadel/acceptance"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	mgmt "github.com/zitadel/zitadel-go/v2/pkg/client/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

const (
	orgName = "terraform-tests"
)

type OrgTestFrame struct {
	BaseTestFrame
	*mgmt.Client
	OrgID string
}

func NewOrgTestFrame(resourceType string) (*OrgTestFrame, error) {
	ctx := context.Background()
	cfg := acceptance.GetConfig().OrgLevel
	baseFrame, err := NewBaseTestFrame(ctx, resourceType, cfg.Domain, cfg.AdminSAJSON)
	if err != nil {
		return nil, err
	}
	mgmtClient, err := helper.GetManagementClient(baseFrame.ClientInfo, "")
	if err != nil {
		return nil, err
	}
	org, err := mgmtClient.GetOrgByDomainGlobal(baseFrame, &management.GetOrgByDomainGlobalRequest{Domain: "zitadel." + cfg.Domain})
	if err != nil {
		return nil, err
	}
	orgID := org.GetOrg().GetId()
	mgmtClient, err = helper.GetManagementClient(baseFrame.ClientInfo, orgID)
	return &OrgTestFrame{
		BaseTestFrame: *baseFrame,
		Client:        mgmtClient,
		OrgID:         orgID,
	}, err
}
