package helper

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"
)

func generateTestRSAKeys(t *testing.T) (privatePEM, publicPEM []byte) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	privatePEM = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	pubASN1, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		t.Fatalf("failed to marshal public key: %v", err)
	}
	publicPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})
	return
}

func generateTestPKCS8Key(t *testing.T) []byte {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	pkcs8, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("failed to marshal PKCS8 key: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8,
	})
}

func TestParseRSAPrivateKey_PKCS1(t *testing.T) {
	privatePEM, _ := generateTestRSAKeys(t)
	key, err := parseRSAPrivateKey(privatePEM)
	if err != nil {
		t.Fatalf("parseRSAPrivateKey(PKCS1) error = %v", err)
	}
	if key == nil {
		t.Fatal("parseRSAPrivateKey(PKCS1) returned nil key")
	}
}

func TestParseRSAPrivateKey_PKCS8(t *testing.T) {
	pkcs8PEM := generateTestPKCS8Key(t)
	key, err := parseRSAPrivateKey(pkcs8PEM)
	if err != nil {
		t.Fatalf("parseRSAPrivateKey(PKCS8) error = %v", err)
	}
	if key == nil {
		t.Fatal("parseRSAPrivateKey(PKCS8) returned nil key")
	}
}

func TestParseRSAPrivateKey_InvalidPEM(t *testing.T) {
	_, err := parseRSAPrivateKey([]byte("not a pem"))
	if err == nil {
		t.Fatal("parseRSAPrivateKey(invalid) expected error, got nil")
	}
}

func TestParseRSAPrivateKey_UnsupportedType(t *testing.T) {
	badPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: []byte("fake"),
	})
	_, err := parseRSAPrivateKey(badPEM)
	if err == nil {
		t.Fatal("parseRSAPrivateKey(EC) expected error, got nil")
	}
}

func TestParseRSAPublicKey_PKIX(t *testing.T) {
	_, publicPEM := generateTestRSAKeys(t)
	key, err := parseRSAPublicKey(publicPEM)
	if err != nil {
		t.Fatalf("parseRSAPublicKey(PKIX) error = %v", err)
	}
	if key == nil {
		t.Fatal("parseRSAPublicKey(PKIX) returned nil key")
	}
}

func TestParseRSAPublicKey_PKCS1(t *testing.T) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}
	pkcs1PEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&rsaKey.PublicKey),
	})
	key, err := parseRSAPublicKey(pkcs1PEM)
	if err != nil {
		t.Fatalf("parseRSAPublicKey(PKCS1) error = %v", err)
	}
	if key == nil {
		t.Fatal("parseRSAPublicKey(PKCS1) returned nil key")
	}
}

func TestParseRSAPublicKey_InvalidPEM(t *testing.T) {
	_, err := parseRSAPublicKey([]byte("not a pem"))
	if err == nil {
		t.Fatal("parseRSAPublicKey(invalid) expected error, got nil")
	}
}

func TestParseRSAPublicKey_UnsupportedType(t *testing.T) {
	badPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: []byte("fake"),
	})
	_, err := parseRSAPublicKey(badPEM)
	if err == nil {
		t.Fatal("parseRSAPublicKey(CERTIFICATE) expected error, got nil")
	}
}

func TestVerifyMatchingPublicKey_Match(t *testing.T) {
	privatePEM, publicPEM := generateTestRSAKeys(t)
	err := verifyMatchingPublicKey(privatePEM, publicPEM)
	if err != nil {
		t.Fatalf("verifyMatchingPublicKey() error = %v", err)
	}
}

func TestVerifyMatchingPublicKey_Mismatch(t *testing.T) {
	privatePEM, _ := generateTestRSAKeys(t)
	_, otherPublicPEM := generateTestRSAKeys(t)
	err := verifyMatchingPublicKey(privatePEM, otherPublicPEM)
	if err == nil {
		t.Fatal("verifyMatchingPublicKey(mismatch) expected error, got nil")
	}
}

func TestVerifyMatchingPublicKey_InvalidPrivateKey(t *testing.T) {
	_, publicPEM := generateTestRSAKeys(t)
	err := verifyMatchingPublicKey([]byte("bad"), publicPEM)
	if err == nil {
		t.Fatal("verifyMatchingPublicKey(bad private) expected error, got nil")
	}
}

func TestVerifyMatchingPublicKey_InvalidPublicKey(t *testing.T) {
	privatePEM, _ := generateTestRSAKeys(t)
	err := verifyMatchingPublicKey(privatePEM, []byte("bad"))
	if err == nil {
		t.Fatal("verifyMatchingPublicKey(bad public) expected error, got nil")
	}
}

func TestNewSystemAPITokenSourceFromPEM(t *testing.T) {
	privatePEM, _ := generateTestRSAKeys(t)
	ts, err := NewSystemAPITokenSourceFromPEM(privatePEM, "user-id", "https://example.com")
	if err != nil {
		t.Fatalf("NewSystemAPITokenSourceFromPEM() error = %v", err)
	}
	if ts == nil {
		t.Fatal("NewSystemAPITokenSourceFromPEM() returned nil")
	}
}

func TestNewSystemAPITokenSourceFromPEM_EmptyKey(t *testing.T) {
	_, err := NewSystemAPITokenSourceFromPEM(nil, "user-id", "https://example.com")
	if err == nil {
		t.Fatal("NewSystemAPITokenSourceFromPEM(empty key) expected error, got nil")
	}
}

func TestNewSystemAPITokenSourceFromPEM_EmptyUser(t *testing.T) {
	privatePEM, _ := generateTestRSAKeys(t)
	_, err := NewSystemAPITokenSourceFromPEM(privatePEM, "", "https://example.com")
	if err == nil {
		t.Fatal("NewSystemAPITokenSourceFromPEM(empty user) expected error, got nil")
	}
}

func TestNewSystemAPITokenSourceFromPEM_EmptyAudience(t *testing.T) {
	privatePEM, _ := generateTestRSAKeys(t)
	_, err := NewSystemAPITokenSourceFromPEM(privatePEM, "user-id", "")
	if err == nil {
		t.Fatal("NewSystemAPITokenSourceFromPEM(empty audience) expected error, got nil")
	}
}

func TestSystemAPITokenSource_Token(t *testing.T) {
	privatePEM, _ := generateTestRSAKeys(t)
	ts, err := NewSystemAPITokenSourceFromPEM(privatePEM, "user-id", "https://example.com")
	if err != nil {
		t.Fatalf("NewSystemAPITokenSourceFromPEM() error = %v", err)
	}
	token, err := ts.Token()
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}
	if token.AccessToken == "" {
		t.Fatal("Token() returned empty access token")
	}
	if token.TokenType != "Bearer" {
		t.Fatalf("Token() type = %s, want Bearer", token.TokenType)
	}
}
