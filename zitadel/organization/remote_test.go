package organization_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/authn"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	org "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/org/v2"
	userv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/user/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

// TestOrgLevelPermissionsCanReadOwnOrg validates that a service account with only
// ORG_OWNER role (no instance-level permissions) can read their own organization
// using the V2 API. This test validates the fix for issue #245.
func TestOrgLevelPermissionsCanReadOwnOrg(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	frame := test_utils.NewInstanceTestFrame(t, "zitadel_organization")
	ctx := frame.Context

	// Step 1: Create a new organization
	orgClient, err := helper.GetOrgClient(ctx, frame.ClientInfo)
	if err != nil {
		t.Fatalf("Failed to get org client: %v", err)
	}

	timestamp := time.Now().Unix()
	orgName := fmt.Sprintf("test-org-perms-%d", timestamp)

	createOrgResp, err := orgClient.AddOrganization(ctx, &org.AddOrganizationRequest{
		Name: orgName,
	})
	if err != nil {
		t.Fatalf("Failed to create org: %v", err)
	}
	testOrgID := createOrgResp.OrganizationId
	t.Logf("Created org: %s (ID: %s)", orgName, testOrgID)

	defer func() {
		_, _ = orgClient.DeleteOrganization(ctx, &org.DeleteOrganizationRequest{
			OrganizationId: testOrgID,
		})
	}()

	// Step 2: Create service account using UserV2 API
	userClient, err := helper.GetUserV2Client(ctx, frame.ClientInfo)
	if err != nil {
		t.Fatalf("Failed to get user client: %v", err)
	}

	username := fmt.Sprintf("test-sa-%d", timestamp)
	createUserResp, err := userClient.CreateUser(helper.CtxSetOrgID(ctx, testOrgID), &userv2.CreateUserRequest{
		OrganizationId: testOrgID,
		Username:       &username,
		UserType: &userv2.CreateUserRequest_Machine_{
			Machine: &userv2.CreateUserRequest_Machine{
				Name: "Test Service Account",
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create service account: %v", err)
	}
	saUserID := createUserResp.Id
	t.Logf("Created service account: %s", saUserID)

	// Step 3: Grant ORG_OWNER role using Management API
	mgmtClient, err := helper.GetManagementClient(helper.CtxSetOrgID(ctx, testOrgID), frame.ClientInfo)
	if err != nil {
		t.Fatalf("Failed to get management client: %v", err)
	}

	_, err = mgmtClient.AddOrgMember(helper.CtxSetOrgID(ctx, testOrgID), &management.AddOrgMemberRequest{
		UserId: saUserID,
		Roles:  []string{"ORG_OWNER"},
	})
	if err != nil {
		t.Fatalf("Failed to grant ORG_OWNER role: %v", err)
	}
	t.Logf("Granted ORG_OWNER role")

	// Step 4: Create key using Management API
	addKeyResp, err := mgmtClient.AddMachineKey(helper.CtxSetOrgID(ctx, testOrgID), &management.AddMachineKeyRequest{
		UserId: saUserID,
		Type:   authn.KeyType_KEY_TYPE_JSON,
	})
	if err != nil {
		t.Fatalf("Failed to create key: %v", err)
	}
	saKeyJSON := string(addKeyResp.KeyDetails)
	t.Logf("Created service account key")

	// Step 5: Create ClientInfo for org-level SA
	orgLevelClientInfo, err := helper.GetClientInfo(ctx, false, frame.ClientInfo.Domain, "", "", "", "", saKeyJSON, "", false, nil)
	if err != nil {
		t.Fatalf("Failed to create org-level client info: %v", err)
	}

	// Step 6: Create org client using org-level credentials
	orgLevelOrgClient, err := helper.GetOrgClient(ctx, orgLevelClientInfo)
	if err != nil {
		t.Fatalf("Failed to create org-level org client: %v", err)
	}

	// Step 7: Try to read org with org-level credentials
	t.Log("Attempting to read org using org-level credentials...")
	listResp, err := orgLevelOrgClient.ListOrganizations(ctx, &org.ListOrganizationsRequest{
		Queries: []*org.SearchQuery{
			{
				Query: &org.SearchQuery_IdQuery{
					IdQuery: &org.OrganizationIDQuery{Id: testOrgID},
				},
			},
		},
	})

	if err != nil {
		t.Fatalf("❌ FAILED: %v\nIssue #245 is NOT fixed.", err)
	}

	if len(listResp.Result) == 0 {
		t.Fatal("❌ FAILED: No orgs returned\nIssue #245 is NOT fixed.")
	}

	t.Logf("✅ SUCCESS: Org-level SA can read their own org")
	t.Logf("   ID: %s", listResp.Result[0].Id)
	t.Logf("   Name: %s", listResp.Result[0].Name)
	t.Log("✅ Issue #245 IS FIXED by V2 migration (#305)")
}

// TestOrgLevelPermissionsCannotReadOwnOrgLegacy validates that the legacy org resource
// (using Admin API) CANNOT be read by org-level permissions, confirming the original issue #245.
func TestOrgLevelPermissionsCannotReadOwnOrgLegacy(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	frame := test_utils.NewInstanceTestFrame(t, "zitadel_org")
	ctx := frame.Context

	// Step 1: Create a new organization
	orgClient, err := helper.GetOrgClient(ctx, frame.ClientInfo)
	if err != nil {
		t.Fatalf("Failed to get org client: %v", err)
	}

	timestamp := time.Now().Unix()
	orgName := fmt.Sprintf("test-legacy-perms-%d", timestamp)

	createOrgResp, err := orgClient.AddOrganization(ctx, &org.AddOrganizationRequest{
		Name: orgName,
	})
	if err != nil {
		t.Fatalf("Failed to create org: %v", err)
	}
	testOrgID := createOrgResp.OrganizationId
	t.Logf("Created org: %s (ID: %s)", orgName, testOrgID)

	defer func() {
		_, _ = orgClient.DeleteOrganization(ctx, &org.DeleteOrganizationRequest{
			OrganizationId: testOrgID,
		})
	}()

	// Step 2: Create service account
	userClient, err := helper.GetUserV2Client(ctx, frame.ClientInfo)
	if err != nil {
		t.Fatalf("Failed to get user client: %v", err)
	}

	username := fmt.Sprintf("test-legacy-sa-%d", timestamp)
	createUserResp, err := userClient.CreateUser(helper.CtxSetOrgID(ctx, testOrgID), &userv2.CreateUserRequest{
		OrganizationId: testOrgID,
		Username:       &username,
		UserType: &userv2.CreateUserRequest_Machine_{
			Machine: &userv2.CreateUserRequest_Machine{
				Name: "Test Legacy Service Account",
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create service account: %v", err)
	}
	saUserID := createUserResp.Id
	t.Logf("Created service account: %s", saUserID)

	// Step 3: Grant ORG_OWNER role (no instance-level permissions)
	mgmtClient, err := helper.GetManagementClient(helper.CtxSetOrgID(ctx, testOrgID), frame.ClientInfo)
	if err != nil {
		t.Fatalf("Failed to get management client: %v", err)
	}

	_, err = mgmtClient.AddOrgMember(helper.CtxSetOrgID(ctx, testOrgID), &management.AddOrgMemberRequest{
		UserId: saUserID,
		Roles:  []string{"ORG_OWNER"},
	})
	if err != nil {
		t.Fatalf("Failed to grant ORG_OWNER role: %v", err)
	}
	t.Logf("Granted ORG_OWNER role")

	// Step 4: Create key
	addKeyResp, err := mgmtClient.AddMachineKey(helper.CtxSetOrgID(ctx, testOrgID), &management.AddMachineKeyRequest{
		UserId: saUserID,
		Type:   authn.KeyType_KEY_TYPE_JSON,
	})
	if err != nil {
		t.Fatalf("Failed to create key: %v", err)
	}
	saKeyJSON := string(addKeyResp.KeyDetails)
	t.Logf("Created service account key")

	// Step 5: Create ClientInfo for org-level SA
	orgLevelClientInfo, err := helper.GetClientInfo(ctx, false, frame.ClientInfo.Domain, "", "", "", "", saKeyJSON, "", false, nil)
	if err != nil {
		t.Fatalf("Failed to create org-level client info: %v", err)
	}

	// Step 6: Try to read org using LEGACY Admin API (should fail)
	adminClient, err := helper.GetAdminClient(ctx, orgLevelClientInfo)
	if err != nil {
		t.Fatalf("Failed to create admin client: %v", err)
	}

	t.Log("Attempting to read org using legacy Admin API with org-level credentials...")
	_, err = adminClient.GetOrgByID(ctx, &admin.GetOrgByIDRequest{
		Id: testOrgID,
	})

	if err == nil {
		t.Fatal("❌ UNEXPECTED: Legacy Admin API allowed org-level SA to read org\nThis contradicts the original issue #245")
	}

	// Verify it's a permission denied error
	if !strings.Contains(err.Error(), "PermissionDenied") && !strings.Contains(err.Error(), "permission") {
		t.Fatalf("Expected PermissionDenied error, got: %v", err)
	}

	t.Logf("✅ CONFIRMED: Legacy Admin API correctly denies org-level SA from reading org")
	t.Logf("   Error: %v", err)
	t.Log("✅ This confirms the original issue #245 existed in the legacy 'org' resource")
}

// TestOrgLevelPermissionsCannotGetDefaultOrg validates that org-level permissions
// cannot call GetDefaultOrg, which is required by the legacy org resource read operation.
// This confirms the root cause of issue #245.
func TestOrgLevelPermissionsCannotGetDefaultOrg(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	frame := test_utils.NewInstanceTestFrame(t, "zitadel_org")
	ctx := frame.Context

	// Create org and service account with ORG_OWNER (same setup as before)
	orgClient, err := helper.GetOrgClient(ctx, frame.ClientInfo)
	if err != nil {
		t.Fatalf("Failed to get org client: %v", err)
	}

	timestamp := time.Now().Unix()
	orgName := fmt.Sprintf("test-default-org-%d", timestamp)

	createOrgResp, err := orgClient.AddOrganization(ctx, &org.AddOrganizationRequest{
		Name: orgName,
	})
	if err != nil {
		t.Fatalf("Failed to create org: %v", err)
	}
	testOrgID := createOrgResp.OrganizationId
	t.Logf("Created org: %s (ID: %s)", orgName, testOrgID)

	defer func() {
		_, _ = orgClient.DeleteOrganization(ctx, &org.DeleteOrganizationRequest{
			OrganizationId: testOrgID,
		})
	}()

	userClient, err := helper.GetUserV2Client(ctx, frame.ClientInfo)
	if err != nil {
		t.Fatalf("Failed to get user client: %v", err)
	}

	username := fmt.Sprintf("test-default-sa-%d", timestamp)
	createUserResp, err := userClient.CreateUser(helper.CtxSetOrgID(ctx, testOrgID), &userv2.CreateUserRequest{
		OrganizationId: testOrgID,
		Username:       &username,
		UserType: &userv2.CreateUserRequest_Machine_{
			Machine: &userv2.CreateUserRequest_Machine{
				Name: "Test Service Account",
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create service account: %v", err)
	}
	saUserID := createUserResp.Id

	mgmtClient, err := helper.GetManagementClient(helper.CtxSetOrgID(ctx, testOrgID), frame.ClientInfo)
	if err != nil {
		t.Fatalf("Failed to get management client: %v", err)
	}

	_, err = mgmtClient.AddOrgMember(helper.CtxSetOrgID(ctx, testOrgID), &management.AddOrgMemberRequest{
		UserId: saUserID,
		Roles:  []string{"ORG_OWNER"},
	})
	if err != nil {
		t.Fatalf("Failed to grant ORG_OWNER role: %v", err)
	}

	addKeyResp, err := mgmtClient.AddMachineKey(helper.CtxSetOrgID(ctx, testOrgID), &management.AddMachineKeyRequest{
		UserId: saUserID,
		Type:   authn.KeyType_KEY_TYPE_JSON,
	})
	if err != nil {
		t.Fatalf("Failed to create key: %v", err)
	}
	saKeyJSON := string(addKeyResp.KeyDetails)

	orgLevelClientInfo, err := helper.GetClientInfo(ctx, false, frame.ClientInfo.Domain, "", "", "", "", saKeyJSON, "", false, nil)
	if err != nil {
		t.Fatalf("Failed to create org-level client info: %v", err)
	}

	adminClient, err := helper.GetAdminClient(ctx, orgLevelClientInfo)
	if err != nil {
		t.Fatalf("Failed to create admin client: %v", err)
	}

	// THIS is what the legacy org resource does during read - and this should fail
	t.Log("Attempting to call GetDefaultOrg (required by legacy org resource)...")
	_, err = adminClient.GetDefaultOrg(ctx, &admin.GetDefaultOrgRequest{})

	if err == nil {
		t.Fatal("❌ UNEXPECTED: GetDefaultOrg succeeded with org-level permissions")
	}

	if !strings.Contains(err.Error(), "PermissionDenied") && !strings.Contains(err.Error(), "permission") {
		t.Logf("Got error (not permission-related): %v", err)
		t.Fatal("Expected PermissionDenied error for GetDefaultOrg")
	}

	t.Logf("✅ CONFIRMED: GetDefaultOrg fails with org-level permissions")
	t.Logf("   Error: %v", err)
	t.Log("✅ This is the root cause of issue #245 - legacy org resource calls GetDefaultOrg during read")
}

// TestIamOrgManagerCanDestroyOrg validates the exact scenario from issue #245:
// 1. Service account with IAM_ORG_MANAGER and IAM_USER_MANAGER (instance-level roles)
// 2. Creates an org via Terraform
// 3. Attempts to destroy it (which requires reading it first)
func TestIamOrgManagerCanDestroyOrg(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	frame := test_utils.NewInstanceTestFrame(t, "zitadel_organization")
	ctx := frame.Context

	// Get the instance's default org ID
	adminClient, err := helper.GetAdminClient(ctx, frame.ClientInfo)
	if err != nil {
		t.Fatalf("Failed to get admin client: %v", err)
	}

	defaultOrgResp, err := adminClient.GetDefaultOrg(ctx, &admin.GetDefaultOrgRequest{})
	if err != nil {
		t.Fatalf("Failed to get default org: %v", err)
	}
	defaultOrgID := defaultOrgResp.Org.Id
	t.Logf("Default org ID: %s", defaultOrgID)

	// Step 1: Create service account with IAM_ORG_MANAGER and IAM_USER_MANAGER
	userClient, err := helper.GetUserV2Client(ctx, frame.ClientInfo)
	if err != nil {
		t.Fatalf("Failed to get user client: %v", err)
	}

	timestamp := time.Now().Unix()
	username := fmt.Sprintf("test-iam-mgr-%d", timestamp)

	// Create service account in the default org
	createUserResp, err := userClient.CreateUser(helper.CtxSetOrgID(ctx, defaultOrgID), &userv2.CreateUserRequest{
		OrganizationId: defaultOrgID,
		Username:       &username,
		UserType: &userv2.CreateUserRequest_Machine_{
			Machine: &userv2.CreateUserRequest_Machine{
				Name: "Test IAM Manager",
			},
		},
	})
	if err != nil {
		t.Fatalf("Failed to create service account: %v", err)
	}
	saUserID := createUserResp.Id
	t.Logf("Created service account: %s", saUserID)

	// Grant IAM_ORG_MANAGER and IAM_USER_MANAGER roles (instance-level)
	_, err = adminClient.AddIAMMember(ctx, &admin.AddIAMMemberRequest{
		UserId: saUserID,
		Roles:  []string{"IAM_ORG_MANAGER", "IAM_USER_MANAGER"},
	})
	if err != nil {
		t.Fatalf("Failed to grant IAM roles: %v", err)
	}
	t.Logf("Granted IAM_ORG_MANAGER and IAM_USER_MANAGER roles")

	// Create key for the service account
	mgmtClient, err := helper.GetManagementClient(helper.CtxSetOrgID(ctx, defaultOrgID), frame.ClientInfo)
	if err != nil {
		t.Fatalf("Failed to get management client: %v", err)
	}

	addKeyResp, err := mgmtClient.AddMachineKey(helper.CtxSetOrgID(ctx, defaultOrgID), &management.AddMachineKeyRequest{
		UserId: saUserID,
		Type:   authn.KeyType_KEY_TYPE_JSON,
	})
	if err != nil {
		t.Fatalf("Failed to create key: %v", err)
	}
	saKeyJSON := string(addKeyResp.KeyDetails)
	t.Logf("Created service account key")

	// Step 2: Create ClientInfo for IAM manager SA
	// Extract port from frame's domain if present
	domain := frame.ClientInfo.Domain
	port := ""
	if idx := strings.LastIndex(domain, ":"); idx != -1 {
		port = domain[idx+1:]
		domain = domain[:idx]
	}

	iamMgrClientInfo, err := helper.GetClientInfo(
		ctx,
		true, // insecure - matches test environment
		domain,
		"", "", "", "",
		saKeyJSON,
		port,
		false,
		nil,
	)
	if err != nil {
		t.Fatalf("Failed to create IAM manager client info: %v", err)
	}

	// Step 3: Create an org as the IAM_ORG_MANAGER
	iamOrgClient, err := helper.GetOrgClient(ctx, iamMgrClientInfo)
	if err != nil {
		t.Fatalf("Failed to create org client: %v", err)
	}

	orgName := fmt.Sprintf("test-iam-org-%d", timestamp)
	createOrgResp, err := iamOrgClient.AddOrganization(ctx, &org.AddOrganizationRequest{
		Name: orgName,
	})
	if err != nil {
		t.Fatalf("Failed to create org as IAM_ORG_MANAGER: %v", err)
	}
	testOrgID := createOrgResp.OrganizationId
	t.Logf("Created org as IAM_ORG_MANAGER: %s (ID: %s)", orgName, testOrgID)

	// Step 4: Try to READ the org (required for Terraform destroy/refresh)
	t.Log("Attempting to read org as IAM_ORG_MANAGER (simulating terraform destroy)...")
	listResp, err := iamOrgClient.ListOrganizations(ctx, &org.ListOrganizationsRequest{
		Queries: []*org.SearchQuery{
			{
				Query: &org.SearchQuery_IdQuery{
					IdQuery: &org.OrganizationIDQuery{Id: testOrgID},
				},
			},
		},
	})

	if err != nil {
		t.Logf("❌ FAILED: IAM_ORG_MANAGER cannot read org: %v", err)
		t.Log("This reproduces issue #245 - IAM_ORG_MANAGER can CREATE but not READ orgs")

		// Cleanup with admin credentials
		cleanupClient, _ := helper.GetOrgClient(ctx, frame.ClientInfo)
		_, _ = cleanupClient.DeleteOrganization(ctx, &org.DeleteOrganizationRequest{
			OrganizationId: testOrgID,
		})

		t.FailNow()
	}

	if len(listResp.Result) == 0 {
		t.Log("❌ CONFIRMED: IAM_ORG_MANAGER cannot read orgs they created")
		t.Log("This reproduces issue #245 exactly as reported")
		t.Log("The user could CREATE orgs but terraform destroy fails because it can't READ them")

		// Cleanup with admin credentials
		cleanupClient, _ := helper.GetOrgClient(ctx, frame.ClientInfo)
		_, _ = cleanupClient.DeleteOrganization(ctx, &org.DeleteOrganizationRequest{
			OrganizationId: testOrgID,
		})

		t.Fatal("Issue #245 is NOT fixed - IAM_ORG_MANAGER still cannot read/destroy orgs")
	}

	t.Logf("✅ SUCCESS: IAM_ORG_MANAGER can read the org they created")
	t.Logf("   ID: %s", listResp.Result[0].Id)
	t.Logf("   Name: %s", listResp.Result[0].Name)

	// Step 5: Try to DELETE (complete the destroy operation)
	_, err = iamOrgClient.DeleteOrganization(ctx, &org.DeleteOrganizationRequest{
		OrganizationId: testOrgID,
	})
	if err != nil {
		t.Fatalf("Failed to delete org: %v", err)
	}

	t.Log("✅ SUCCESS: Full lifecycle works - IAM_ORG_MANAGER can create, read, and delete orgs")
	t.Log("✅ Issue #245 is resolved for the reported scenario")
}
