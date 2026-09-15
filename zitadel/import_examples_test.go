package zitadel

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// importCommand matches a documented import command with a single quoted, double quoted or bare ID.
// Commands that build the ID with a shell substitution are not matched.
var importCommand = regexp.MustCompile(`^terraform import (\S+) (?:'([^']*)'|"([^"$]*)"|([^\s'"$]+))$`)

// TestResourceImportExamples runs the documented import command of every resource through the resource's
// importer, so the documented ID format and the parsing of the import ID cannot drift apart.
func TestResourceImportExamples(t *testing.T) {
	resources := Provider().ResourcesMap
	names := make([]string, 0, len(resources))
	for name := range resources {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		res := resources[name]
		if res.Importer == nil || res.Importer.StateContext == nil {
			continue
		}
		examplePath := filepath.Join("..", "examples", "provider", "resources", strings.TrimPrefix(name, "zitadel_")+"-import.sh")
		content, err := os.ReadFile(examplePath)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			t.Fatalf("failed to read %s: %v", examplePath, err)
		}
		for _, line := range strings.Split(string(content), "\n") {
			match := importCommand.FindStringSubmatch(line)
			if match == nil {
				continue
			}
			address, id := match[1], match[2]+match[3]+match[4]
			t.Run(name+"/"+id, func(t *testing.T) {
				if !strings.HasPrefix(address, name+".") {
					t.Errorf("%s imports %s instead of a %s resource", examplePath, address, name)
				}
				data := res.TestResourceData()
				data.SetId(id)
				_, err := res.Importer.StateContext(context.Background(), data, nil)
				// Importers that look up the remote resource while importing fail without a client after parsing
				// the ID, those lookups are covered by the acceptance tests.
				if err != nil && err.Error() != "failed to get client" {
					t.Errorf("documented import ID %q does not match the importer of %s: %v", id, name, err)
				}
			})
		}
	}
}
