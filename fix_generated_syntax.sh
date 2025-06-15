#!/bin/bash

GENFILE="gen/github.com/zitadel/zitadel/pkg/grpc/text/text_terraform.go"

echo "Fixing generated file syntax..."

# Fix the main issues in generated file
sed -i 's/schema\.Attribute{/schema.StringAttribute{/g' "$GENFILE"
sed -i 's/github_com_hashicorp_terraform_plugin_framework_provider_schema\.SingleNestedAttributes/schema.SingleNestedAttribute/g' "$GENFILE"
sed -i 's/github_com_hashicorp_terraform_plugin_framework_provider_schema\.Attribute/schema.StringAttribute/g' "$GENFILE"

# Fix type declarations
sed -i 's/map\[string\]schema\.StringAttribute/map[string]schema.Attribute/g' "$GENFILE"

echo "Generated file syntax fixed!"
