package org_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/org"
)

func TestAccOrgDatasource_ID(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org")
	orgName := "org_datasource_" + frame.UniqueResourcesID
	otherFrame := frame.AnotherOrg(t, orgName)
	test_utils.RunDatasourceTest(
		t,
		otherFrame.BaseTestFrame,
		otherFrame.AsOrgDefaultDependency,
		nil,
		map[string]string{
			"id":    otherFrame.OrgID,
			"name":  orgName,
			"state": "ORG_STATE_ACTIVE",
		},
	)
}

func TestAccOrgsDatasources_ID_Name_Match(t *testing.T) {
	datasourceName := "zitadel_orgs"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	config, attributes := test_utils.ReadExample(t, test_utils.Datasources, datasourceName)
	exampleName := test_utils.AttributeValue(t, org.NameVar, attributes).AsString()
	exampleDomain := test_utils.AttributeValue(t, org.DomainVar, attributes).AsString()
	orgName := fmt.Sprintf("%s-%s", exampleName, frame.UniqueResourcesID)
	// for-each is not supported in acceptance tests, so we cut the example down to the first block
	// https://github.com/hashicorp/terraform-plugin-sdk/issues/536
	config = strings.Join(strings.Split(config, "\n")[0:7], "\n")
	config = strings.Replace(config, exampleName, orgName, 1)
	config = strings.Replace(config, exampleDomain, orgName, 1)
	otherFrame := frame.AnotherOrg(t, orgName)
	test_utils.RunDatasourceTest(
		t,
		otherFrame.BaseTestFrame,
		config,
		checkRemoteProperty(otherFrame, idFromFrame(otherFrame))(orgName),
		map[string]string{
			"ids.0": otherFrame.OrgID,
			"ids.#": "1",
		},
	)
}

func TestAccOrgsDatasources_ID_Name_Mismatch(t *testing.T) {
	datasourceName := "zitadel_orgs"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)
	config, attributes := test_utils.ReadExample(t, test_utils.Datasources, datasourceName)
	orgName := fmt.Sprintf("%s_%s", test_utils.AttributeValue(t, org.NameVar, attributes), frame.UniqueResourcesID)
	otherFrame := frame.AnotherOrg(t, orgName)
	test_utils.RunDatasourceTest(
		t,
		otherFrame.BaseTestFrame,
		config,
		checkRemoteProperty(otherFrame, idFromFrame(otherFrame))(orgName),
		map[string]string{"ids.#": "0"},
	)
}

func idFromFrame(frame *test_utils.OrgTestFrame) func(state *terraform.State) string {
	return func(state *terraform.State) string {
		return frame.OrgID
	}
}
