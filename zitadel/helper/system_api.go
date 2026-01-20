package helper

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"
)

const systemAPITokenTTL = time.Hour

// NewSystemAPITokenSourceFromPEM returns a TokenSource that signs short-lived JWTs with the provided RSA private key.
// The issuer and subject are set to systemAPIUser, the audience defaults to the provider issuer if not overridden.
func NewSystemAPITokenSourceFromPEM(keyPEM []byte, systemAPIUser, audience string) (oauth2.TokenSource, error) {
	if len(keyPEM) == 0 {
		return nil, fmt.Errorf("system api key is empty")
	}
	if systemAPIUser == "" {
		return nil, fmt.Errorf("system api user is empty")
	}
	if audience == "" {
		return nil, fmt.Errorf("system api audience is empty")
	}

	signer, err := newSystemAPISigner(keyPEM)
	if err != nil {
		return nil, err
	}

	return oauth2.ReuseTokenSource(nil, &systemAPITokenSource{
		signer:  signer,
		issuer:  systemAPIUser,
		aud:     audience,
		expTime: systemAPITokenTTL,
	}), nil
}

func newSystemAPISigner(keyPEM []byte) (jose.Signer, error) {
	key, err := parseRSAPrivateKey(keyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse system api key: %w", err)
	}

	signer, err := jose.NewSigner(
		jose.SigningKey{Algorithm: jose.RS256, Key: key},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create signer for system api key: %w", err)
	}
	return signer, nil
}

func parseRSAPrivateKey(keyPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(keyPEM)
	if block == nil {
		return nil, fmt.Errorf("no valid PEM data found")
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("unsupported private key type %T", key)
		}
		return rsaKey, nil
	default:
		return nil, fmt.Errorf("unsupported private key type %q", block.Type)
	}
}

type systemAPITokenSource struct {
	signer  jose.Signer
	issuer  string
	aud     string
	expTime time.Duration
}

func (s *systemAPITokenSource) Token() (*oauth2.Token, error) {
	now := time.Now()
	claims := jwt.Claims{
		Issuer:   s.issuer,
		Subject:  s.issuer,
		Audience: jwt.Audience{s.aud},
		IssuedAt: jwt.NewNumericDate(now),
		Expiry:   jwt.NewNumericDate(now.Add(s.expTime)),
	}

	raw, err := jwt.Signed(s.signer).Claims(claims).Serialize()
	if err != nil {
		return nil, fmt.Errorf("failed to sign system api jwt: %w", err)
	}

	return &oauth2.Token{
		AccessToken: raw,
		TokenType:   oidc.BearerToken,
		Expiry:      now.Add(s.expTime),
	}, nil
}
