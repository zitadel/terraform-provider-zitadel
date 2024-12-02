package test_utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

var (
	_                           resource.ImportStateIdFunc = ImportNothing
	ZitadelGeneratedIdPattern                              = `\d{18}`
	ZitadelGeneratedIdOnlyRegex                            = regexp.MustCompile(fmt.Sprintf(`^%s$`, ZitadelGeneratedIdPattern))
)

func ImportResourceId(frame BaseTestFrame) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		return frame.State(state).ID, nil
	}
}

func ImportOrgId(frame *OrgTestFrame) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		return frame.OrgID, nil
	}
}

func ImportStateAttribute(frame BaseTestFrame, attr string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		primary := frame.State(state)
		val, ok := primary.Attributes[attr]
		if !ok {
			return "", fmt.Errorf("attribute %s not found in attributes %+v", attr, primary.Attributes)
		}
		if val != "" {
			val = fmt.Sprintf(`"%s"`, strings.ReplaceAll(val, ":", helper.SemicolonPlaceholder))
		}
		return val, nil
	}
}

func ImportNothing(_ *terraform.State) (string, error) { return "", nil }

// ChainImportStateIdFuncs returns a function that composes an ID that has exactly the same
// length of semicolon separated parts as the passed number of funcs
func ChainImportStateIdFuncs(funcs ...resource.ImportStateIdFunc) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		parts := make([]string, len(funcs))
		for i, f := range funcs {
			part, err := f(state)
			if err != nil {
				return "", err
			}
			parts[i] = part
		}
		return strings.Join(parts, ":"), nil
	}
}
