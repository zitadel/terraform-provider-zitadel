package idp_test_utils

// Throwaway EC P-256 (PKCS#8) keys generated only for acceptance tests. They are
// not real Apple keys and grant no access. ZITADEL (>= v4.17.0) validates the
// Apple IDP private key format (apple.BytesToPrivateKey / PKCS#8 ECDSA), so the
// test key must be a parseable PKCS#8 EC key rather than an arbitrary string.
const (
	AppleTestPrivateKey = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQggyIGmyz7WNh0s+QO
eb8jFtYMFtfx//DVKpslP9OpvhehRANCAAQXXnHyzVHP/ktHZdtHfeuGAllVARd5
ERNWy0mrrAChvyRM0NDRinucx4AFOw+fkyjTruAjYfRzdUngUCHh0t6L
-----END PRIVATE KEY-----
`
	AppleTestPrivateKeyUpdated = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgx+exjEhUyw5eAG0I
bdfweBmXQFAGrAej4Iidh4nbHr+hRANCAAQxp47ZssQKCyVkuLQ9spjYDt9ouBSY
fqxaUrPBcWS/SKcnN8ydpJY7UPVqOTdO+3xrvTHHjnbSO2+dvP59n393
-----END PRIVATE KEY-----
`
)
