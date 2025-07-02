#!/bin/bash

echo "=== Fixing ZITADEL Official Repo ==="

# Step 1: Update dependencies
echo "1. Updating dependencies..."
go get github.com/hashicorp/terraform-plugin-framework@v1.15.0
go get github.com/hashicorp/terraform-plugin-mux@v0.18.0  
go get github.com/hashicorp/terraform-plugin-sdk/v2@v2.33.0
go mod tidy

# Step 2: Fix generated file
echo "2. Fixing generated file..."
GENFILE="gen/github.com/zitadel/zitadel/pkg/grpc/text/text_terraform.go"
if [ -f "$GENFILE" ]; then
    cp "$GENFILE" "$GENFILE.backup"
    sed -i 's|github.com/hashicorp/terraform-plugin-framework/provider/schema|github.com/hashicorp/terraform-plugin-framework/provider/schema|g' "$GENFILE"
    sed -i 's/github_com_hashicorp_terraform_plugin_framework_tfsdk/github_com_hashicorp_terraform_plugin_framework_provider_schema/g' "$GENFILE"
    sed -i 's/tfsdk\./schema./g' "$GENFILE"
    echo "Generated file fixed!"
else
    echo "Generated file not found, skipping..."
fi

# Step 3: Test build
echo "3. Testing build..."
go build . 2>&1 | head -20

echo "=== Done! Check output above ==="
