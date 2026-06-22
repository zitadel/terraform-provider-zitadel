package helper

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	actionV2 "github.com/zitadel/zitadel-go/v3/pkg/client/action/v2"
	"github.com/zitadel/zitadel-go/v3/pkg/client/admin"
	appv2 "github.com/zitadel/zitadel-go/v3/pkg/client/application/v2"
	featurev2 "github.com/zitadel/zitadel-go/v3/pkg/client/feature/v2"
	instanceV2 "github.com/zitadel/zitadel-go/v3/pkg/client/instance/v2"
	"github.com/zitadel/zitadel-go/v3/pkg/client/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/middleware"
	orgV2 "github.com/zitadel/zitadel-go/v3/pkg/client/org/v2"
	projectv2 "github.com/zitadel/zitadel-go/v3/pkg/client/project/v2"
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
	AccessTokenDescription    = "Personal Access Token to connect to ZITADEL. Either 'access_token', 'jwt_file', 'jwt_profile_file', 'jwt_profile_json' or 'system_api' is required"
	TokenVar                  = "token"
	TokenDescription          = "Deprecated: Use 'access_token' for Personal Access Tokens or 'jwt_profile_file' for JWT Profile credentials instead. Path to the file containing a JWT Profile key to connect to ZITADEL."
	PortVar                   = "port"
	PortDescription           = "Used port if not the default ports 80 or 443 are configured"
	JWTFileVar                = "jwt_file"
	JWTFileDescription        = "Path to the file containing presigned JWT to connect to ZITADEL. Either 'access_token', 'jwt_file', 'jwt_profile_file', 'jwt_profile_json' or 'system_api' is required"
	JWTProfileFileVar         = "jwt_profile_file"
	JWTProfileFileDescription = "Path to the file containing credentials to connect to ZITADEL. Either 'access_token', 'jwt_file', 'jwt_profile_file', 'jwt_profile_json' or 'system_api' is required"
	JWTProfileJSONVar         = "jwt_profile_json"
	JWTProfileJSONDescription = "JSON value of credentials to connect to ZITADEL. Either 'access_token', 'jwt_file', 'jwt_profile_file', 'jwt_profile_json' or 'system_api' is required"
	SystemAPIVar              = "system_api"
	SystemAPIDescription      = "Configuration block for authenticating with the ZITADEL System API using a PEM encoded key."
	SystemAPIKeyFileAttr      = "key_file"
	SystemAPIKeyFileDesc      = "Path to the PEM encoded private key for a ZITADEL System API user. Either 'key_file'/'key' or the 'private_key'+'public_key' pair is required when using System API authentication."
	SystemAPIKeyAttr          = "key"
	SystemAPIKeyDesc          = "PEM encoded private key for a ZITADEL System API user. Either 'key_file'/'key' or the 'private_key'+'public_key' pair is required when using System API authentication."
	SystemAPIPrivateKeyAttr   = "private_key"
	SystemAPIPrivateKeyDesc   = "PEM encoded private key for a ZITADEL System API user when provided separately from the public key. Use together with 'public_key'."
	SystemAPIPublicKeyAttr    = "public_key"
	SystemAPIPublicKeyDesc    = "PEM encoded public key for a ZITADEL System API user when provided separately from the private key. Use together with 'private_key'."
	SystemAPIUserAttr         = "user"
	SystemAPIUserDesc         = "User ID configured for the System API key. Used as both issuer and subject in the self-signed JWT."
	SystemAPIAudienceAttr     = "audience"
	SystemAPIAudienceDesc     = "Audience to set on the System API JWT. Defaults to the issuer derived from domain/port if omitted."
)

// ClientInfo holds everything the provider needs to talk to ZITADEL: the
// connection target (Domain/Issuer), the gRPC dial Options, and the credential
// material. The gRPC client authenticates purely from Options, but asset uploads
// run over a separate plain-HTTP path (see form.go) that cannot read Options, so
// the fields below expose the same credential and headers to that path.
type ClientInfo struct {
	Domain  string
	Issuer  string
	KeyPath string
	Data    []byte
	Options []zitadel.Option
	// TokenSource is the bearer credential for the auth modes whose token the
	// asset-upload path cannot derive from KeyPath/Data: access_token (PAT),
	// jwt_file and system_api. It is nil for the JWT-profile file modes, which
	// form.go builds from KeyPath or Data instead.
	TokenSource oauth2.TokenSource
	// TransportHeaders is the provider's transport_headers. The gRPC client gets
	// them via Options; the asset-upload path applies them from here so both
	// transports send the same headers (e.g. proxy auth like GCP IAP).
	TransportHeaders map[string]string
}

