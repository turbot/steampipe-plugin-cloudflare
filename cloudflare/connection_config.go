package cloudflare

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type cloudflareConfig struct {
	Token     *string `hcl:"token"`
	Email     *string `hcl:"email"`
	APIKey    *string `hcl:"api_key"`
	AccessKey *string `hcl:"access_key"`
	SecretKey *string `hcl:"secret_key"`
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
