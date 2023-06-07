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
		Name:             "cloudflare_user_audit_log",
		Description:      "Cloudflare User Audit Logs",
		DefaultTransform: transform.FromCamel(),
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:       "when",
					Operators:  []string{">", ">=", "<"},
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:    "actor_ip",
					Require: plugin.Optional,
				},
				{
					Name:    "actor_email",
					Require: plugin.Optional,
				},
				{
					Name:    "id",
					Require: plugin.Optional,
				},
			},
			Hydrate: listUserAuditLogs,
		},
		Columns: []*plugin.Column{
			{
				Name:        "actor_email",
				Description: "Email of the actor.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Actor.Email"),
			},
			{
				Name:        "actor_id",
				Description: "Unique identifier of the actor in Cloudflare’s system.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Actor.ID"),
			},
			{
				Name:        "actor_ip",
				Description: "Physical network address of the actor.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Actor.IP"),
			},
			{
				Name:        "actor_type",
				Description: "Type of user that started the audit trail.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Actor.Type"),
			},
			{
				Name:        "id",
				Description: "Unique identifier of an audit log.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "new_value",
				Description: "Contains the new value for the audited item.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "old_value",
				Description: "Contains the old value for the audited item.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "owner_id",
				Description: "The identifier of the user that was acting or was acted on behalf of. If a user did the action themselves, this value will be the same as the ActorID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Owner.ID"),
			},
			{
				Name:        "when",
				Description: "When the change happened.",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			{
				Name:        "action",
				Description: "The action that was taken.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "metadata",
				Description: "Additional audit log-specific information. Metadata is organized in key:value pairs. Key and Value formats can vary by ResourceType.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "new_value_json",
				Description: "Contains the new value for the audited item in JSON format.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("NewValueJSON"),
			},
			{
				Name:        "old_value_json",
				Description: "Contains the old value for the audited item in JSON format.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("OldValueJSON"),
			},
			{
				Name:        "resource",
				Description: "Target resource the action was performed on.",
				Type:        proto.ColumnType_JSON,
			},
		},
	}
}

func listUserAuditLogs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	quals := d.KeyColumnQuals

	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	opts := cloudflare.AuditLogFilter{}
	if quals["actor_ip"] != nil && quals["actor_ip"].GetStringValue() != "" {
		opts.ActorIP = quals["actor_ip"].GetStringValue()
	}
	if quals["actor_email"] != nil && quals["actor_email"].GetStringValue() != "" {
		opts.ActorEmail = quals["actor_email"].GetStringValue()
	}
	if quals["id"] != nil && quals["id"].GetStringValue() != "" {
		opts.ID = quals["id"].GetStringValue()
	}
	if d.Quals["when"] != nil {
		for _, q := range d.Quals["when"].Quals {
			timestamp := q.Value.GetStringValue()
			switch q.Operator {
			case ">=", ">":
				opts.Since = timestamp
			case "<":
				opts.Before = timestamp
			}
		}
	}

	items, err := conn.GetUserAuditLogs(ctx, opts)

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
