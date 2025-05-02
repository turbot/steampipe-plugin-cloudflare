package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/accounts"
	"github.com/cloudflare/cloudflare-go/v4/audit_logs"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableCloudflareUserAuditLog() *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_user_audit_log",
		Description: "User audit logs provide a record of actions performed by users.",
		List: &plugin.ListConfig{
			ParentHydrate: listAccountForParent,
			Hydrate:       listUserAuditLogs,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the audit log entry."},
			{Name: "account_id", Type: proto.ColumnType_STRING, Description: "The account ID associated with the audit log entry."},
			{Name: "action", Type: proto.ColumnType_JSON, Description: "Details about the action performed."},
			{Name: "actor", Type: proto.ColumnType_JSON, Description: "Information about the user who performed the action."},
			{Name: "resource", Type: proto.ColumnType_JSON, Description: "Details about the resource that was affected."},
			{Name: "interface", Type: proto.ColumnType_STRING, Description: "The interface through which the action was performed."},
			{Name: "metadata", Type: proto.ColumnType_JSON, Description: "Additional metadata about the audit log entry."},
			{Name: "when", Type: proto.ColumnType_TIMESTAMP, Description: "When the action was performed."},
		},
	}
}

func listUserAuditLogs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accountData := h.Item.(accounts.Account)

	auditService := audit_logs.NewAuditLogService()
	result, err := auditService.List(ctx, audit_logs.AuditLogListParams{
		AccountID: cloudflare.F(accountData.ID),
	})
	if err != nil {
		plugin.Logger(ctx).Error("cloudflare_user_audit_log.listUserAuditLogs", "list_error", err)
		return nil, err
	}

	for _, log := range result.Result {
		d.StreamListItem(ctx, log)
	}

	return nil, nil
}
