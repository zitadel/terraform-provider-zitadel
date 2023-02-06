package helper

import (
	"fmt"

	"github.com/zitadel/oidc/pkg/oidc"
	"github.com/zitadel/zitadel-go/v2/pkg/client/admin"
	"github.com/zitadel/zitadel-go/v2/pkg/client/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/middleware"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DomainVar      = "domain"
	InsecureVar    = "insecure"
	TokenVar       = "token"
	PortVar        = "port"
	JWTProfileFile = "jwt_profile_file"
	JWTProfileJSON = "jwt_profile_json"
)

type ClientInfo struct {
	Domain  string
	Issuer  string
	KeyPath string
	Options []zitadel.Option
}

func GetClientInfo(insecure bool, domain string, jwtProfileFile string, jwtProfileJSON string, token string, port string) (*ClientInfo, error) {
	options := []zitadel.Option{}
	if jwtProfileFile != "" {
		options = append(options, zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(jwtProfileFile)))
	} else if jwtProfileJSON != "" {
		options = append(options, zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromFileData([]byte(jwtProfileJSON))))
	} else {
		options = append(options, zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(token)))
	}

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
		token,
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

func IgnoreIfNotFoundError(err error) error {
	//permission denied included as nothing can be found then as well
	if code := status.Code(err); code == codes.NotFound || code == codes.PermissionDenied {
		return nil
	}
	return err
}

func IgnorePreconditionError(err error) error {
	if code := status.Code(err); code == codes.FailedPrecondition {
		return nil
	}
	return err
}

func IgnoreAlreadyExistsError(err error) error {
	if code := status.Code(err); code == codes.AlreadyExists {
		return nil
	}
	return err
}
