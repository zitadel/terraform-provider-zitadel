package test_utils

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/zitadel/zitadel-go/v3/pkg/client/admin"
	mgmt "github.com/zitadel/zitadel-go/v3/pkg/client/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/acceptance"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
)

type OrgTestFrame struct {
	BaseTestFrame
	*mgmt.Client
	Admin                  *admin.Client
	OrgID                  string
	AsOrgDefaultDependency string
}

func (o *OrgTestFrame) useOrgContext(orgID string) (err error) {
	o.Client, err = helper.GetManagementClient(o.Context, o.BaseTestFrame.ClientInfo)
	if err != nil {
		return err
	}
	o.Context = helper.CtxSetOrgID(o.Context, orgID)
	o.Admin, err = helper.GetAdminClient(o.Context, o.BaseTestFrame.ClientInfo)
	o.AsOrgDefaultDependency = strings.Replace(o.AsOrgDefaultDependency, o.OrgID, orgID, 1)
	o.OrgID = orgID
	return err
}

func NewOrgTestFrame(t *testing.T, resourceType string) *OrgTestFrame {
	ctx := context.Background()
	cfg := acceptance.GetConfig().OrgLevel
	baseFrame, err := NewBaseTestFrame(ctx, resourceType, cfg.Domain, cfg.AdminSAJSON)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	orgFrame := &OrgTestFrame{
		BaseTestFrame: *baseFrame,
	}
	if err = orgFrame.useOrgContext(""); err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	org, err := orgFrame.GetOrgByDomainGlobal(baseFrame, &management.GetOrgByDomainGlobalRequest{Domain: "zitadel." + cfg.Domain})
	if err != nil {
		t.Fatalf("failed to get org by domain: %v", err)
	}
	orgFrame.OrgID = org.GetOrg().GetId()
	orgFrame.AsOrgDefaultDependency = fmt.Sprintf(`
data "zitadel_org" "default" {
	id = "%s"
}
`, orgFrame.OrgID)
	return orgFrame
}

func (o OrgTestFrame) AnotherOrg(t *testing.T, name string) *OrgTestFrame {
	org, err := o.Client.AddOrg(o, &management.AddOrgRequest{
		Name: name,
	})
	if err != nil {
		t.Fatalf("failed to create org: %v", err)
	}
	if err := o.useOrgContext(org.GetId()); err != nil {
		t.Fatalf("failed to use org context: %v", err)
	}
	return &o
}
