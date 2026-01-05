package helper

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	actionV2 "github.com/zitadel/zitadel-go/v3/pkg/client/action/v2"
	"github.com/zitadel/zitadel-go/v3/pkg/client/admin"
	featurev2 "github.com/zitadel/zitadel-go/v3/pkg/client/feature/v2"
	instanceV2 "github.com/zitadel/zitadel-go/v3/pkg/client/instance/v2"
	"github.com/zitadel/zitadel-go/v3/pkg/client/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/middleware"
	orgV2 "github.com/zitadel/zitadel-go/v3/pkg/client/org/v2"
	settingsv2 "github.com/zitadel/zitadel-go/v3/pkg/client/settings/v2"
	userv2 "github.com/zitadel/zitadel-go/v3/pkg/client/user/v2"
	webkeys "github.com/zitadel/zitadel-go/v3/pkg/client/webkey/v2"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel"
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
	Domain  string
	Issuer  string
	KeyPath string
	Data    []byte
	Options []zitadel.Option
}

func GetClientInfo(ctx context.Context, insecure bool, domain string, accessToken string, token string, jwtFile string, jwtProfileFile string, jwtProfileJSON string, port string) (*ClientInfo, error) {
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	options := make([]zitadel.Option, 0)
	keyPath := ""
	if accessToken != "" {
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: accessToken,
			TokenType:   "Bearer",
		})
		options = append(options, zitadel.WithTokenSource(tokenSource))
	} else if token != "" {
		if _, err := os.Stat(token); err != nil {
			return nil, fmt.Errorf("failed to read token file: %v", err)
		}
		options = append(options, zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(context.Background(), token)))
		keyPath = token
	} else if jwtFile != "" {
		jwt, err := os.ReadFile(jwtFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read JWT file: %v", err)
		}
		options = append(options, zitadel.WithJWTDirectTokenSource(string(jwt)))
	} else if jwtProfileFile != "" {
		if _, err := os.Stat(jwtProfileFile); err != nil {
			return nil, fmt.Errorf("failed to read jwt_profile_file: %v", err)
		}
		options = append(options, zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(context.Background(), jwtProfileFile)))
		keyPath = jwtProfileFile
	} else if jwtProfileJSON != "" {
		options = append(options, zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromFileData(context.Background(), []byte(jwtProfileJSON))))
	} else {
		return nil, fmt.Errorf("either 'access_token', 'jwt_file', 'jwt_profile_file' or 'jwt_profile_json' is required")
	}

	issuerScheme := "https://"
	if insecure {
		options = append(options, zitadel.WithInsecure())
		issuerScheme = "http://"
	}

	issuerPort := port
	if port == "80" && insecure || port == "443" && !insecure {
		issuerPort = ""
	}

	issuer := issuerScheme + domain
	if issuerPort != "" {
		issuer += ":" + issuerPort
	}

	clientDomain := domain + ":" + port
	if port == "" {
		clientDomain = domain + ":443"
		if insecure {
			clientDomain = domain + ":80"
		}
	}

	return &ClientInfo{
		clientDomain,
		issuer,
		keyPath,
		[]byte(jwtProfileJSON),
		options,
	}, nil
}

var securitySettingsClientLock = &sync.Mutex{}
var securitySettingsClient *settingsv2.Client

