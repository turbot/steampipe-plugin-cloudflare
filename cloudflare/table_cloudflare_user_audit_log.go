package cloudflare

import (
	"context"
	"time"

	"github.com/cloudflare/cloudflare-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
					Operators:  []string{">", ">=", "<", "<=", "="},
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
		Columns: commonColumns([]*plugin.Column{
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
				Transform:   transform.From(convertAuditLogTimeToRFC3339Timestamp),
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
		}),
	}
}

func listUserAuditLogs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("cloudflare_user_audit_log.listUserAuditLogs", "connection_error", err)
		return nil, err
	}

	opts := cloudflare.AuditLogFilter{
		Page:    1,
		PerPage: 1000,
	}
	if d.EqualsQualString("actor_ip") != "" {
		opts.ActorIP = d.EqualsQualString("actor_ip")
	}
	if d.EqualsQualString("actor_email") != "" {
		opts.ActorEmail = d.EqualsQualString("actor_email")
	}
	if d.EqualsQualString("id") != "" {
		opts.ID = d.EqualsQualString("id")
	}
	if d.Quals["when"] != nil {
		for _, q := range d.Quals["when"].Quals {
			timestamp := q.Value.GetTimestampValue().AsTime().Format(time.RFC3339)
			timestampAdd := q.Value.GetTimestampValue().AsTime().Add(time.Second).Format(time.RFC3339)
			switch q.Operator {
			case ">=", ">":
				opts.Since = timestamp
			case "<":
				opts.Before = timestamp
			case "<=":
				opts.Before = timestampAdd
			case "=":
				opts.Since = timestamp
				opts.Before = timestampAdd
			}
		}
	}

	for {
		items, err := conn.GetUserAuditLogs(ctx, opts)
		if err != nil {
			plugin.Logger(ctx).Error("cloudflare_user_audit_log.listUserAuditLogs", "api_error", err)
			return nil, err
		}

		// Value of items.TotalPages is always 0, so we need to check if we have reached the end of the list
		if len(items.Result) == 0 {
			break
		}

		for _, i := range items.Result {
			d.StreamListItem(ctx, i)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		opts.Page = opts.Page + 1
	}

	return nil, nil
}

//// TRANSFORM FUNCTION

func convertAuditLogTimeToRFC3339Timestamp(_ context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(cloudflare.AuditLog)
	return data.When.Format(time.RFC3339), nil
}
