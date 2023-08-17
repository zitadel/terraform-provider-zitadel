package test_utils

import (
	"context"
	"fmt"

	"github.com/zitadel/zitadel-go/v2/pkg/client/admin"
	mgmt "github.com/zitadel/zitadel-go/v2/pkg/client/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/acceptance"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

type OrgTestFrame struct {
	BaseTestFrame
	*mgmt.Client
	Admin                *admin.Client
	OrgID                string
	OrgExampleDatasource string
}

func (o *OrgTestFrame) useOrgContext(orgID string) (err error) {
	o.Client, err = helper.GetManagementClient(o.BaseTestFrame.ClientInfo, orgID)
	if err != nil {
		return err
	}
	o.Admin, err = helper.GetAdminClient(o.BaseTestFrame.ClientInfo)
	o.OrgID = orgID
	return err
}

func NewOrgTestFrame(resourceType string) (*OrgTestFrame, error) {
	ctx := context.Background()
	cfg := acceptance.GetConfig().OrgLevel
	baseFrame, err := NewBaseTestFrame(ctx, resourceType, cfg.Domain, cfg.AdminSAJSON)
	if err != nil {
		return nil, err
	}
	orgFrame := &OrgTestFrame{
		BaseTestFrame: *baseFrame,
	}
	if err = orgFrame.useOrgContext(""); err != nil {
		return nil, err
	}
	org, err := orgFrame.GetOrgByDomainGlobal(baseFrame, &management.GetOrgByDomainGlobalRequest{Domain: "zitadel." + cfg.Domain})
	orgFrame.OrgID = org.GetOrg().GetId()
	orgFrame.OrgExampleDatasource = fmt.Sprintf(`
data "zitadel_org" "org" {
	id = "%s"
}
`, orgFrame.OrgID)
	return orgFrame, err
}

func (o OrgTestFrame) AnotherOrg(name string) (*OrgTestFrame, error) {
	org, err := o.Client.AddOrg(o, &management.AddOrgRequest{
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	return &o, o.useOrgContext(org.GetId())
}
