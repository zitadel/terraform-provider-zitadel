protoc \
    -I$(pwd) \
    -I$(go env GOPATH)/src/github.com/gogo/protobuf \
    -I$(go env GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway \
    -I$(go env GOPATH)/src/github.com/envoyproxy/protoc-gen-validate \
    -I$(go env GOPATH)/src/github.com/zitadel/zitadel/proto \
    --plugin=$(go env GOBIN)/protoc-gen-terraform \
    --terraform_out=config=gen/config.yaml:gen \
    $(go env GOPATH)/src/github.com/zitadel/zitadel/proto/zitadel/text.proto

sed -i 's#_ "github.com/zitadel/zitadel/pkg/grpc/object"##g' gen/github.com/zitadel/zitadel/pkg/grpc/text/text_terraform.go
sed -i 's#textpb "textpb"#textpb "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/text"#g' gen/github.com/zitadel/zitadel/pkg/grpc/text/text_terraform.go
sed -i 's/U2f/U2F/g' gen/github.com/zitadel/zitadel/pkg/grpc/text/text_terraform.go

