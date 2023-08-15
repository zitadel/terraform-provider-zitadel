package test_utils

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var _ resource.ImportStateIdFunc = ImportNothing

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
		return val, nil
	}
}

func ImportNothing(_ *terraform.State) (string, error) { return "", nil }

func ConcatImportStateIdFuncs(funcs ...resource.ImportStateIdFunc) resource.ImportStateIdFunc {
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
