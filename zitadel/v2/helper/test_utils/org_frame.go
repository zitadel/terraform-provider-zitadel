package test_utils

import (
	"fmt"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	mgmt "github.com/zitadel/zitadel-go/v2/pkg/client/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	baseFrame, err := NewBaseTestFrame(resourceType)
	if err != nil {
		return nil, err
	}
	mgmtClient, err := helper.GetManagementClient(baseFrame.ClientInfo, "")
	if err != nil {
		return nil, err
	}
	org, err := mgmtClient.GetOrgByDomainGlobal(baseFrame, &management.GetOrgByDomainGlobalRequest{Domain: fmt.Sprintf("%s.%s", orgName, domain)})
	orgID := org.GetOrg().GetId()
	if status.Code(err) == codes.NotFound {
		var newOrg *management.AddOrgResponse
		newOrg, err = mgmtClient.AddOrg(baseFrame, &management.AddOrgRequest{Name: orgName})
		orgID = newOrg.GetId()
	}
	if err != nil {
		return nil, err
	}
	mgmtClient, err = helper.GetManagementClient(baseFrame.ClientInfo, orgID)
	return &OrgTestFrame{
		BaseTestFrame: *baseFrame,
		Client:        mgmtClient,
		OrgID:         orgID,
	}, err
}
