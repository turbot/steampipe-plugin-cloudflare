package cloudflare

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableCloudflareAPIToken(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_api_token",
		Description: "API tokens for the user.",
		List: &plugin.ListConfig{
			Hydrate: listAPIToken,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "ID of the API token."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the API token."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the API token."},

			// Other columns
			{Name: "condition", Type: proto.ColumnType_JSON, Description: "Conditions (e.g. client IP ranges) associated with the API token."},
			{Name: "expires_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the API token expires."},
			{Name: "issued_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the API token was issued."},
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the API token was last modified."},
			{Name: "not_before", Type: proto.ColumnType_TIMESTAMP, Description: "When the API token becomes valid."},

			// JSON columns
			{Name: "policies", Type: proto.ColumnType_JSON, Description: "Policies associated with this API token."},
		}),
	}
}

func listAPIToken(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	items, err := conn.APITokens(ctx)
	if err != nil {
		return nil, err
	}
	for _, i := range items {
		d.StreamListItem(ctx, i)
	}
	return nil, nil
}
