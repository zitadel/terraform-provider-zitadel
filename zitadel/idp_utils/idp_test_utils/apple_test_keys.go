package idp_test_utils

// Throwaway EC P-256 (PKCS#8) keys generated only for acceptance tests. They are
// not real Apple keys and grant no access. ZITADEL (>= v4.17.0) validates the
// Apple IDP private key format (apple.BytesToPrivateKey / PKCS#8 ECDSA), so the
// test key must be a parseable PKCS#8 EC key rather than an arbitrary string.
//
// The keys are single-line with literal "\n" escapes so they can be injected into
// double-quoted HCL attributes (private_key = "..."); Terraform decodes the
// escapes back into the newlines a PEM block requires.
const (
	AppleTestPrivateKey        = `-----BEGIN PRIVATE KEY-----\nMIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQggyIGmyz7WNh0s+QO\neb8jFtYMFtfx//DVKpslP9OpvhehRANCAAQXXnHyzVHP/ktHZdtHfeuGAllVARd5\nERNWy0mrrAChvyRM0NDRinucx4AFOw+fkyjTruAjYfRzdUngUCHh0t6L\n-----END PRIVATE KEY-----\n`
	AppleTestPrivateKeyUpdated = `-----BEGIN PRIVATE KEY-----\nMIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgx+exjEhUyw5eAG0I\nbdfweBmXQFAGrAej4Iidh4nbHr+hRANCAAQxp47ZssQKCyVkuLQ9spjYDt9ouBSY\nfqxaUrPBcWS/SKcnN8ydpJY7UPVqOTdO+3xrvTHHjnbSO2+dvP59n393\n-----END PRIVATE KEY-----\n`
)
