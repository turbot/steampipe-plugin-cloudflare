package cloudflare

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-cloudflare",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
		DefaultTransform: transform.FromJSONTag(),
		TableMap: map[string]*plugin.Table{
			"cloudflare_account":       tableCloudflareAccount(ctx),
			"cloudflare_account_role":  tableCloudflareAccountRole(ctx),
			"cloudflare_api_token":     tableCloudflareAPIToken(ctx),
			"cloudflare_dns_record":    tableCloudflareDNSRecord(ctx),
			"cloudflare_firewall_rule": tableCloudflareFirewallRule(ctx),
			"cloudflare_user":          tableCloudflareUser(ctx),
			"cloudflare_zone":          tableCloudflareZone(ctx),
		},
	}
	return p
}