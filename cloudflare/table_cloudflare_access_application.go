package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableCloudflareAccessApplication(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_access_application",
		Description: "Access Applications are used to restrict access to a whole application using an authorisation gateway managed by Cloudflare.",
		List: &plugin.ListConfig{
			Hydrate: listAccessApplications,
		},
		// Get Config - Currently SDK is not supporting get call
		Columns: []*plugin.Column{
			// Top fields
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Application API uuid."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Friendly name of the Access Application."},
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "The domain and path that Access will block."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The application type. Defaults to self_hosted. Valid values are self_hosted, ssh, vnc, or file."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the application was created."},

			// Other fields
			{Name: "aud", Type: proto.ColumnType_STRING, Description: "Audience tag."},
			{Name: "auto_redirect_to_identity", Type: proto.ColumnType_BOOL, Description: " Option to skip identity provider selection if only one is configured in allowed_idps. Defaults to false (disabled)."},
			{Name: "custom_deny_message", Type: proto.ColumnType_STRING, Description: "Option that returns a custom error message when a user is denied access to the application."},
			{Name: "custom_deny_url", Type: proto.ColumnType_STRING, Description: "Option that redirects to a custom URL when a user is denied access to the application."},
			{Name: "enable_binding_cookie", Type: proto.ColumnType_BOOL, Description: "Option to provide increased security against compromised authorization tokens and CSRF attacks by requiring an additional \"binding\" cookie on requests. Defaults to false."},
			{Name: "session_duration", Type: proto.ColumnType_STRING, Description: "How often a user will be forced to re-authorise. Must be in the format \"48h\" or \"2h45m\". Valid time units are ns, us (or µs), ms, s, m, h. Defaults to 24h."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the application was last modified."},

			// JSON fields
			{Name: "allowed_idps", Type: proto.ColumnType_JSON, Description: "The identity providers selected for the application."},
			{Name: "cors_headers", Type: proto.ColumnType_JSON, Description: "CORS configuration for the Access Application. See below for reference structure."},
		},
	}
}

func listAccessApplications(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	conn, err := connect(ctx, d)
	if err != nil {
		logger.Error("listAccessApplications", "connection error", err)
		return nil, err
	}

	opts := cloudflare.PaginationOptions{
		PerPage: 100,
		Page:    1,
	}

	for {
		items, result_info, err := conn.AccessApplications(ctx, conn.AccountID, opts)
		if err != nil {
			logger.Error("listAccessApplications", "AccessApplications api error", err)
			return nil, err
		}
		for _, i := range items {
			d.StreamListItem(ctx, i)
		}

		if result_info.Page >= result_info.TotalPages {
			break
		}
		opts.Page = opts.Page + 1
	}

	return nil, nil
}
