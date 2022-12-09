package helper

import (
	"fmt"

	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/zitadel-go/v2/pkg/client/admin"
	"github.com/zitadel/zitadel-go/v2/pkg/client/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/middleware"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel"
)

const (
	DomainVar   = "domain"
	InsecureVar = "insecure"
	TokenVar    = "token"
	PortVar     = "port"
)

type ClientInfo struct {
	Domain  string
	Issuer  string
	Options []zitadel.Option
}

func GetClientInfo(insecure bool, domain string, token string, port string) (*ClientInfo, error) {
	options := []zitadel.Option{zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(token))}

	issuer := ""
	if port != "" {
		domain = domain + ":" + port
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

	return &ClientInfo{
		domain,
		issuer,
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
