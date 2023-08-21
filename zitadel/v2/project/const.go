package project

import "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/project"

const (
	ProjectIDVar              = "project_id"
	NameVar                   = "name"
	stateVar                  = "state"
	roleAssertionVar          = "project_role_assertion"
	roleCheckVar              = "project_role_check"
	hasProjectCheckVar        = "has_project_check"
	privateLabelingSettingVar = "private_labeling_setting"
)

var (
	defaultPrivateLabelingSetting = project.PrivateLabelingSetting_name[0]
)
