package v2

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/zitadel-go/v2/pkg/client/admin"
	"github.com/zitadel/zitadel-go/v2/pkg/client/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/middleware"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel"
)

const (
	DomainVar   = "domain"
	InsecureVar = "insecure"
	ProjectVar  = "project"
	TokenVar    = "token"
)

type ClientInfo struct {
	Domain   string
	Insecure bool
	Project  string
	Token    string
}

func GetClientInfo(d *schema.ResourceData) (*ClientInfo, error) {
	return &ClientInfo{
		d.Get(DomainVar).(string),
		d.Get(InsecureVar).(bool),
		d.Get(ProjectVar).(string),
		d.Get(TokenVar).(string),
	}, nil
}

func getAdminClient(info *ClientInfo) (*admin.Client, error) {
	options := []zitadel.Option{zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(info.Token))}
	issuer := info.Domain
	if info.Insecure {
		options = append(options, zitadel.WithInsecure())
		issuer = "http://" + issuer
	} else {
		issuer = "https://" + issuer
	}

	client, err := admin.NewClient(
		issuer, info.Domain,
		[]string{oidc.ScopeOpenID, zitadel.ScopeProjectID(info.Project)},
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start zitadel client: %v", err)
	}
	return client, nil
}

func getManagementClient(info *ClientInfo, orgID string) (*management.Client, error) {
	options := []zitadel.Option{zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(info.Token))}
	issuer := info.Domain
	if info.Insecure {
		options = append(options, zitadel.WithInsecure())
		issuer = "http://" + issuer
	} else {
		issuer = "https://" + issuer
	}
	if orgID != "" {
		options = append(options, zitadel.WithOrgID(orgID))
	}

	client, err := management.NewClient(
		issuer, info.Domain,
		[]string{oidc.ScopeOpenID, zitadel.ScopeProjectID(info.Project)},
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start zitadel client: %v", err)
	}
	return client, nil
}
