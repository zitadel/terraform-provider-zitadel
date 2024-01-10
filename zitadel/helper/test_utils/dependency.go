package test_utils

import (
	"strings"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
)

func CreateDefaultDependency(t *testing.T, datasourceType string, idField string, newDependencyID func() (string, error)) (string, string) {
	datasourceExample, datasourceExampleAttributes := ReadExample(t, Datasources, datasourceType)
	dependencyID, err := newDependencyID()
	if err != nil {
		t.Fatalf("failed to create dependency for %s: %v", datasourceType, err)
	}
	return strings.Replace(datasourceExample, AttributeValue(t, idField, datasourceExampleAttributes).AsString(), dependencyID, 1), dependencyID
}

func CreateOrgDefaultDependency(t *testing.T, datasourceType string, orgID string, idField string, newDependencyID func() (string, error)) (string, string) {
	datasourceExample, datasourceExampleAttributes := ReadExample(t, Datasources, datasourceType)
	dependencyID, err := newDependencyID()
	if err != nil {
		t.Fatalf("failed to create dependency for %s: %v", datasourceType, err)
	}
	// org exampleID does always have to be different then the idField value
	dep := strings.Replace(datasourceExample, AttributeValue(t, idField, datasourceExampleAttributes).AsString(), dependencyID, 1)
	// replace the exmaple OrgID with the given org
	return strings.Replace(dep, AttributeValue(t, helper.OrgIDVar, datasourceExampleAttributes).AsString(), orgID, 1), dependencyID
}

func CreateProjectDefaultDependency(t *testing.T, datasourceType string, orgID string, projectField string, projectID string, idField string, newDependencyID func() (string, error)) (string, string) {
	datasourceExample, datasourceExampleAttributes := ReadExample(t, Datasources, datasourceType)
	dependencyID, err := newDependencyID()
	if err != nil {
		t.Fatalf("failed to create dependency for %s: %v", datasourceType, err)
	}
	// org exampleID does always have to be different then the idField value
	dep := strings.Replace(datasourceExample, AttributeValue(t, idField, datasourceExampleAttributes).AsString(), dependencyID, 1)
	dep = strings.Replace(dep, AttributeValue(t, projectField, datasourceExampleAttributes).AsString(), projectID, 1)
	// replace the exmaple OrgID with the given org
	return strings.Replace(dep, AttributeValue(t, helper.OrgIDVar, datasourceExampleAttributes).AsString(), orgID, 1), dependencyID
}
