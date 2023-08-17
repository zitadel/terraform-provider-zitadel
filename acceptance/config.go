package acceptance

import (
	_ "embed"
	"encoding/json"
)

var (
	//go:embed keys/instance-level-admin-sa.json
	instanceLevelAdminSAJSON []byte

	//go:embed keys/org-level-admin-sa.json
	orgLevelAdminSAJSON []byte

	//go:embed config.json
	configJson []byte
)

type Config struct {
	OrgLevel      IsolatedInstance
	InstanceLevel IsolatedInstance
}

type IsolatedInstance struct {
	Domain      string
	AdminSAJSON []byte
}

func GetConfig() Config {
	val := Config{
		OrgLevel:      IsolatedInstance{AdminSAJSON: orgLevelAdminSAJSON},
		InstanceLevel: IsolatedInstance{AdminSAJSON: instanceLevelAdminSAJSON},
	}
	if err := json.Unmarshal(configJson, &val); err != nil {
		panic(err)
	}
	return val
}
