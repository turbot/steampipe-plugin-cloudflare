package cloudflare

import (
	"context"
	"time"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/page_rules"
	"github.com/cloudflare/cloudflare-go/v4/zones"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type pageRuleInfo = struct {
	ID         string                      `json:"id,omitempty"`
	Targets    []page_rules.Target         `json:"targets"`
	Actions    []page_rules.PageRuleAction `json:"actions"`
	Priority   int64                       `json:"priority"`
	Status     string                      `json:"status"`
	ModifiedOn time.Time                   `json:"modified_on,omitempty"`
	CreatedOn  time.Time                   `json:"created_on,omitempty"`
	ZoneID     string
}

//// TABLE DEFINITION

func tableCloudflarePageRule(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_page_rule",
		Description: "Page Rules gives the ability to control how Cloudflare works on a URL or subdomain basis.",
		List: &plugin.ListConfig{
			Hydrate:       listPageRules,
			ParentHydrate: listZones,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"zone_id", "id"}),
			ShouldIgnoreError: isNotFoundError([]string{"HTTP status 404"}),
			Hydrate:           getPageRule,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "Specifies the Page Rule identifier."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Specifies the status of the page rule."},

			// Other columns
			{Name: "zone_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ZoneID"), Description: "Specifies the zone identifier."},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the page rule is created."},
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Description: "The time when the page rule was last modified."},
			{Name: "priority", Type: proto.ColumnType_INT, Description: "A number that indicates the preference for a page rule over another."},
			{Name: "title", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "Title of the resource."},

			// JSON columns
			{Name: "actions", Type: proto.ColumnType_JSON, Description: "A list of actions to perform if the targets of this rule match the request. Actions can redirect the url to another url or override settings (but not both)."},
			{Name: "targets", Type: proto.ColumnType_JSON, Description: "A list of targets to evaluate on a request."},
		}),
	}
}

//// LIST FUNCTION

func listPageRules(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_page_rule.listPageRules", "connection error", err)
		return nil, err
	}
	zoneDetails := h.Item.(zones.Zone)

	input := page_rules.PageRuleListParams{
		ZoneID: cloudflare.F(zoneDetails.ID),
	}

	resp, err := conn.PageRules.List(ctx, input)
	if err != nil {
		logger.Error("cloudflare_page_rule.listPageRules", "PageRules api error", err)
		return nil, err
	}

	for _, rule := range *resp {
		d.StreamLeafListItem(ctx, pageRuleInfo{
			ID:         rule.ID,
			Status:     string(rule.Status),
			CreatedOn:  rule.CreatedOn,
			ModifiedOn: rule.ModifiedOn,
			Priority:   rule.Priority,
			Actions:    rule.Actions,
			Targets:    rule.Targets,
			ZoneID:     zoneDetails.ID,
		})
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getPageRule(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_page_rule.getPageRule", "connection error", err)
		return nil, err
	}

	zoneID := d.EqualsQuals["zone_id"].GetStringValue()
	id := d.EqualsQuals["id"].GetStringValue()

	input := page_rules.PageRuleGetParams{
		ZoneID: cloudflare.F(zoneID),
	}

	rule, err := conn.PageRules.Get(ctx, id, input)
	if err != nil {
		logger.Error("cloudflare_page_rule.getPageRule", "PageRule api error", err)
		return nil, err
	}
	return pageRuleInfo{
		ID:         rule.ID,
		Status:     string(rule.Status),
		CreatedOn:  rule.CreatedOn,
		ModifiedOn: rule.ModifiedOn,
		Priority:   rule.Priority,
		Actions:    rule.Actions,
		Targets:    rule.Targets,
		ZoneID:     zoneID,
	}, nil
}
