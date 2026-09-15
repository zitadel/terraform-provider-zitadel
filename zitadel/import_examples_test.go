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

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

var (
	// importCommand matches a documented import command with a single quoted, double quoted or bare ID.
	importCommand = regexp.MustCompile(`^terraform import (\S+) (?:'([^']*)'|"(.*)"|([^\s'"]+))$`)
	// escapedFileSubstitution matches the documented way to pass a file that contains :, which escapes every : with
	// the SemicolonPlaceholder.
	escapedFileSubstitution = regexp.MustCompile(`\$\(cat \S+ \| sed -e 's/:/` + helper.SemicolonPlaceholder + `/g'\)`)
)

// exampleKeyDetails stands in for the key file of an escapedFileSubstitution.
const exampleKeyDetails = `{"type":"serviceaccount","keyId":"123456789012345678","key":"-----BEGIN RSA PRIVATE KEY-----\nMIIEpQ...\n-----END RSA PRIVATE KEY-----\n","userId":"123456789012345678"}`

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
			if !strings.HasPrefix(line, "terraform import ") {
				continue
			}
			line = escapedFileSubstitution.ReplaceAllLiteralString(line, strings.ReplaceAll(exampleKeyDetails, ":", helper.SemicolonPlaceholder))
			match := importCommand.FindStringSubmatch(line)
			if match == nil || strings.Contains(line, "$(") {
				t.Errorf("%s contains an import command that cannot be verified: %s", examplePath, line)
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
				if strings.Contains(data.Id(), helper.SemicolonPlaceholder) {
					t.Errorf("%s was imported with the escaped ID %q", name, data.Id())
				}
				for key := range res.Schema {
					if value, ok := data.GetOk(key); ok {
						if s, isString := value.(string); isString && strings.Contains(s, helper.SemicolonPlaceholder) {
							t.Errorf("%s was imported without unescaping %s: %q", name, key, s)
						}
					}
				}
			})
		}
	}
}
