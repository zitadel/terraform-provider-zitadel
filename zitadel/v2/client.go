package v2

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/zitadel-go/v2/pkg/client/admin"
	"github.com/zitadel/zitadel-go/v2/pkg/client/auth"
	"github.com/zitadel/zitadel-go/v2/pkg/client/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/middleware"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel"
)

type ClientInfo struct {
	Issuer  string
	Address string
	Project string
	Token   string
}

func GetClientInfo(d *schema.ResourceData) (*ClientInfo, error) {
	issuer := d.Get(IssuerVar).(string)
	address := d.Get(AddressVar).(string)
	projectID := d.Get(ProjectVar).(string)
	token := d.Get(TokenVar).(string)

	return &ClientInfo{
		issuer,
		address,
		projectID,
		token,
	}, nil
}

func getAuthClient(info *ClientInfo) (*auth.Client, error) {
	client, err := auth.NewClient(
		info.Issuer, info.Address,
		[]string{oidc.ScopeOpenID, zitadel.ScopeProjectID(info.Project)},
		zitadel.WithCustomURL(info.Issuer, info.Address),
		zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(info.Token)),
		zitadel.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start zitadel client: %v", err)
	}
	return client, nil
}

func getAdminClient(info *ClientInfo) (*admin.Client, error) {
	client, err := admin.NewClient(
		info.Issuer, info.Address,
		[]string{oidc.ScopeOpenID, zitadel.ScopeProjectID(info.Project)},
		//zitadel.WithCustomURL(info.Issuer, info.Address),
		zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(info.Token)),
		zitadel.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start zitadel client: %v", err)
	}
	return client, nil
}

func getManagementClient(info *ClientInfo, orgID string) (*management.Client, error) {
	opts := []zitadel.Option{
		zitadel.WithInsecure(),
		//zitadel.WithCustomURL(info.Issuer, info.Address),
		zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(info.Token)),
	}
	if orgID != "" {
		opts = append(opts, zitadel.WithOrgID(orgID))
	}

	client, err := management.NewClient(
		info.Issuer, info.Address,
		[]string{oidc.ScopeOpenID, zitadel.ScopeProjectID(info.Project)},
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start zitadel client: %v", err)
	}
	return client, nil
}
