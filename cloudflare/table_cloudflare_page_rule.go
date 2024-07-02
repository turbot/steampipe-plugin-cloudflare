package cloudflare

import (
	"context"
	"time"

	"github.com/cloudflare/cloudflare-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type pageRuleInfo = struct {
	ID         string                      `json:"id,omitempty"`
	Targets    []cloudflare.PageRuleTarget `json:"targets"`
	Actions    []cloudflare.PageRuleAction `json:"actions"`
	Priority   int                         `json:"priority"`
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
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	zoneDetails := h.Item.(cloudflare.Zone)

	resp, err := conn.ListPageRules(ctx, zoneDetails.ID)
	if err != nil {
		return nil, err
	}
	for _, i := range resp {
		d.StreamLeafListItem(ctx, pageRuleInfo{
			ID:         i.ID,
			Status:     i.Status,
			CreatedOn:  i.CreatedOn,
			ModifiedOn: i.ModifiedOn,
			Priority:   i.Priority,
			Actions:    i.Actions,
			Targets:    i.Targets,
			ZoneID:     zoneDetails.ID,
		})
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getPageRule(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	zoneID := d.EqualsQuals["zone_id"].GetStringValue()
	id := d.EqualsQuals["id"].GetStringValue()

	op, err := conn.PageRule(ctx, zoneID, id)
	if err != nil {
		return nil, err
	}
	return pageRuleInfo{
		ID:         op.ID,
		Status:     op.Status,
		CreatedOn:  op.CreatedOn,
		ModifiedOn: op.ModifiedOn,
		Priority:   op.Priority,
		Actions:    op.Actions,
		Targets:    op.Targets,
		ZoneID:     zoneID,
	}, nil
}
