package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableCloudflareAccount(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_account",
		Description: "Accounts the user has access to.",
		List: &plugin.ListConfig{
			Hydrate: listAccount,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			Hydrate:           getAccount,
			ShouldIgnoreError: isNotFoundError([]string{"HTTP status 404"}),
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "ID of the account."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the account."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Type of the account."},

			// JSON columns
			{Name: "settings", Type: proto.ColumnType_JSON, Description: "Settings for the account."},
		}),
	}
}

func listAccount(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	items, _, err := conn.Accounts(ctx, cloudflare.PaginationOptions{})
	if err != nil {
		return nil, err
	}
	for _, i := range items {
		d.StreamListItem(ctx, i)
	}
	return nil, nil
}

func getAccount(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	quals := d.EqualsQuals
	id := quals["id"].GetStringValue()
	account, _, err := conn.Account(ctx, id)
	if err != nil {
		return nil, err
	}
	return account, nil
}
