# Live binding acceptance tests

`TestAccAuthorizationV2` covers human, machine and delegated role assignments:
create, complete role replacement without identity replacement, import,
server-side drift correction and remote deletion. `TestAccAdministratorV2`
covers create, full role replacement, import and remote deletion for instance,
organization, project and project-grant scopes. Every check reads the real
Authorization/InternalPermission v2 API; these are not mock-server tests.

The existing pull-request acceptance job provisions Zitadel with
`acceptance/docker-compose.yaml` and runs these cases with the rest of
`TF_ACC=1 go test ./...`. They use its instance-level signing-key fixture by
default. In-process gRPC fault tests remain in the resource packages.

```sh
TF_ACC=1 go test ./zitadel/bindings -v -count=1
```

To target an already-running disposable instance instead, set
`ZITADEL_TEST_DOMAIN`, `ZITADEL_TEST_ACCESS_TOKEN`, `ZITADEL_TEST_PORT`
(default 443), `ZITADEL_TEST_INSECURE` (default false), and optionally
`ZITADEL_TEST_TRANSPORT_HEADERS` (a JSON object). The PAT is passed to Terraform
through a sensitive input variable, not inserted into the HCL. Do not enable
Terraform debug logging or publish test state/credential files.

The caller must be allowed to create/delete organizations, projects, users,
roles, delegation and all four administrator scopes. Tests create unique,
disposable organizations and verify remote assignment cleanup on destroy.
They do not adopt or modify pre-existing application identities.
