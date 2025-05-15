package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hirosi1900day/terraform-provider-devin-knowledge/internal/provider"
)

// Version information will be set at build time
var (
	version string = "0.0.2"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "Run in debug mode")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/hirosi1900day/devin",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
