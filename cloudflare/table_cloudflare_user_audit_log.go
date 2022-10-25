package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableCloudflareUserAuditLog(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_user_audit_log",
		Description: "Audit logs",
		List: &plugin.ListConfig{
			Hydrate: listUserAuditLogs,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "action", Type: proto.ColumnType_JSON, Transform: transform.FromField("Action"), Description: "ID of the account."},
			{Name: "actor_email", Type: proto.ColumnType_STRING, Transform: transform.FromField("Actor.Email"), Description: "Name of the account."},
			{Name: "actor_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Actor.ID"), Description: "Name of the account."},
			{Name: "actor_ip", Type: proto.ColumnType_STRING, Transform: transform.FromField("Actor.IP"), Description: "Name of the account."},
			{Name: "actor_type", Type: proto.ColumnType_STRING, Transform: transform.FromField("Actor.Type"), Description: "Name of the account."},
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "Type of the account."},
			{Name: "new_value", Type: proto.ColumnType_STRING, Transform: transform.FromField("NewValue"), Description: "Type of the account."},
			{Name: "old_value", Type: proto.ColumnType_STRING, Transform: transform.FromField("OldValue"), Description: "Type of the account."},
			{Name: "owner_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Owner.ID"), Description: "Type of the account."},
			{Name: "when", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("When"), Description: "Type of the account."},
		},
	}
}

func listUserAuditLogs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	items, err := conn.GetUserAuditLogs(ctx, cloudflare.AuditLogFilter{})
	if err != nil {
		return nil, err
	}

	for _, i := range items.Result {
		d.StreamListItem(ctx, i)

		// Context can be cancelled due to manual cancellation or thelimit has been hit
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}