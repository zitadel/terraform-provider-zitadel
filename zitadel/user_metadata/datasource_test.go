package user_metadata_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/human_user/human_user_test_dep"
)

func TestAccUserMetadataDatasource(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_user_metadata")
	userDep, userID := human_user_test_dep.Create(t, frame)

	key := "test_key_" + frame.UniqueResourcesID
	value := "test_value"

	_, err := frame.SetUserMetadata(frame, &management.SetUserMetadataRequest{
		Id:    userID,
		Key:   key,
		Value: []byte(value),
	})
	if err != nil {
		t.Fatal(err)
	}

	config := fmt.Sprintf(`
data "zitadel_user_metadata" "default" {
  org_id  = "%s"
  user_id = "%s"
  key     = "%s"
}
`, frame.OrgID, userID, key)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, userDep},
		checkMetadataExists(frame, userID, key, value),
		map[string]string{
			"value": value,
		},
	)
}

func TestAccUserMetadatasDatasource_All(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_user_metadatas")
	userDep, userID := human_user_test_dep.Create(t, frame)

	keys := []string{"key1_" + frame.UniqueResourcesID, "key2_" + frame.UniqueResourcesID, "key3_" + frame.UniqueResourcesID}
	for _, key := range keys {
		_, err := frame.SetUserMetadata(frame, &management.SetUserMetadataRequest{
			Id:    userID,
			Key:   key,
			Value: []byte("value"),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	config := fmt.Sprintf(`
data "zitadel_user_metadatas" "default" {
  org_id  = "%s"
  user_id = "%s"
}
`, frame.OrgID, userID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, userDep},
		nil,
		map[string]string{
			"metadata.#": "3",
		},
	)
}

func TestAccUserMetadatasDatasource_Filtered(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_user_metadatas")
	userDep, userID := human_user_test_dep.Create(t, frame)

	matchingKey := "matching_" + frame.UniqueResourcesID
	keys := []string{matchingKey, "other1_" + frame.UniqueResourcesID, "other2_" + frame.UniqueResourcesID}
	for _, key := range keys {
		_, err := frame.SetUserMetadata(frame, &management.SetUserMetadataRequest{
			Id:    userID,
			Key:   key,
			Value: []byte("value"),
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	config := fmt.Sprintf(`
data "zitadel_user_metadatas" "default" {
  org_id  = "%s"
  user_id = "%s"
  key     = "%s"
}
`, frame.OrgID, userID, matchingKey)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, userDep},
		nil,
		map[string]string{
			"metadata.#": "1",
		},
	)
}

func checkMetadataExists(frame *test_utils.OrgTestFrame, userID, key, expectedValue string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		resp, err := frame.GetUserMetadata(frame, &management.GetUserMetadataRequest{
			Id:  userID,
			Key: key,
		})
		if err != nil {
			return err
		}

		actual := string(resp.GetMetadata().GetValue())
		if expectedValue != actual {
			return fmt.Errorf("expected value %s, but got %s", expectedValue, actual)
		}
		return nil
	}
}
