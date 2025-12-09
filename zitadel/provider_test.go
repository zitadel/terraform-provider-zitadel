package zitadel

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

// TestProvider_AllResourcesRegistered scans the filesystem for resource
// packages and verifies that each discovered resource is properly registered
// in the provider's ResourcesMap. This test ensures no resource package is
// accidentally omitted from provider registration.
func TestProvider_AllResourcesRegistered(t *testing.T) {
	p := Provider()

	discoveredResources, err := discoverResources(".")
	require.NoError(t, err, "Failed to discover resources")

	var missingResources []string
	for _, resourceName := range discoveredResources {
		if !slices.Contains(getRegisteredResourceNames(p), resourceName) {
			missingResources = append(missingResources, resourceName)
		}
	}

	assert.Empty(t, missingResources,
		"The following resources were discovered but not registered: %v",
		missingResources)

	t.Logf("Total resources registered: %d", len(p.ResourcesMap))
	t.Logf("Total resources discovered: %d", len(discoveredResources))
}

// TestProvider_AllDatasourcesRegistered scans the filesystem for datasource
// packages and verifies that each discovered datasource is properly registered
// in the provider's DataSourcesMap. This test ensures no datasource package
// is accidentally omitted from provider registration.
func TestProvider_AllDatasourcesRegistered(t *testing.T) {
	p := Provider()

	discoveredDatasources, err := discoverDatasources(".")
	require.NoError(t, err, "Failed to discover datasources")

	var missingDatasources []string
	for _, datasourceName := range discoveredDatasources {
		if !slices.Contains(getRegisteredDatasourceNames(p), datasourceName) {
			missingDatasources = append(missingDatasources, datasourceName)
		}
	}

	assert.Empty(t, missingDatasources,
		"The following datasources were discovered but not registered: %v",
		missingDatasources)

	t.Logf("Total datasources registered: %d", len(p.DataSourcesMap))
	t.Logf("Total datasources discovered: %d", len(discoveredDatasources))
}

// TestProvider_ResourceSchemaExactlyOneOfConsistency validates that all
// registered resources have consistent ExactlyOneOf field references. This
// ensures that mutually exclusive field groups are properly configured across
// all resource schemas.
func TestProvider_ResourceSchemaExactlyOneOfConsistency(t *testing.T) {
	p := Provider()

	for resourceName, resource := range p.ResourcesMap {
		t.Run(resourceName, func(t *testing.T) {
			checkSchemaExactlyOneOfConsistency(t, resource.Schema)
		})
	}
}

func getRegisteredResourceNames(p *schema.Provider) []string {
	names := make([]string, 0, len(p.ResourcesMap))
	for name := range p.ResourcesMap {
		names = append(names, name)
	}
	return names
}

func getRegisteredDatasourceNames(p *schema.Provider) []string {
	names := make([]string, 0, len(p.DataSourcesMap))
	for name := range p.DataSourcesMap {
		names = append(names, name)
	}
	return names
}

func discoverResources(baseDir string) ([]string, error) {
	excludeDirs := []string{"helper", "pat"}

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	var resources []string
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		if slices.Contains(excludeDirs, entry.Name()) {
			continue
		}

		resourcePath := filepath.Join(baseDir, entry.Name(), "resource.go")
		if _, err := os.Stat(resourcePath); err != nil {
			continue
		}

		content, err := os.ReadFile(resourcePath)
		if err != nil {
			continue
		}

		if strings.Contains(string(content), "func GetResource()") {
			resources = append(resources, "zitadel_"+entry.Name())
		}
	}

	return resources, nil
}

func discoverDatasources(baseDir string) ([]string, error) {
	excludeDirs := []string{"helper", "pat"}

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil, err
	}

	var datasources []string
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		if slices.Contains(excludeDirs, entry.Name()) {
			continue
		}

		datasourcePath := filepath.Join(baseDir, entry.Name(), "datasource.go")
		if _, err := os.Stat(datasourcePath); err != nil {
			continue
		}

		content, err := os.ReadFile(datasourcePath)
		if err != nil {
			continue
		}

		if strings.Contains(string(content), "func GetDatasource()") {
			datasources = append(datasources, "zitadel_"+entry.Name())
		}
	}

	return datasources, nil
}

func checkSchemaExactlyOneOfConsistency(t *testing.T, schemaMap map[string]*schema.Schema) {
	tests := []struct {
		name      string
		checkFunc func(t *testing.T, fieldName string, fieldSchema *schema.Schema)
	}{
		{
			name: "referenced fields exist",
			checkFunc: func(t *testing.T, fieldName string, fieldSchema *schema.Schema) {
				for _, refFieldName := range fieldSchema.ExactlyOneOf {
					assert.Contains(t, schemaMap, refFieldName,
						"Field %q references non-existent field %q in ExactlyOneOf",
						fieldName, refFieldName)
				}
			},
		},
		{
			name: "field includes itself",
			checkFunc: func(t *testing.T, fieldName string, fieldSchema *schema.Schema) {
				assert.Contains(t, fieldSchema.ExactlyOneOf, fieldName,
					"Field %q has ExactlyOneOf set but does not include itself in the list",
					fieldName)
			},
		},
		{
			name: "group consistency",
			checkFunc: func(t *testing.T, fieldName string, fieldSchema *schema.Schema) {
				for _, refFieldName := range fieldSchema.ExactlyOneOf {
					if refFieldName == fieldName {
						continue
					}

					refSchema, exists := schemaMap[refFieldName]
					if !exists {
						continue
					}

					if len(refSchema.ExactlyOneOf) == 0 {
						t.Errorf("Field %q references %q in ExactlyOneOf, but %q does not have ExactlyOneOf set",
							fieldName, refFieldName, refFieldName)
						continue
					}

					assert.ElementsMatch(t, fieldSchema.ExactlyOneOf, refSchema.ExactlyOneOf,
						"Field %q and %q are in the same ExactlyOneOf group but have inconsistent field lists",
						fieldName, refFieldName)
				}
			},
		},
	}

	for fieldName, fieldSchema := range schemaMap {
		if len(fieldSchema.ExactlyOneOf) == 0 {
			continue
		}

		for _, tt := range tests {
			t.Run(fieldName+"_"+tt.name, func(t *testing.T) {
				tt.checkFunc(t, fieldName, fieldSchema)
			})
		}
	}
}
