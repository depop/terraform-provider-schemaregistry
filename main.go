package main

import (
	schemaregisty "github.com/depop/terraform-provider-schemaregistry/schemaregistry"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: schemaregisty.Provider})
}
