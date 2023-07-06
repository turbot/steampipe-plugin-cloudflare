package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
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

	items, _, err := conn.Accounts(ctx, cloudflare.PaginationOptions{})
	if err != nil {
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
