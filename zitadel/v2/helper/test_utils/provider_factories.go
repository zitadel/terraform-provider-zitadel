package test_utils

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func ZitadelProviderFactories(provider *schema.Provider) map[string]func() (*schema.Provider, error) {
	return map[string]func() (*schema.Provider, error){
		"zitadel": func() (*schema.Provider, error) {
			return provider, nil
		},
	}
}
