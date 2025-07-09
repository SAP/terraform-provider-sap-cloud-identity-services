//go:generate go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
//go:generate tfplugindocs generate --provider-name "sci" --rendered-provider-name "SAP Cloud Identity Services"

package main

import (
	"context"
	"flag"
	"log"

	"github.com/SAP/terraform-provider-sap-cloud-identity-services/sci/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

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
		Address:         "registry.terraform.io/sap/sap-cloud-identity-services",
		Debug:           debug,
		ProtocolVersion: 6,
	}
	err := providerserver.Serve(context.Background(), provider.New, opts)
	if err != nil {
		log.Fatal(err.Error())
	}

}
