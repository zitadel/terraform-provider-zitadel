package test_utils

import (
	"strings"
	"testing"
)

func CreateDefaultDependency(t *testing.T, datasourceType string, idField string, newDependencyID func() (string, error)) (string, string) {
	datasourceExample, datasourceExampleAttributes := ReadExample(t, Datasources, datasourceType)
	dependencyID, err := newDependencyID()
	if err != nil {
		t.Fatalf("failed to create dependency for %s: %v", datasourceType, err)
	}
	return strings.Replace(datasourceExample, AttributeValue(t, idField, datasourceExampleAttributes).AsString(), dependencyID, 1), dependencyID
}
