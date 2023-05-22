package project

import "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/project"

const (
	projectIDVar              = "project_id"
	nameVar                   = "name"
	stateVar                  = "state"
	orgIDVar                  = "org_id"
	roleAssertionVar          = "project_role_assertion"
	roleCheckVar              = "project_role_check"
	hasProjectCheckVar        = "has_project_check"
	privateLabelingSettingVar = "private_labeling_setting"
)

var (
	defaultPrivateLabelingSetting = project.PrivateLabelingSetting_name[0]
)