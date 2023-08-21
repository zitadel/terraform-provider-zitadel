package machine_user

import (
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/user"
)

const (
	UserIDVar             = "user_id"
	orgIDVar              = "org_id"
	userStateVar          = "state"
	UserNameVar           = "user_name"
	loginNamesVar         = "login_names"
	preferredLoginNameVar = "preferred_login_name"
	nameVar               = "name"
	DescriptionVar        = "description"
	accessTokenTypeVar    = "access_token_type"
)

var (
	defaultAccessTokenType = user.AccessTokenType_name[0]
)
