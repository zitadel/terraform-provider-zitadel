package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"

	"github.com/zitadel/terraform-provider-zitadel/zitadel"
)

func main() {
	ctx := context.Background()
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	upgradedSdkProvider, err := tf5to6server.UpgradeServer(ctx, zitadel.Provider().GRPCProvider)
	if err != nil {
		log.Fatal(err)
	}

	providers := []func() tfprotov6.ProviderServer{
		func() tfprotov6.ProviderServer {
			return upgradedSdkProvider
		},
		providerserver.NewProtocol6(zitadel.NewProviderPV6()),
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)

	opts := []tf6server.ServeOpt{}
	if debug {
		opts = append(opts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve("registry.terraform.io/zitadel/zitadel", muxServer.ProviderServer, opts...)

	if err != nil {
		log.Fatalln(err.Error())
	}
}
