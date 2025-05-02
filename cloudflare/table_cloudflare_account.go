package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/accounts"
	"github.com/cloudflare/cloudflare-go/v4/option"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableCloudflareAccount(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_account",
		Description: "Accounts are separate organizational units within Cloudflare.",
		List: &plugin.ListConfig{
			Hydrate: listAccount,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"Invalid account identifier"}),
			Hydrate:           getAccount,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Account identifier tag."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Account name."},

			// Other columns
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the account was created."},
			{Name: "enforce_twofactor", Type: proto.ColumnType_BOOL, Description: "Whether 2FA is enforced."},
			{Name: "settings", Type: proto.ColumnType_JSON, Description: "Settings for this account."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Account type."},
		}),
	}
}

func listAccount(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get the client's options from the connection
	opts := []option.RequestOption{}
	if conn != nil {
		opts = append(opts, conn.Options...)
	}

	accountService := accounts.NewAccountService(opts...)
	iter := accountService.ListAutoPaging(ctx, accounts.AccountListParams{})
	for iter.Next() {
		account := iter.Current()
		d.StreamListItem(ctx, account)
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func getAccount(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get the client's options from the connection
	opts := []option.RequestOption{}
	if conn != nil {
		opts = append(opts, conn.Options...)
	}

	quals := d.EqualsQuals
	accountID := quals["id"].GetStringValue()

	accountService := accounts.NewAccountService(opts...)
	account, err := accountService.Get(ctx, accounts.AccountGetParams{
		AccountID: cloudflare.F(accountID),
	})
	if err != nil {
		return nil, err
	}

	return account, nil
}
