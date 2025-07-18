package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/accounts"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "ID of the account."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the account."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Type of the account.", Transform: transform.FromP(getExtraFieldsFromAPIresponse, "type")},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "The create time when account was created."},

			// JSON columns
			{Name: "settings", Type: proto.ColumnType_JSON, Description: "Settings for the account."},
		}),
	}
}

func listAccount(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_account.listAccount", "connection error", err)
		return nil, err
	}
	maxLimit := int32(500)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := accounts.AccountListParams{
		PerPage: cloudflare.F(float64(maxLimit)),
	}

	iter := conn.Accounts.ListAutoPaging(ctx, input)
	for iter.Next() {
		account := iter.Current()
		d.StreamListItem(ctx, account)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	if err := iter.Err(); err != nil {
		logger.Error("cloudflare_account.listAccount", "Accounts api error", err)
		return nil, err
	}

	return nil, nil
}

func getAccount(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_account.getAccount", "connection error", err)
		return nil, err
	}
	quals := d.EqualsQuals
	id := quals["id"].GetStringValue()

	input := accounts.AccountGetParams{
		AccountID: cloudflare.F(id),
	}
	account, err := conn.Accounts.Get(ctx, input)
	if err != nil {
		logger.Error("cloudflare_account.getAccount", "Account api error", err)
		return nil, err
	}
	return account, nil
}

//// TRANSFORM FUNCTIONS

func getExtraFieldsFromAPIresponse(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	response := d.HydrateItem.(accounts.Account)
	param := d.Param.(string)

	extraFields, err := toMap(response.JSON.RawJSON())
	if err != nil {
		logger.Error("cloudflare_account.getExtraFieldsFromAPIresponse", "JSON parsing error", err)
		return nil, err
	}

	return extraFields[param], nil
}
