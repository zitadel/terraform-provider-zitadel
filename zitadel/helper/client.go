package helper

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oidcclient "github.com/zitadel/oidc/v3/pkg/client"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/client"
	"github.com/zitadel/zitadel-go/v3/pkg/client/middleware"
	actionV2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	featurev2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/feature/v2"
	instanceV2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/instance/v2"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	orgV2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/org/v2"
	settingsv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/settings/v2"
	userv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/user/v2"
	webkeys "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/webkey/v2"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DomainVar                 = "domain"
	DomainDescription         = "Domain used to connect to the ZITADEL instance"
	InsecureVar               = "insecure"
	InsecureDescription       = "Use insecure connection"
	AccessTokenVar            = "access_token"
	AccessTokenDescription    = "Personal Access Token to connect to ZITADEL. Either 'access_token', 'jwt_file', 'jwt_profile_file' or 'jwt_profile_json' is required"
	TokenVar                  = "token"
	TokenDescription          = "Path to the file containing credentials to connect to ZITADEL"
	PortVar                   = "port"
	PortDescription           = "Used port if not the default ports 80 or 443 are configured"
	JWTFileVar                = "jwt_file"
	JWTFileDescription        = "Path to the file containing presigned JWT to connect to ZITADEL. Either 'access_token', 'jwt_file', 'jwt_profile_file' or 'jwt_profile_json' is required"
	JWTProfileFileVar         = "jwt_profile_file"
	JWTProfileFileDescription = "Path to the file containing credentials to connect to ZITADEL. Either 'access_token', 'jwt_file', 'jwt_profile_file' or 'jwt_profile_json' is required"
	JWTProfileJSONVar         = "jwt_profile_json"
	JWTProfileJSONDescription = "JSON value of credentials to connect to ZITADEL. Either 'access_token', 'jwt_file', 'jwt_profile_file' or 'jwt_profile_json' is required"
)

type ClientInfo struct {
	// Configuration for lazy client creation
	zitadelConfig *zitadel.Zitadel
	authOptions   []client.Option

	// Lazy client instance (created on first use)
	clientOnce sync.Once
	clientInst *client.Client
	clientErr  error

	// Legacy fields for backward compatibility
	Domain  string
	Issuer  string
	KeyPath string
	Data    []byte
}

// ensureClient creates the client on first use (lazy initialization)
func (c *ClientInfo) ensureClient(ctx context.Context) error {
	c.clientOnce.Do(func() {
		c.clientInst, c.clientErr = client.New(ctx, c.zitadelConfig, c.authOptions...)
	})
	return c.clientErr
}

func GetClientInfo(ctx context.Context, insecure bool, domain string, accessToken string, token string, jwtFile string, jwtProfileFile string, jwtProfileJSON string, port string, insecureSkipVerifyTLS bool, transportHeaders map[string]string) (*ClientInfo, error) {
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")

	// Build zitadel.Zitadel options
	zitadelOpts := make([]zitadel.Option, 0)
	keyPath := ""

	// Handle port configuration
	if insecure {
		if port == "" || port == "80" {
			zitadelOpts = append(zitadelOpts, zitadel.WithInsecure("80"))
		} else {
			zitadelOpts = append(zitadelOpts, zitadel.WithInsecure(port))
		}
	} else if port != "" && port != "443" {
		portNum, err := strconv.Atoi(port)
		if err == nil {
			zitadelOpts = append(zitadelOpts, zitadel.WithPort(uint16(portNum)))
		}
	}

	// Add insecure_skip_verify_tls option
	if insecureSkipVerifyTLS {
		zitadelOpts = append(zitadelOpts, zitadel.WithInsecureSkipVerifyTLS())
	}

	// Add transport headers
	for k, v := range transportHeaders {
		zitadelOpts = append(zitadelOpts, zitadel.WithTransportHeader(k, v))
	}

	// Create zitadel.Zitadel instance
	z := zitadel.New(domain, zitadelOpts...)

	// Build client options for authentication
	clientOpts := make([]client.Option, 0)

	// Set up authentication
	var tokenSourceInit client.TokenSourceInitializer
	if accessToken != "" {
		tokenSourceInit = func(ctx context.Context, issuer string) (oauth2.TokenSource, error) {
			return oauth2.StaticTokenSource(&oauth2.Token{
				AccessToken: accessToken,
				TokenType:   "Bearer",
			}), nil
		}
		keyPath = ""
	} else if token != "" {
		if _, err := os.Stat(token); err != nil {
			return nil, fmt.Errorf("failed to read token file: %v", err)
		}
		keyFile, err := oidcclient.ConfigFromKeyFile(token)
		if err != nil {
			return nil, fmt.Errorf("failed to load key file: %v", err)
		}
		// CRITICAL: Include scopes for JWT authentication
		tokenSourceInit = client.JWTAuthentication(keyFile, oidc.ScopeOpenID, client.ScopeZitadelAPI())
		keyPath = token
	} else if jwtFile != "" {
		jwt, err := os.ReadFile(jwtFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read JWT file: %v", err)
		}
		tokenSourceInit = func(ctx context.Context, issuer string) (oauth2.TokenSource, error) {
			return oauth2.StaticTokenSource(&oauth2.Token{
				AccessToken: string(jwt),
				TokenType:   "Bearer",
			}), nil
		}
	} else if jwtProfileFile != "" {
		if _, err := os.Stat(jwtProfileFile); err != nil {
			return nil, fmt.Errorf("failed to read jwt_profile_file: %v", err)
		}
		keyFile, err := oidcclient.ConfigFromKeyFile(jwtProfileFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load jwt profile file: %v", err)
		}
		// CRITICAL: Include scopes for JWT authentication
		tokenSourceInit = client.JWTAuthentication(keyFile, oidc.ScopeOpenID, client.ScopeZitadelAPI())
		keyPath = jwtProfileFile
	} else if jwtProfileJSON != "" {
		keyFile, err := oidcclient.ConfigFromKeyFileData([]byte(jwtProfileJSON))
		if err != nil {
			return nil, fmt.Errorf("failed to parse JWT profile JSON: %v", err)
		}
		// CRITICAL: Include scopes for JWT authentication
		tokenSourceInit = client.JWTAuthentication(keyFile, oidc.ScopeOpenID, client.ScopeZitadelAPI())
	} else {
		return nil, fmt.Errorf("either 'access_token', 'jwt_file', 'jwt_profile_file' or 'jwt_profile_json' is required")
	}

	clientOpts = append(clientOpts, client.WithAuth(tokenSourceInit))

	// Return ClientInfo with lazy client creation (no global singleton)
	return &ClientInfo{
		zitadelConfig: z,
		authOptions:   clientOpts,
		Domain:        z.Host(),
		Issuer:        z.Origin(),
		KeyPath:       keyPath,
		Data:          []byte(jwtProfileJSON),
	}, nil
}

