package helper

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/zitadel-go/v2/pkg/client/admin"
	"github.com/zitadel/zitadel-go/v2/pkg/client/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/middleware"
	"github.com/zitadel/zitadel-go/v2/pkg/client/system"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel"
)

const (
	DomainVar       = "domain"
	InsecureVar     = "insecure"
	TokenVar        = "token"
	PortVar         = "port"
	SystemAPIUserID = "systemapi_user_id"
	SystemAPIKey    = "systemapi_certificate"
)

type ClientInfo struct {
	Domain          string
	Issuer          string
	SystemAPIUserID string
	SystemAPIKey    string
	Options         []zitadel.Option
}

func GetClientInfo(d *schema.ResourceData) (*ClientInfo, error) {
	insecure := d.Get(InsecureVar).(bool)
	domain := d.Get(DomainVar).(string)
	options := []zitadel.Option{zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(d.Get(TokenVar).(string)))}

	portStr := ""
	port := d.Get(PortVar)
	if port != nil {
		portStr = port.(string)
	}

	issuer := ""
	if portStr != "" {
		domain = domain + ":" + portStr
		issuer = domain
	} else {
		issuer = domain
		if insecure {
			domain = domain + ":80"
		} else {
			domain = domain + ":443"
		}
	}

	if insecure {
		options = append(options, zitadel.WithInsecure())
		issuer = "http://" + issuer
	} else {
		issuer = "https://" + issuer
	}

	userID := ""
	if userIDValue, ok := d.GetOk(SystemAPIUserID); ok {
		userID = userIDValue.(string)
	}
	key := ""
	if keyValue, ok := d.GetOk(SystemAPIKey); ok {
		key = keyValue.(string)
	}

	return &ClientInfo{
		domain,
		issuer,
		userID,
		key,
		options,
	}, nil
}

func GetAdminClient(info *ClientInfo) (*admin.Client, error) {
	client, err := admin.NewClient(
		info.Issuer, info.Domain,
		[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
		info.Options...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start zitadel client: %v", err)
	}

	return client, nil
}

func GetManagementClient(info *ClientInfo, orgID string) (*management.Client, error) {
	options := info.Options
	if orgID != "" {
		options = append(options, zitadel.WithOrgID(orgID))
	}

	client, err := management.NewClient(
		info.Issuer, info.Domain,
		[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
		options...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start zitadel client: %v", err)
	}
	return client, nil
}

func GetSystemClient(info *ClientInfo) (*system.Client, error) {
	if info.SystemAPIUserID == "" || info.SystemAPIKey == "" {
		return nil, fmt.Errorf("Missing userID or key for systemAPI")
	}

	client, err := system.NewClient(
		info.Issuer, info.Domain, system.JWTProfileFromPath(info.SystemAPIKey, info.SystemAPIUserID),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start zitadel client: %v", err)
	}
	return client, nil
}
