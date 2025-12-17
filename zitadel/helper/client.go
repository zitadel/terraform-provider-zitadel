package helper

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	actionV2 "github.com/zitadel/zitadel-go/v3/pkg/client/action/v2"
	"github.com/zitadel/zitadel-go/v3/pkg/client/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/middleware"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel"
	"golang.org/x/net/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DomainVar                  = "domain"
	DomainDescription          = "Domain used to connect to the ZITADEL instance"
	InsecureVar                = "insecure"
	InsecureDescription        = "Use insecure connection"
	TokenVar                   = "token"
	TokenDescription           = "Path to the file containing credentials to connect to ZITADEL"
	PortVar                    = "port"
	PortDescription            = "Used port if not the default ports 80 or 443 are configured"
	JWTFileVar                 = "jwt_file"
	JWTFileDescription         = "Path to the file containing presigned JWT to connect to ZITADEL. Either 'jwt_file', 'jwt_profile_file' or 'jwt_profile_json' is required"
	JWTProfileFileVar          = "jwt_profile_file"
	JWTProfileFileDescription  = "Path to the file containing credentials to connect to ZITADEL. Either 'jwt_file', 'jwt_profile_file' or 'jwt_profile_json' is required"
	JWTProfileJSONVar          = "jwt_profile_json"
	JWTProfileJSONDescription  = "JSON value of credentials to connect to ZITADEL. Either 'jwt_file', 'jwt_profile_file' or 'jwt_profile_json' is required"
	ProxyBlockVar              = "proxy"
	ProxyBlockDescription      = "Proxy configuration."
	ProxyURLVar                = "url"
	ProxyURLDescription        = "URL of the proxy to be used."
	ProxyAuthHeaderVar         = "auth_header"
	ProxyAuthHeaderDescription = "Proxy authentication header."
)

type ClientInfo struct {
	Domain  string
	Issuer  string
	KeyPath string
	Data    []byte
	Options []zitadel.Option
}

func GetClientInfo(ctx context.Context, insecure bool, domain string, token string, jwtFile string, jwtProfileFile string, jwtProfileJSON string, port string, proxyURL string, proxyAuth string) (*ClientInfo, error) {
	options := make([]zitadel.Option, 0)
	keyPath := ""
	if token != "" {
		options = append(options, zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(context.Background(), token)))
		keyPath = token
	} else if jwtFile != "" {
		jwt, err := os.ReadFile(jwtFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read JWT file: %v", err)
		}
		options = append(options, zitadel.WithJWTDirectTokenSource(string(jwt)))
	} else if jwtProfileFile != "" {
		options = append(options, zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromPath(context.Background(), jwtProfileFile)))
		keyPath = jwtProfileFile
	} else if jwtProfileJSON != "" {
		options = append(options, zitadel.WithJWTProfileTokenSource(middleware.JWTProfileFromFileData(context.Background(), []byte(jwtProfileJSON))))
	} else {
		return nil, fmt.Errorf("either 'jwt_file', 'jwt_profile_file' or 'jwt_profile_json' is required")
	}

	if proxyURL != "" {
		parsedURL, err := url.Parse(proxyURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse proxy url: %w", err)
		}

		var grpcDialer func(context.Context, string) (net.Conn, error)

		switch parsedURL.Scheme {
		case "http", "https":
			// This is the corrected logic for HTTP/HTTPS proxies.
			grpcDialer = func(ctx context.Context, addr string) (net.Conn, error) {
				proxyConn, err := net.Dial("tcp", parsedURL.Host)
				if err != nil {
					return nil, fmt.Errorf("failed to connect to http proxy: %w", err)
				}
				connectReq := &http.Request{
					Method: "CONNECT",
					URL:    &url.URL{Opaque: addr},
					Host:   addr,
					Header: make(http.Header),
				}
				// Use the provided value to set the standard Proxy-Authorization header.
				if proxyAuth != "" {
					connectReq.Header.Set("Proxy-Authorization", proxyAuth)
				}
				if err := connectReq.Write(proxyConn); err != nil {
					return nil, err
				}
				br := bufio.NewReader(proxyConn)
				resp, err := http.ReadResponse(br, connectReq)
				if err != nil {
					return nil, err
				}
				if resp.StatusCode != http.StatusOK {
					return nil, fmt.Errorf("http proxy connect failed with status %s", resp.Status)
				}
				return proxyConn, nil
			}

		case "socks5":
			// This logic for SOCKS5 proxies is already correct.
			dialer, err := proxy.FromURL(parsedURL, proxy.Direct)
			if err != nil {
				return nil, fmt.Errorf("failed to create socks5 dialer: %w", err)
			}
			grpcDialer = func(ctx context.Context, addr string) (net.Conn, error) {
				return dialer.Dial("tcp", addr)
			}

		default:
			return nil, fmt.Errorf("unsupported proxy scheme: %s. Use 'http', 'https', or 'socks5'", parsedURL.Scheme)
		}
		options = append(options, zitadel.WithDialOptions(grpc.WithContextDialer(grpcDialer)))
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
