// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/plain-insure/terraform-provider-cronjoborg/provider"
)

// Run "go generate" in the tools directory to format example terraform files and generate the docs for the registry/website
// Use `go generate ./tools` or `make generate` to run the generation process.

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev" //nolint:unused // Set by build via ldflags

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		ProviderFunc: provider.Provider,
		// TODO: update this string with the full name of your provider as used in your configs
		ProviderAddr: "registry.terraform.io/plain-insure/cronjoborg",
		Debug:        debugMode,
	}

	plugin.Serve(opts)
}
