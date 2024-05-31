package main

import (
	// "cloudIdentityServices/internal/cli"
	"terraform-provider-cloudidentityservices/internal/provider"
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	// "cloudIdentityServices/cli/apiObjects/users"
	// "net/http"
	// "net/url"
	// "context"
	// "encoding/json"
	// "fmt"
)

// const host = "https://iasprovidertestblr.accounts400.ondemand.com/"

func main () {
	
	var debug bool
    flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
    flag.Parse()
    opts := providerserver.ServeOpts{
        // NOTE: This is not a typical Terraform Registry provider address,
        // such as registry.terraform.io/hashicorp/hashicups. This specific
        // provider address is used in these tutorials in conjunction with a
        // specific Terraform CLI configuration for manual development testing
        // of this provider.
        Address: "registry.terraform.io/sap/cloudIdentityServices",
        Debug:   debug,
		ProtocolVersion: 6,
    }
    err := providerserver.Serve(context.Background(), provider.New, opts)
    if err != nil {
        log.Fatal(err.Error())
    }

}