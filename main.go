package main

import (
	"github.com/turbot/steampipe-plugin-cloudflare/cloudflare"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: cloudflare.Plugin})
}
