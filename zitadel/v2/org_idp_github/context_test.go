package org_idp_github_test

import (
	"context"
	"fmt"
	"os"

	"github.com/zitadel/terraform-provider-zitadel/zitadel"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	mgmt "github.com/zitadel/zitadel-go/v2/pkg/client/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const zitadelCtxKey = "zitadel"

type zitadelContext struct {
	client *mgmt.Client
	orgID,
	terraformType, terraformID, terraformName,
	providerSnippet string
	zitadelProvider *schema.Provider
}

func fromZitadelContext(ctx context.Context) *zitadelContext {
	return ctx.Value(zitadelCtxKey).(*zitadelContext)
}

func toZitadelContext() (context.Context, error) {
	const (
		orgName  = "terraform-tests"
		domain   = "localhost"
		insecure = true
		port     = "8080"
	)
	ctx := context.Background()
	tokenPath := os.Getenv("TF_ACC_ZITADEL_TOKEN")
	zitadelProvider := zitadel.Provider()
	diag := zitadelProvider.Configure(ctx, terraform.NewResourceConfigRaw(map[string]interface{}{
		"domain":   domain,
		"insecure": insecure,
		"port":     port,
		"token":    tokenPath,
	}))
	providerSnippet := fmt.Sprintf(`
provider "zitadel" {
  domain   = "%s"
  insecure = "%t"
  port     = "%s"
  token    = "%s"
}
`, domain, insecure, port, tokenPath)
	if diag.HasError() {
		return nil, fmt.Errorf("unknown error configuring the test provider: %v", diag)
	}
	clientInfo := zitadelProvider.Meta().(*helper.ClientInfo)
	mgmtClient, err := helper.GetManagementClient(clientInfo, "")
	if err != nil {
		return nil, err
	}
	org, err := mgmtClient.GetOrgByDomainGlobal(ctx, &management.GetOrgByDomainGlobalRequest{Domain: fmt.Sprintf("%s.%s", orgName, domain)})
	orgID := org.GetOrg().GetId()
	if status.Code(err) == codes.NotFound {
		var newOrg *management.AddOrgResponse
		newOrg, err = mgmtClient.AddOrg(ctx, &management.AddOrgRequest{Name: orgName})
		orgID = newOrg.GetId()
	}
	if err != nil {
		return nil, err
	}
	mgmtClient, err = helper.GetManagementClient(clientInfo, orgID)
	if err != nil {
		return nil, err
	}
	terraformType := "zitadel_org_idp_github"
	terraformID := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	terraformName := fmt.Sprintf("%s.%s", terraformType, terraformID)
	return context.WithValue(ctx, zitadelCtxKey, &zitadelContext{
		client:          mgmtClient,
		orgID:           orgID,
		terraformType:   terraformType,
		terraformID:     terraformID,
		terraformName:   terraformName,
		providerSnippet: providerSnippet,
		zitadelProvider: zitadelProvider,
	}), err
}
