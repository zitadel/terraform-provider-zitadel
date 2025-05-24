package main

import (
	"context"
	"flag"
	"log"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/zitadel/zitadel",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), func() provider.Provider {
		return zitadel.NewProviderPV6()
	}, opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
