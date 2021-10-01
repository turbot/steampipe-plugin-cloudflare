package cloudflare

import (
	"context"
	"time"

	"github.com/cloudflare/cloudflare-go"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
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
		Description: "Cloudflare Page Rule",
		List: &plugin.ListConfig{
			Hydrate:       listPageRules,
			ParentHydrate: listZones,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"zone_id", "id"}),
			ShouldIgnoreError: isNotFoundError([]string{"HTTP status 404"}),
			Hydrate:           getPageRule,
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "Specifies the Page Rule identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "status",
				Description: "Specifies the status of the page rule.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "zone_id",
				Description: "Specifies the zone identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ZoneID"),
			},
			{
				Name:        "created_on",
				Description: "The time when the page rule is created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "modified_on",
				Description: "The time when the page rule was last modified.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "priority",
				Description: "A number that indicates the preference for a page rule over another.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "actions",
				Description: "A list of actions to perform if the targets of this rule match the request. Actions can redirect the url to another url or override settings (but not both).",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "targets",
				Description: "A list of targets to evaluate on a request.",
				Type:        proto.ColumnType_JSON,
			},

			// steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
		},
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

	zoneID := d.KeyColumnQuals["zone_id"].GetStringValue()
	id := d.KeyColumnQuals["id"].GetStringValue()

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
