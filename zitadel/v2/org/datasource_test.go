package org_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccOrgDatasource_ID(t *testing.T) {
	datasourceName := "zitadel_org"
	frame, err := test_utils.NewOrgTestFrame(datasourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	orgName := "org_datasource_" + frame.UniqueResourcesID
	otherFrame, err := frame.AnotherOrg(orgName)
	if err != nil {
		t.Fatalf("could not switch to another org: %v", err)
	}
	test_utils.RunDatasourceTest(
		t,
		otherFrame.BaseTestFrame,
		fmt.Sprintf(`
data "%s" "%s" {
  id = "%s"
}
`, datasourceName, otherFrame.UniqueResourcesID, otherFrame.OrgID),
		nil,
		map[string]string{
			"id":    otherFrame.OrgID,
			"name":  orgName,
			"state": "ORG_STATE_ACTIVE",
		},
		nil,
	)
}

func TestAccOrgDatasource_Name(t *testing.T) {
	datasourceName := "zitadel_org"
	frame, err := test_utils.NewOrgTestFrame(datasourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	orgName := "org_datasource_" + frame.UniqueResourcesID
	otherFrame, err := frame.AnotherOrg(orgName)
	if err != nil {
		t.Fatalf("could not switch to another org: %v", err)
	}
	test_utils.RunDatasourceTest(
		t,
		otherFrame.BaseTestFrame,
		fmt.Sprintf(`
data "%s" "%s" {
  name = "%s"
}
`, datasourceName, otherFrame.UniqueResourcesID, orgName),
		checkRemoteProperty(otherFrame)(orgName),
		map[string]string{
			"id":    otherFrame.OrgID,
			"name":  orgName,
			"state": "ORG_STATE_ACTIVE",
		},
		nil,
	)
}

func TestAccOrgDatasource_ID_Name_Match(t *testing.T) {
	datasourceName := "zitadel_org"
	frame, err := test_utils.NewOrgTestFrame(datasourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	orgName := "org_datasource_" + frame.UniqueResourcesID
	otherFrame, err := frame.AnotherOrg(orgName)
	if err != nil {
		t.Fatalf("could not switch to another org: %v", err)
	}
	test_utils.RunDatasourceTest(
		t,
		otherFrame.BaseTestFrame,
		fmt.Sprintf(`
data "%s" "%s" {
  name = "%s"
}
`, datasourceName, otherFrame.UniqueResourcesID, orgName),
		nil,
		map[string]string{
			"id":    otherFrame.OrgID,
			"name":  orgName,
			"state": "ORG_STATE_ACTIVE",
		},
		nil,
	)
}

func TestAccOrgDatasource_ID_Name_Mismatch(t *testing.T) {
	datasourceName := "zitadel_org"
	frame, err := test_utils.NewOrgTestFrame(datasourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	orgName := "org_datasource_" + frame.UniqueResourcesID
	otherFrame, err := frame.AnotherOrg(orgName)
	if err != nil {
		t.Fatalf("could not switch to another org: %v", err)
	}
	test_utils.RunDatasourceTest(
		t,
		otherFrame.BaseTestFrame,
		fmt.Sprintf(`
data "%s" "%s" {
  name = "mismatching_org_name"
}
`, datasourceName, otherFrame.UniqueResourcesID),
		nil,
		nil,
		regexp.MustCompile("the filters don't match exactly 1 org, but 0 orgs"),
	)
}
