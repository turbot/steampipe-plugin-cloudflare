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
			"cloudflare_account":               tableCloudflareAccount(ctx),
			"cloudflare_account_role":          tableCloudflareAccountRole(ctx),
			"cloudflare_api_token":             tableCloudflareAPIToken(ctx),
			"cloudflare_dns_record":            tableCloudflareDNSRecord(ctx),
			"cloudflare_firewall_rule":         tableCloudflareFirewallRule(ctx),
			"cloudflare_load_balancer":         tableCloudflareLoadBalancer(ctx),
			"cloudflare_load_balancer_monitor": tableCloudflareLoadBalancerMonitor(ctx),
			"cloudflare_load_balancer_pool":    tableCloudflareLoadBalancerPool(ctx),
			"cloudflare_page_rule":             tableCloudflarePageRule(ctx),
			"cloudflare_user":                  tableCloudflareUser(ctx),
			"cloudflare_worker_route":          tableCloudflareWorkerRoute(ctx),
			"cloudflare_worker_script":         tableCloudflareWorkerScript(ctx),
			"cloudflare_workers_kv_namespace":  tableCloudflareWorkersKVNamespace(ctx),
			"cloudflare_zone":                  tableCloudflareZone(ctx),
		},
	}
	return p
}
