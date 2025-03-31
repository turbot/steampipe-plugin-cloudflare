package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4/accounts"
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

	conn, err := connectV4(ctx, d)
	if err != nil {
		return nil
	}

	page, err := conn.Accounts.List(context.TODO(), accounts.AccountListParams{})
	if err != nil {
		panic(err.Error())
	}
	matrix := make([]map[string]interface{}, len(page.Result))
	for page != nil {
		for i, account := range page.Result {
			matrix[i] = map[string]interface{}{matrixKeyAccount: account.ID}
		}
		page, err = page.GetNextPage()
	}

	// set cache
	d.ConnectionManager.Cache.Set(cacheKey, matrix)
	return matrix
}
