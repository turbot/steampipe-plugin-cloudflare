package cloudflare

import (
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/schema"
)

type cloudflareConfig struct {
	Token  *string `cty:"token"`
	Email  *string `cty:"email"`
	APIKey *string `cty:"api_key"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"token": {
		Type: schema.TypeString,
	},
	"email": {
		Type: schema.TypeString,
	},
	"api_key": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &cloudflareConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) cloudflareConfig {
	if connection == nil || connection.Config == nil {
		return cloudflareConfig{}
	}
	config, _ := connection.Config.(cloudflareConfig)
	return config
}
