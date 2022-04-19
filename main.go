package main

import (
	"github.com/turbot/steampipe-plugin-cloudflare/cloudflare"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: cloudflare.Plugin})
}
