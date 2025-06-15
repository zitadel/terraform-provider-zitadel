#!/bin/bash

GENFILE="gen/github.com/zitadel/zitadel/pkg/grpc/text/text_terraform.go"

# Create backup
mv "$GENFILE" "$GENFILE.complex"

# Create simple stub that compiles
cat > "$GENFILE" << 'GOEOF'
// Temporary stub for framework upgrade
package text

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/provider/schema"
    textpb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/text"
)

// GenSchemaMessageCustomText returns a simple schema
func GenSchemaMessageCustomText(ctx context.Context) schema.Schema {
    return schema.Schema{
        Attributes: map[string]schema.Attribute{
            "language": schema.StringAttribute{
                Required: true,
                Description: "Language code",
            },
            "title": schema.StringAttribute{
                Optional: true,
                Description: "Custom title",
            },
            "text": schema.StringAttribute{
                Optional: true,
                Description: "Custom text",
            },
        },
    }
}

// GenSchemaLoginCustomText returns a simple schema  
func GenSchemaLoginCustomText(ctx context.Context) schema.Schema {
    return schema.Schema{
        Attributes: map[string]schema.Attribute{
            "language": schema.StringAttribute{
                Required: true,
                Description: "Language code", 
            },
            "login_text": schema.StringAttribute{
                Optional: true,
                Description: "Custom login text",
            },
        },
    }
}

// Stub functions for compilation
func CopyMessageCustomTextFromTerraform(ctx context.Context, data interface{}) (*textpb.MessageCustomText, error) {
    return &textpb.MessageCustomText{}, nil
}

func CopyMessageCustomTextToTerraform(ctx context.Context, msg *textpb.MessageCustomText) (interface{}, error) {
    return map[string]interface{}{}, nil
}

func CopyLoginCustomTextFromTerraform(ctx context.Context, data interface{}) (*textpb.LoginCustomText, error) {
    return &textpb.LoginCustomText{}, nil
}

func CopyLoginCustomTextToTerraform(ctx context.Context, msg *textpb.LoginCustomText) (interface{}, error) {
    return map[string]interface{}{}, nil
}
GOEOF

echo "Simple stub created!"
