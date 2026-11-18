package idp_ldap_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_ldap"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils/idp_test_utils"
)

func TestAccInstanceIdPLDAP(t *testing.T) {
	idp_test_utils.RunInstanceIDPLifecyleTest(t, "zitadel_idp_ldap", idp_ldap.BindPasswordVar)
}

func TestAccInstanceIdPLDAPUpdateKeepsUnchangedSecret(t *testing.T) {
	idp_test_utils.RunInstanceIDPUnchangedSecretTest(t, "zitadel_idp_ldap", "instance.idp.ldap.v2.changed", idp_ldap.BindPasswordVar, "bindPassword")
}
