package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4/accounts"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
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

	conn, err := connectV4(ctx, d)
	if err != nil {
		return nil
	}

	page, err := conn.Accounts.List(ctx, accounts.AccountListParams{})
	if err != nil {
		panic(err.Error())
	}
	var matrix []map[string]interface{}
	// GetNextPage always issues another API call and only returns a nil page on
	// error; it never signals "no more pages" by itself. Cloudflare returns an
	// empty result array once the page range is exhausted, so that emptiness is
	// what we must use to stop, otherwise this loops forever.
	for page != nil && len(page.Result) > 0 {
		for _, account := range page.Result {
			matrix = append(matrix, map[string]interface{}{matrixKeyAccount: account.ID})
		}
		if page, err = page.GetNextPage(); err != nil {
			panic(err.Error())
		}
	}

	// set cache
	d.ConnectionManager.Cache.Set(cacheKey, matrix)
	return matrix
}