func GetSecuritySettingsClient(ctx context.Context, info *ClientInfo) (*settingsv2.Client, error) {
	if securitySettingsClient == nil {
		securitySettingsClientLock.Lock()
		defer securitySettingsClientLock.Unlock()
		if securitySettingsClient == nil {
			client, err := settingsv2.NewClient(ctx,
				info.Issuer, info.Domain,
				[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
				info.Options...,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to start zitadel client: %v", err)
			}
			time.Sleep(time.Second * 2)
			securitySettingsClient = client
		}
	}
	return securitySettingsClient, nil
}

var orgClientLock = &sync.Mutex{}
var orgClient *orgV2.Client

func GetOrgClient(ctx context.Context, info *ClientInfo) (*orgV2.Client, error) {
	if orgClient == nil {
		orgClientLock.Lock()
		defer orgClientLock.Unlock()
		if orgClient == nil {
			client, err := orgV2.NewClient(ctx,
				info.Issuer, info.Domain,
				[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
				info.Options...,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to start zitadel client: %v", err)
			}
			time.Sleep(time.Second * 2)
			orgClient = client
		}
	}
	return orgClient, nil
}

var actionClientLock = &sync.Mutex{}
var actionClient *actionV2.Client

func GetActionClient(ctx context.Context, info *ClientInfo) (*actionV2.Client, error) {
	if actionClient == nil {
		actionClientLock.Lock()
		defer actionClientLock.Unlock()
		if actionClient == nil {
			client, err := actionV2.NewClient(ctx,
				info.Issuer, info.Domain,
				[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
				info.Options...,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to start zitadel client: %v", err)
			}
			time.Sleep(time.Second * 2)
			actionClient = client
		}
	}
	return actionClient, nil
}

var featureClientLock = &sync.Mutex{}
var featureClient *featurev2.Client

func GetFeatureClient(ctx context.Context, info *ClientInfo) (*featurev2.Client, error) {
	if featureClient == nil {
		featureClientLock.Lock()
		defer featureClientLock.Unlock()
		if featureClient == nil {
			client, err := featurev2.NewClient(ctx,
				info.Issuer, info.Domain,
				[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
				info.Options...,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to start zitadel feature client: %v", err)
			}
			time.Sleep(time.Second * 2)
			featureClient = client
		}
	}
	return featureClient, nil
}

var webkeyClientLock = &sync.Mutex{}
var webkeyClient *webkeys.Client

func GetWebKeyClient(ctx context.Context, info *ClientInfo) (*webkeys.Client, error) {
	if webkeyClient == nil {
		webkeyClientLock.Lock()
		defer webkeyClientLock.Unlock()
		if webkeyClient == nil {
			client, err := webkeys.NewClient(ctx,
				info.Issuer, info.Domain,
				[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
				info.Options...,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to start zitadel client: %v", err)
			}
			time.Sleep(time.Second * 2)
			webkeyClient = client
		}
	}
	return webkeyClient, nil
}

var mgmtClientLock = &sync.Mutex{}
var mgmtClient *management.Client

func GetManagementClient(ctx context.Context, info *ClientInfo) (*management.Client, error) {
	if mgmtClient == nil {
		mgmtClientLock.Lock()
		defer mgmtClientLock.Unlock()
		if mgmtClient == nil {
			client, err := management.NewClient(ctx,
				info.Issuer, info.Domain,
				[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
				info.Options...,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to start zitadel client: %v", err)
			}
			time.Sleep(time.Second * 2)
			mgmtClient = client
		}
	}
	return mgmtClient, nil
}

var orgV2ClientLock = &sync.Mutex{}
var orgV2Client *orgV2.Client

func GetOrgV2Client(ctx context.Context, info *ClientInfo) (*orgV2.Client, error) {
	if orgV2Client == nil {
		orgV2ClientLock.Lock()
		defer orgV2ClientLock.Unlock()
		if orgV2Client == nil {
			client, err := orgV2.NewClient(ctx,
				info.Issuer, info.Domain,
				[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
				info.Options...,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to start zitadel org v2 client: %v", err)
			}
			time.Sleep(time.Second * 2)
			orgV2Client = client
		}
	}
	return orgV2Client, nil
}

var userV2ClientLock = &sync.Mutex{}
var userV2Client *userv2.Client

func GetUserV2Client(ctx context.Context, info *ClientInfo) (*userv2.Client, error) {
	if userV2Client == nil {
		userV2ClientLock.Lock()
		defer userV2ClientLock.Unlock()
		if userV2Client == nil {
			client, err := userv2.NewClient(ctx,
				info.Issuer, info.Domain,
				[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
				info.Options...,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to start zitadel user v2 client: %v", err)
			}
			time.Sleep(time.Second * 2)
			userV2Client = client
		}
	}
	return userV2Client, nil
}

var adminClientLock = &sync.Mutex{}
var adminClient *admin.Client

func GetAdminClient(ctx context.Context, info *ClientInfo) (*admin.Client, error) {
	if adminClient == nil {
		adminClientLock.Lock()
		defer adminClientLock.Unlock()
		if adminClient == nil {
			client, err := admin.NewClient(ctx,
				info.Issuer, info.Domain,
				[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
				info.Options...,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to start zitadel client: %v", err)
			}
			time.Sleep(time.Second * 2)
			adminClient = client
		}
	}
	return adminClient, nil
}

var instanceClientLock = &sync.Mutex{}
var instanceClient *instanceV2.Client

func GetInstanceClient(ctx context.Context, info *ClientInfo) (*instanceV2.Client, error) {
	if instanceClient == nil {
		instanceClientLock.Lock()
		defer instanceClientLock.Unlock()
		if instanceClient == nil {
			client, err := instanceV2.NewClient(ctx,
				info.Issuer, info.Domain,
				[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
				info.Options...,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to start zitadel client: %v", err)
			}
			time.Sleep(time.Second * 2)
			instanceClient = client
		}
	}
	return instanceClient, nil
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
