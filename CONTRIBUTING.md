# Debug

1. Run the local plugin code with your debugging IDE of choice with something similar to `go run ./... -debug`.
2. Set breakpoints in your IDE.
3. In your shell, apply the resource you are working on.
   ```
   # export the printed environment variable from the go run ./... -debug command above. E.g.
   export TF_REATTACH_PROVIDERS='{"registry.terraform.io/zitadel/zitadel":{"Protocol":"grpc","ProtocolVersion":6,"Pid":8123,"Test":true,"Addr":{"Network":"unix","String":"/tmp/plugin275634719"}}}'
   
   # go to a directory containing .tf files.
   cd /my-zitadel-terraform-files
   
   # apply them
   terraform apply
   ```
4. The execution stops at your breakpoints.

# Run Acceptance Tests

```bash
TF_ACC=1 TF_ACC_ZITADEL_TOKEN=/my-token.json go test ./...
```

The tests are flaky when resources should be cleaned up.
This results in dangling resources.
