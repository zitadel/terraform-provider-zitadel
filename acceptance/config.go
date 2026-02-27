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

	//go:embed keys/system-api-sa.pem
	systemAPIKeyPEM []byte

	//go:embed config.json
	configJson []byte
)

type Config struct {
	OrgLevel      IsolatedInstance  `json:"orgLevel"`
	InstanceLevel IsolatedInstance  `json:"instanceLevel"`
	SystemAPI     SystemAPIInstance `json:"systemAPI"`
}

type IsolatedInstance struct {
	Domain      string
	AdminSAJSON []byte
}

type SystemAPIInstance struct {
	Domain string `json:"domain"`
	User   string `json:"user"`
	KeyPEM []byte `json:"-"`
}

func GetConfig() Config {
	val := Config{
		OrgLevel:      IsolatedInstance{AdminSAJSON: orgLevelAdminSAJSON},
		InstanceLevel: IsolatedInstance{AdminSAJSON: instanceLevelAdminSAJSON},
		SystemAPI: SystemAPIInstance{
			KeyPEM: systemAPIKeyPEM,
			User:   "system-api-sa",
		},
	}
	if err := json.Unmarshal(configJson, &val); err != nil {
		panic(err)
	}
	if val.SystemAPI.Domain == "" {
		val.SystemAPI.Domain = val.OrgLevel.Domain
	}
	if len(val.SystemAPI.KeyPEM) == 0 {
		val.SystemAPI.KeyPEM = systemAPIKeyPEM
	}
	if val.SystemAPI.User == "" {
		val.SystemAPI.User = "system-api-sa"
	}
	return val
}
