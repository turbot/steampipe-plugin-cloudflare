package cloudflare

import (
	"context"
	"errors"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

// The Firewall Rules API and Filters API will still work until 2025-06-15. There will be a single list of rules for both firewall rules and WAF custom rules, and this list contains WAF custom rules.
// https://developers.cloudflare.com/waf/reference/legacy/firewall-rules-upgrade/#new-api-and-terraform-resources
func tableCloudflareFirewallRule(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_firewall_rule",
		Description: "[DEPRECATED] Cloudflare Firewall Rules is a flexible and intuitive framework for filtering HTTP requests.",
		List: &plugin.ListConfig{
			Hydrate: listFirewallRules,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "Specifies the Firewall Rule identifier."},
			{Name: "zone_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ZoneID"), Description: "Specifies the zone identifier."},
			{Name: "paused", Type: proto.ColumnType_BOOL, Description: "Indicates whether the firewall rule is currently paused."},
			{Name: "priority", Type: proto.ColumnType_INT, Description: "The priority of the rule to allow control of processing order. A lower number indicates high priority. If not provided, any rules with a priority will be sequenced before those without."},
			{Name: "action", Type: proto.ColumnType_STRING, Description: "The action to apply to a matched request."},

			// Other columns
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the firewall rule is created."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "A description of the rule to help identify it."},
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the firewall rule is updated."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "Title of the resource."},

			// JSON columns
			{Name: "filter", Type: proto.ColumnType_JSON, Description: "A set of firewall properties."},
			{Name: "products", Type: proto.ColumnType_JSON, Description: "A list of products to bypass for a request when the bypass action is used."},
		}),
	}
}

//// LIST FUNCTION

func listFirewallRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	err := errors.New("The cloudflare_firewall_rule table has been deprecated and removed, please use cloudflare_ruleset table instead.")
	return nil, err
}