func GetClientInfo(ctx context.Context, insecure bool, domain string, accessToken string, token string, jwtFile string, jwtProfileFile string, jwtProfileJSON string, systemAPIKeyFile string, systemAPIKey string, systemAPIPrivateKey string, systemAPIPublicKey string, systemAPIUser string, systemAPIAudience string, port string, insecureSkipVerifyTLS bool, transportHeaders map[string]string) (*ClientInfo, error) {
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")

	options := make([]zitadel.Option, 0)

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

	if insecureSkipVerifyTLS {
		options = append(options, zitadel.WithInsecureSkipVerifyTLS())
	}

	for k, v := range transportHeaders {
		options = append(options, zitadel.WithTransportHeader(k, v))
	}

	// keyPath/keyData feed the asset-upload path for the JWT-profile file modes;
	// assetTokenSource feeds it for the bearer-token modes (set in the branches
	// below). Exactly one of them ends up populated.
	keyPath := ""
	var keyData []byte
	var assetTokenSource oauth2.TokenSource

	switch {
	case accessToken != "":
		// Personal Access Token: a static bearer token, reused for both transports.
		tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken, TokenType: "Bearer"})
		options = append(options, zitadel.WithTokenSource(tokenSource))
		assetTokenSource = tokenSource
	case token != "":
		// Deprecated 'token': path to a JWT-profile key. Asset uploads read it via keyPath.
		if _, err := os.Stat(token); err != nil {
			return nil, fmt.Errorf("failed to read token file: %v", err)
		}
		options = append(options, zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(context.Background(), token)))
		keyPath = token
	case jwtFile != "":
		// Presigned JWT used directly as a bearer token. Trim once so the gRPC
		// option and the asset-upload bearer use an identical value; a trailing
		// newline (common in CLI-written files) is invalid in an auth header.
		jwt, err := os.ReadFile(jwtFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read JWT file: %v", err)
		}
		presignedJWT := strings.TrimSpace(string(jwt))
		options = append(options, zitadel.WithJWTDirectTokenSource(presignedJWT))
		assetTokenSource = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: presignedJWT, TokenType: "Bearer"})
	case jwtProfileFile != "":
		// JWT-profile key on disk. Asset uploads exchange it for a token via keyPath.
		if _, err := os.Stat(jwtProfileFile); err != nil {
			return nil, fmt.Errorf("failed to read jwt_profile_file: %v", err)
		}
		options = append(options, zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(context.Background(), jwtProfileFile)))
		keyPath = jwtProfileFile
	case jwtProfileJSON != "":
		// JWT-profile key inline. Asset uploads exchange it for a token via keyData.
		options = append(options, zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromFileData(context.Background(), []byte(jwtProfileJSON))))
		keyData = []byte(jwtProfileJSON)
	case systemAPIKeyFile != "" || systemAPIKey != "" || systemAPIPrivateKey != "" || systemAPIPublicKey != "":
		// System API: a self-signed JWT token source, reused for both transports.
		if systemAPIUser == "" {
			return nil, fmt.Errorf("system_api.user is required when using System API authentication")
		}

		var keyPEM []byte
		switch {
		case systemAPIPrivateKey != "" || systemAPIPublicKey != "":
			if systemAPIPrivateKey == "" || systemAPIPublicKey == "" {
				return nil, fmt.Errorf("both system_api.private_key and system_api.public_key must be set when providing split keys")
			}
			keyPEM = []byte(systemAPIPrivateKey)
			if err := verifyMatchingPublicKey(keyPEM, []byte(systemAPIPublicKey)); err != nil {
				return nil, fmt.Errorf("system api keys do not match: %w", err)
			}
		case systemAPIKeyFile != "":
			data, err := os.ReadFile(systemAPIKeyFile)
			if err != nil {
				return nil, fmt.Errorf("failed to read system_api_key_file: %v", err)
			}
			keyPEM = data
		default:
			keyPEM = []byte(systemAPIKey)
		}

		audience := systemAPIAudience
		if audience == "" {
			audience = issuer
		}

		ts, err := NewSystemAPITokenSourceFromPEM(keyPEM, systemAPIUser, audience)
		if err != nil {
			return nil, fmt.Errorf("failed to create system api token source: %w", err)
		}
		options = append(options, zitadel.WithTokenSource(ts))
		assetTokenSource = ts
	default:
		return nil, fmt.Errorf("either 'access_token', 'jwt_file', 'jwt_profile_file', 'jwt_profile_json' or 'system_api' (with 'key', 'key_file', or both 'private_key' and 'public_key') is required")
	}

	return &ClientInfo{
		Domain:           clientDomain,
		Issuer:           issuer,
		KeyPath:          keyPath,
		Data:             keyData,
		Options:          options,
		TokenSource:      assetTokenSource,
		TransportHeaders: transportHeaders,
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

var projectV2ClientLock sync.Mutex
var projectV2Client atomic.Pointer[projectv2.Client]

func GetProjectV2Client(ctx context.Context, info *ClientInfo) (*projectv2.Client, error) {
	// Lock-free fast path once the client is initialised: the atomic load
	// synchronises with the atomic store below, so there is no data race and
	// initialised reads never contend on the mutex. The mutex only guards
	// initialisation, and a failed init is not cached so a later call retries.
	if client := projectV2Client.Load(); client != nil {
		return client, nil
	}
	projectV2ClientLock.Lock()
	defer projectV2ClientLock.Unlock()
	if client := projectV2Client.Load(); client != nil {
		return client, nil
	}
	client, err := projectv2.NewClient(ctx,
		info.Issuer, info.Domain,
		[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
		info.Options...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start zitadel project v2 client: %w", err)
	}
	time.Sleep(time.Second * 2)
	projectV2Client.Store(client)
	return client, nil
}

var appV2ClientLock sync.Mutex
var appV2Client atomic.Pointer[appv2.Client]

func GetAppV2Client(ctx context.Context, info *ClientInfo) (*appv2.Client, error) {
	// Lock-free fast path once the client is initialised: the atomic load
	// synchronises with the atomic store below, so there is no data race and
	// initialised reads never contend on the mutex. The mutex only guards
	// initialisation, and a failed init is not cached so a later call retries.
	if client := appV2Client.Load(); client != nil {
		return client, nil
	}
	appV2ClientLock.Lock()
	defer appV2ClientLock.Unlock()
	if client := appV2Client.Load(); client != nil {
		return client, nil
	}
	client, err := appv2.NewClient(ctx,
		info.Issuer, info.Domain,
		[]string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()},
		info.Options...,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start zitadel app v2 client: %w", err)
	}
	time.Sleep(time.Second * 2)
	appV2Client.Store(client)
	return client, nil
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
