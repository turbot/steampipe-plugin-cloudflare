package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/accounts"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableCloudflareAPIToken(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_api_token",
		Description: "API tokens for the user.",
		List: &plugin.ListConfig{
			ParentHydrate: listAccount,
			Hydrate:       listAPIToken,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "ID of the API token."},
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

func listAPIToken(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_api_token.listAPIToken", "connection error", err)
		return nil, err
	}

	account := h.Item.(accounts.Account)

	iter := conn.Accounts.Tokens.ListAutoPaging(ctx, accounts.TokenListParams{
		AccountID: cloudflare.String(account.ID),
	})
	if err := iter.Err(); err != nil {
		logger.Error("cloudflare_api_token.listAPIToken", "APITokens api error", err)
		return nil, err
	}

	for iter.Next() {
		token := iter.Current()
		d.StreamListItem(ctx, token)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
