package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hirosi1900day/terraform-provider-devin-knowledge/internal/provider"
)

// バージョン情報はビルド時に設定されます
var (
	version string = "0.0.2"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "デバッグモードで実行する")
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
