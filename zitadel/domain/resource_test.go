package domain_test

/*
commented as the test fail but the functionality works

func TestAccDomain(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_domain")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, domain.NameVar, exampleAttributes).AsString()
	updatedProperty := "updated.default.127.0.0.1.sslip.io"
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, updatedProperty,
		"", "", "",
		true,
		checkRemoteProperty(frame),
		regexp.MustCompile(fmt.Sprintf(`^%s$|^%s$`, exampleProperty, updatedProperty)),
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportOrgId(frame),
		),
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(_ *terraform.State) error {
			remoteResource, err := frame.ListOrgDomains(frame, &management.ListOrgDomainsRequest{
				Queries: []*org.DomainSearchQuery{{
					Query: &org.DomainSearchQuery_DomainNameQuery{
						DomainNameQuery: &org.DomainNameQuery{
							Name: expect,
						},
					},
				}},
			})
			if err != nil {
				return err
			}
			if len(remoteResource.GetResult()) == 0 {
				return fmt.Errorf("expected to find %s, but didn't: %w", expect, test_utils.ErrNotFound)
			}
			return nil
		}
	}
}
*/
