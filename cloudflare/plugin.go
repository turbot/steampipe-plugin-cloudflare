package cloudflare

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name: "steampipe-plugin-cloudflare",
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
		},
		ConnectionKeyColumns: []plugin.ConnectionKeyColumn{
			{
				Name:    "user_id",
				Hydrate: getUserId,
			},
		},
		DefaultTransform: transform.FromJSONTag(),
		TableMap: map[string]*plugin.Table{
			"cloudflare_access_application":    tableCloudflareAccessApplication(ctx),
			"cloudflare_access_group":          tableCloudflareAccessGroup(ctx),
			"cloudflare_access_policy":         tableCloudflareAccessPolicy(ctx),
			"cloudflare_account":               tableCloudflareAccount(ctx),
			"cloudflare_account_member":        tableCloudflareAccountMember(),
			"cloudflare_account_role":          tableCloudflareAccountRole(ctx),
			"cloudflare_api_token":             tableCloudflareAPIToken(ctx),
			"cloudflare_dns_record":            tableCloudflareDNSRecord(),
			"cloudflare_firewall_rule":         tableCloudflareFirewallRule(ctx),
			"cloudflare_load_balancer":         tableCloudflareLoadBalancer(),
			"cloudflare_load_balancer_monitor": tableCloudflareLoadBalancerMonitor(ctx),
			"cloudflare_load_balancer_pool":    tableCloudflareLoadBalancerPool(ctx),
			"cloudflare_page_rule":             tableCloudflarePageRule(),
			"cloudflare_r2_bucket":             tableCloudflareR2Bucket(ctx),
			"cloudflare_r2_object":             tableCloudflareR2Object(ctx),
			"cloudflare_r2_object_data":        tableCloudflareR2ObjectData(ctx),
			"cloudflare_user":                  tableCloudflareUser(ctx),
			"cloudflare_user_audit_log":        tableCloudflareUserAuditLog(),
			"cloudflare_worker_route":          tableCloudflareWorkerRoute(),
			"cloudflare_zone":                  TableCloudflareZone(),
		},
	}
	return p
}
