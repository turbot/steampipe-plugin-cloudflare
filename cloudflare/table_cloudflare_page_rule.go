package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/page_rules"
	"github.com/cloudflare/cloudflare-go/v4/zones"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableCloudflarePageRule() *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_page_rule",
		Description: "Page rules allow you to control how Cloudflare works on a URL or subdomain basis.",
		List: &plugin.ListConfig{
			ParentHydrate: listZones,
			Hydrate:       listPageRules,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"zone_id", "id"}),
			Hydrate:    getPageRule,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the page rule."},
			{Name: "zone_id", Type: proto.ColumnType_STRING, Description: "The zone ID where the page rule is configured."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the page rule (active/disabled)."},
			{Name: "priority", Type: proto.ColumnType_INT, Description: "The priority of the page rule."},
			{Name: "targets", Type: proto.ColumnType_JSON, Description: "The URL patterns that trigger the page rule."},
			{Name: "actions", Type: proto.ColumnType_JSON, Description: "The actions to perform when the page rule is triggered."},
		},
	}
}

func listPageRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	zoneData := h.Item.(zones.Zone)

	prService := page_rules.NewPageRuleService()
	result, err := prService.List(ctx, page_rules.PageRuleListParams{
		ZoneID: cloudflare.F(zoneData.ID),
	})
	if err != nil {
		plugin.Logger(ctx).Error("cloudflare_page_rule.listPageRules", "list_error", err)
		return nil, err
	}

	for _, rule := range *result {
		d.StreamListItem(ctx, rule)
	}

	return nil, nil
}

func getPageRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	zoneID := d.EqualsQuals["zone_id"].GetStringValue()
	id := d.EqualsQuals["id"].GetStringValue()

	prService := page_rules.NewPageRuleService()
	rule, err := prService.Get(ctx, id, page_rules.PageRuleGetParams{
		ZoneID: cloudflare.F(zoneID),
	})
	if err != nil {
		plugin.Logger(ctx).Error("cloudflare_page_rule.getPageRule", "get_error", err)
		return nil, err
	}

	return rule, nil
}
