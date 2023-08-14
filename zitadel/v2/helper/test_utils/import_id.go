package test_utils

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var _ resource.ImportStateIdFunc = ImportNothing

func ImportStateIdWithOrg(frame *OrgTestFrame) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		return fmt.Sprintf("%s:%s", frame.State(state).ID, frame.OrgID), nil
	}
}

func ImportStateId(frame BaseTestFrame, withAttribute ...string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		primary := frame.State(state)
		id := primary.ID
		for _, v := range withAttribute {
			id += fmt.Sprintf(":%s", primary.Attributes[v])
		}
		return id, nil
	}
}

func ImportNothing(_ *terraform.State) (string, error) { return "", nil }
