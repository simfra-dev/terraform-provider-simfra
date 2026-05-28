package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/simfra-dev/terraform-provider-simfra/internal/provider"
	"github.com/simfra-dev/terraform-provider-simfra/version"
)

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name simfra

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/simfra-dev/simfra",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version.ProviderVersion), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
