//go:generate go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
//go:generate tfplugindocs generate --provider-name "sci" --rendered-provider-name "SAP Cloud Identity Services"

package main

import (
	// "sci/internal/cli"
	"context"
	"flag"
	"log"
	"terraform-provider-sci/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	// "sci/cli/apiObjects/users"
	// "net/http"
	// "net/url"
	// "context"
	// "encoding/json"
	// "fmt"
)

// const host = "https://iasprovidertestblr.accounts400.ondemand.com/"

func main() {

	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()
	opts := providerserver.ServeOpts{
		// NOTE: This is not a typical Terraform Registry provider address,
		// such as registry.terraform.io/hashicorp/hashicups. This specific
		// provider address is used in these tutorials in conjunction with a
		// specific Terraform CLI configuration for manual development testing
		// of this provider.
		Address:         "registry.terraform.io/sap/sci",
		Debug:           debug,
		ProtocolVersion: 6,
	}
	err := providerserver.Serve(context.Background(), provider.New, opts)
	if err != nil {
		log.Fatal(err.Error())
	}

}
