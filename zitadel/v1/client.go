package v1

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/zitadel-go/pkg/client/management"
	"github.com/zitadel/zitadel-go/pkg/client/middleware"
	"github.com/zitadel/zitadel-go/pkg/client/zitadel"
)

type ClientInfo struct {
	Issuer  string
	Address string
	Project string
	Token   string
}

func GetClientInfo(d *schema.ResourceData) (*ClientInfo, error) {
	issuer := d.Get(issuerVar).(string)
	address := d.Get(addressVar).(string)
	projectID := d.Get(projectVar).(string)
	token := d.Get(tokenVar).(string)

	return &ClientInfo{
		issuer,
		address,
		projectID,
		token,
	}, nil
}

func getManagementClient(clientinfo *ClientInfo, orgID string) (*management.Client, error) {
	opts := []zitadel.Option{
		zitadel.WithCustomURL(clientinfo.Issuer, clientinfo.Address),
		zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(clientinfo.Token)),
	}
	if orgID != "" {
		opts = append(opts, zitadel.WithOrgID(orgID))
	}

	client, err := management.NewClient(
		[]string{oidc.ScopeOpenID, zitadel.ScopeProjectID(clientinfo.Project)},
		opts...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start zitadel client: %v", err)
	}
	/*defer func() {
		err := client.Connection.Close()
		if err != nil {
			log.Println("could not close grpc connection", err)
		}
	}()*/
	return client, nil
}
