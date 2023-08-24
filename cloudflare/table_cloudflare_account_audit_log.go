package cloudflare

import (
	"context"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableCloudflareAccountAuditLog(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_account_audit_log",
		Description: "Cloudflare Account Audit Logs",
		List: &plugin.ListConfig{
			Hydrate:       listAccountAuditLogs,
			ParentHydrate: listAccount,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:       "log_event_time",
					Operators:  []string{">", ">=", "<", "<=", "="},
					Require:    plugin.Optional,
					CacheMatch: "exact",
				},
				{
					Name:    "log_ip_address",
					Require: plugin.Optional,
				},
				{
					Name:    "log_actor_email",
					Require: plugin.Optional,
				},
				{
					Name:    "log_event_id",
					Require: plugin.Optional,
				},
			},
		},
		Columns: []*plugin.Column{
			{
				Name:        "log_actor_email",
				Description: "Email of the actor.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Actor.Email"),
			},
			{
				Name:        "log_actor_username",
				Description: "Username of the actor.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Actor.Email"),
			},
			{
				Name:        "log_actor_id",
				Description: "Unique identifier of the actor in Cloudflare’s system.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Actor.ID"),
			},
			{
				Name:        "log_ip_address",
				Description: "Physical network address of the actor.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Actor.IP"),
			},
			{
				Name:        "actor_type",
				Description: "Type of actor that started the audit trail.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Actor.Type"),
			},
			{
				Name:        "log_event_id",
				Description: "Unique identifier of an audit log.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "owner_id",
				Description: "The identifier of the actor that was acting or was acted on behalf of. If a actor did the action themselves, this value will be the same as the ActorID.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Owner.ID"),
			},
			{
				Name:        "log_event_time",
				Description: "When the change happened.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("When"),
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

//// LIST FUNCTIONS

func listAccountAuditLogs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("cloudflare_account_audit_log.listAccountAuditLogs", "conn_error", err)
		return nil, err
	}
	account := h.Item.(cloudflare.Account)

	opts := cloudflare.AuditLogFilter{
		Page:    1,
		PerPage: 1000,
	}
	if d.EqualsQualString("log_ip_address") != "" {
		opts.ActorIP = d.EqualsQualString("log_ip_address")
	}
	if d.EqualsQualString("log_actor_email") != "" {
		opts.ActorEmail = d.EqualsQualString("log_actor_email")
	}
	if d.EqualsQualString("log_event_id") != "" {
		opts.ID = d.EqualsQualString("log_event_id")
	}
	if d.Quals["log_event_time"] != nil {
		for _, q := range d.Quals["log_event_time"].Quals {
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
		items, err := conn.GetOrganizationAuditLogs(ctx, account.ID, opts)
		if err != nil {
			plugin.Logger(ctx).Error("cloudflare_account_audit_log.GetOrganizationAuditLogs", "api_error", err)
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
