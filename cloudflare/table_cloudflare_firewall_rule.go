package cloudflare

import (
	"context"
	"time"

	"github.com/cloudflare/cloudflare-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type firewallRuleInfo = struct {
	ID          string            `json:"id,omitempty"`
	Paused      bool              `json:"paused"`
	Description string            `json:"description"`
	Action      string            `json:"action"`
	Priority    interface{}       `json:"priority"`
	Filter      cloudflare.Filter `json:"filter"`
	Products    []string          `json:"products,omitempty"`
	CreatedOn   time.Time         `json:"created_on,omitempty"`
	ModifiedOn  time.Time         `json:"modified_on,omitempty"`
	ZoneID      string
}

//// TABLE DEFINITION

func tableCloudflareFirewallRule(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_firewall_rule",
		Description: "Cloudflare Firewall Rules is a flexible and intuitive framework for filtering HTTP requests.",
		List: &plugin.ListConfig{
			Hydrate:       listFirewallRules,
			ParentHydrate: listZones,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"zone_id", "id"}),
			ShouldIgnoreError: isNotFoundError([]string{"HTTP status 404"}),
			Hydrate:           getFirewallRule,
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
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	zoneDetails := h.Item.(cloudflare.Zone)

	resp, err := conn.FirewallRules(ctx, zoneDetails.ID, cloudflare.PaginationOptions{})
	if err != nil {
		return nil, err
	}
	for _, i := range resp {
		d.StreamLeafListItem(ctx, firewallRuleInfo{
			ID:          i.ID,
			Paused:      i.Paused,
			Description: i.Description,
			Action:      i.Action,
			Priority:    i.Priority,
			Filter:      i.Filter,
			Products:    i.Products,
			CreatedOn:   i.CreatedOn,
			ModifiedOn:  i.ModifiedOn,
			ZoneID:      zoneDetails.ID,
		})
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getFirewallRule(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	zoneID := d.EqualsQuals["zone_id"].GetStringValue()
	id := d.EqualsQuals["id"].GetStringValue()

	op, err := conn.FirewallRule(ctx, zoneID, id)
	if err != nil {
		return nil, err
	}
	return firewallRuleInfo{
		ID:          op.ID,
		Paused:      op.Paused,
		Description: op.Description,
		Action:      op.Action,
		Priority:    op.Priority,
		Filter:      op.Filter,
		Products:    op.Products,
		CreatedOn:   op.CreatedOn,
		ModifiedOn:  op.ModifiedOn,
		ZoneID:      zoneID,
	}, nil
}