// Service client getters - all use lazy client creation

func GetSecuritySettingsClient(ctx context.Context, info *ClientInfo) (settingsv2.SettingsServiceClient, error) {
	if err := info.ensureClient(ctx); err != nil {
		return nil, err
	}
	return info.clientInst.SettingsServiceV2(), nil
}

func GetOrgClient(ctx context.Context, info *ClientInfo) (orgV2.OrganizationServiceClient, error) {
	if err := info.ensureClient(ctx); err != nil {
		return nil, err
	}
	return info.clientInst.OrganizationServiceV2(), nil
}

func GetActionClient(ctx context.Context, info *ClientInfo) (actionV2.ActionServiceClient, error) {
	if err := info.ensureClient(ctx); err != nil {
		return nil, err
	}
	return info.clientInst.ActionServiceV2(), nil
}

func GetFeatureClient(ctx context.Context, info *ClientInfo) (featurev2.FeatureServiceClient, error) {
	if err := info.ensureClient(ctx); err != nil {
		return nil, err
	}
	return info.clientInst.FeatureServiceV2(), nil
}

func GetWebKeyClient(ctx context.Context, info *ClientInfo) (webkeys.WebKeyServiceClient, error) {
	if err := info.ensureClient(ctx); err != nil {
		return nil, err
	}
	return info.clientInst.WebkeyServiceV2(), nil
}

func GetManagementClient(ctx context.Context, info *ClientInfo) (management.ManagementServiceClient, error) {
	if err := info.ensureClient(ctx); err != nil {
		return nil, err
	}
	return info.clientInst.ManagementService(), nil
}

func GetOrgV2Client(ctx context.Context, info *ClientInfo) (orgV2.OrganizationServiceClient, error) {
	if err := info.ensureClient(ctx); err != nil {
		return nil, err
	}
	return info.clientInst.OrganizationServiceV2(), nil
}

func GetUserV2Client(ctx context.Context, info *ClientInfo) (userv2.UserServiceClient, error) {
	if err := info.ensureClient(ctx); err != nil {
		return nil, err
	}
	return info.clientInst.UserServiceV2(), nil
}

func GetAdminClient(ctx context.Context, info *ClientInfo) (admin.AdminServiceClient, error) {
	if err := info.ensureClient(ctx); err != nil {
		return nil, err
	}
	return info.clientInst.AdminService(), nil
}

func GetInstanceClient(ctx context.Context, info *ClientInfo) (instanceV2.InstanceServiceClient, error) {
	if err := info.ensureClient(ctx); err != nil {
		return nil, err
	}
	return info.clientInst.InstanceServiceV2(), nil
}

func CtxWithID(ctx context.Context, d *schema.ResourceData) context.Context {
	return CtxSetOrgID(ctx, GetID(d, OrgIDVar))
}

func CtxWithOrgID(ctx context.Context, d *schema.ResourceData) context.Context {
	return CtxSetOrgID(ctx, d.Get(OrgIDVar).(string))
}

func CtxSetOrgID(ctx context.Context, orgID string) context.Context {
	return middleware.SetOrgID(ctx, orgID)
}

func IgnoreIfNotFoundError(err error) error {
	if code := status.Code(err); code == codes.NotFound {
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
