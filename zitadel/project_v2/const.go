package project_v2

import projectpb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/project/v2"

const (
	ProjectIDVar              = "project_id"
	projectIDsVar             = "project_ids"
	NameVar                   = "name"
	nameMethodVar             = "name_method"
	stateVar                  = "state"
	roleAssertionVar          = "project_role_assertion"
	roleCheckVar              = "project_role_check"
	hasProjectCheckVar        = "has_project_check"
	privateLabelingSettingVar = "private_labeling_setting"
)

var (
	defaultPrivateLabelingSetting = projectpb.PrivateLabelingSetting_name[0]
)
