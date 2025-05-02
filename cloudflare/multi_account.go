package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/accounts"
	"github.com/cloudflare/cloudflare-go/v4/option"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

const matrixKeyAccount = "account_id"

// BuildAccountmatrix :: return a list of matrix items, one per account.
// Allows to perform three level resource listing as in case of cloudflare_access_policy
// (i.e List Account -> List Applications -> List Access policies for each application)
func BuildAccountmatrix(ctx context.Context, d *plugin.QueryData) []map[string]interface{} {

	// cache matrix
	cacheKey := "AccountListMatrix"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.([]map[string]interface{})
	}

	conn, err := connect(ctx, d)
	if err != nil {
		return nil
	}

	// Get the client's options from the connection
	opts := []option.RequestOption{}
	if conn != nil {
		opts = append(opts, conn.Options...)
	}

	accountService := accounts.NewAccountService(opts...)
	iter := accountService.ListAutoPaging(ctx, accounts.AccountListParams{})

	var items []accounts.Account
	for iter.Next() {
		items = append(items, iter.Current())
	}
	if err := iter.Err(); err != nil {
		return nil
	}

	matrix := make([]map[string]interface{}, len(items))
	for i, account := range items {
		matrix[i] = map[string]interface{}{matrixKeyAccount: account.ID}
	}

	// set cache
	d.ConnectionManager.Cache.Set(cacheKey, matrix)
	return matrix
}

func listAccountForParent(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
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

func getAccountForParent(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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
	account, err := accountService.Get(ctx, accounts.AccountGetParams{
		AccountID: cloudflare.F(h.Item.(string)),
	})
	if err != nil {
		return nil, err
	}

	return account, nil
}
