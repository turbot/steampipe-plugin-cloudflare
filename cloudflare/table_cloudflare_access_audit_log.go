package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

//// TABLE DEFINITION

func tableCloudflareAccessAuditLog(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_access_audit_log",
		Description: "Cloudflare Access Audit Log",
		List: &plugin.ListConfig{
			Hydrate:       listAccessAuditLogs,
			ParentHydrate: listAccount,
		},
		Columns: []*plugin.Column{
			{
				Name:        "user_email",
				Description: "The email of the user who generated the log record.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ip_address",
				Description: "The IP address from which the user accessed the service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "app_uid",
				Description: "A unique identifier for the app that was accessed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "app_domain",
				Description: "The domain associated with the app.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "action",
				Description: "The action that was performed by the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "connection",
				Description: "The type or state of the connection during the access.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "allowed",
				Description: "A boolean indicating if the action was allowed or not.",
				Type:        proto.ColumnType_BOOL,
			},
			{
				Name:        "created_at",
				Description: "The timestamp when the log record was created.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "ray_id",
				Description: "A unique identifier that Cloudflare uses for tracking individual requests.",
				Type:        proto.ColumnType_STRING,
			},
		},
	}
}

//// LIST FUNCTIONS

func listAccessAuditLogs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	account := h.Item.(cloudflare.Account)

	resp, err := conn.AccessAuditLogs(ctx, account.ID, cloudflare.AccessAuditLogFilterOptions{})
	if err != nil {
		return nil, err
	}

	for _, log := range resp {
		d.StreamListItem(ctx, log)
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS
