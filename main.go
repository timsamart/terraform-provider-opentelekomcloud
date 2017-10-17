package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/gator1/terraform-provider-opentelekomcloud/opentelekomcloud" // TODO: Revert path when merge
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: opentelekomcloud.Provider})
}
