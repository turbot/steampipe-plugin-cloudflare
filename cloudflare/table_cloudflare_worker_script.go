package cloudflare

import (
	"context"
	"encoding/json"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/workers"
	"github.com/cloudflare/cloudflare-go/v4/accounts"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableCloudflareWorkerScript(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_worker_script",
		Description: "Cloudflare Worker Scripts host custom serverless code executed at Cloudflare’s edge for enhanced logic, routing, and performance optimizations.",
		List: &plugin.ListConfig{
			Hydrate:       listWorkerScripts,
			ParentHydrate: listAccount, 
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "account_id", Require: plugin.Optional},
			},
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "The id of the script in the Workers system. Usually the script name."},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the script was created."},
			{Name: "etag", Type: proto.ColumnType_STRING, Description: "Hashed script content, can be used in a If-None-Match header when updating."},
			{Name: "has_assets", Type: proto.ColumnType_BOOL, Description: "Whether a Worker contains assets."},
			{Name: "has_modules", Type: proto.ColumnType_BOOL, Description: "Whether a Worker contains modules."},
			{Name: "logpush", Type: proto.ColumnType_BOOL, Description: "Whether Logpush is turned on for the Worker."},
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the script was last modified."},
			{Name: "usage_model", Type: proto.ColumnType_STRING, Description: "Usage model for the Worker invocations."},
			{Name: "account_name", Type: proto.ColumnType_STRING,  Hydrate: getParentAccountDetails, Transform: transform.FromField("Name"), Description: "Account name."},
		
			// Query columns for filtering
			{Name: "account_id", Type: proto.ColumnType_STRING, Hydrate: getParentAccountDetails, Transform: transform.FromField("ID"), Description: "The account ID to filter Worker scripts."},
		
			// JSON Columns
			{Name: "subdomain", Type: proto.ColumnType_JSON, Hydrate: getWorkerSubdomain, Transform: transform.FromValue(),Description: "Whether the Worker is available on the workers.dev subdomain."},
			{Name: "tail_consumers", Type: proto.ColumnType_JSON, Description: "List of Workers that will consume logs from the attached Worker."},
			{Name: "placement", Type: proto.ColumnType_JSON, Transform: transform.FromField("Placement"), Description: "Configuration for Smart Placement."},
		}),
	}
}

//// LIST FUNCTION

// listWorkerScripts retrieves all Worker scripts across all accounts.
func listWorkerScripts(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	accountDetails := h.Item.(accounts.Account)

	inputAccountId := d.EqualsQualString("account_id")

	// Only list scripts for accounts stated in the input query
	if inputAccountId != "" && inputAccountId != accountDetails.ID {
		return nil, nil
	}

	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_worker_script.listWorkerScripts", "connect error", err)
		return nil, err
	}

	input := workers.ScriptListParams{
		AccountID: cloudflare.F(accountDetails.ID),
	}

	iter := conn.Workers.Scripts.ListAutoPaging(ctx, input)
	if err := iter.Err(); err != nil {
		logger.Error("cloudflare_worker_script.listWorkerScripts", "api call error", err)
		return nil, err
	}

	for iter.Next() {
		script := iter.Current()
		d.StreamListItem(ctx, script)
	}
	return nil, nil
}

//// GET FUNCTION

// getWorkerSubdomain returns if the Worker is available on the workers.dev subdomain.
//
// Parameters:
// - id: The Worker script identifier (required)
func getWorkerSubdomain(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	account := h.ParentItem.(accounts.Account)
    script := h.Item.(workers.Script)

	conn, err := connectV4(ctx, d)
	if err != nil {
		return nil, err
	}
	input := workers.ScriptSubdomainGetParams{
		AccountID: cloudflare.F(account.ID),
	}
	subdomain, err := conn.Workers.Scripts.Subdomain.Get(ctx,script.ID,input)
	if err != nil {
		return nil, err
	}

	// SDK does not map the response correctly, therefore returning the raw json instead
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(subdomain.JSON.RawJSON()), &m); err != nil {
		logger.Error("cloudflare_worker_script.getWorkerSubdomain", err)
		return nil, err
	}
	return m["result"], nil
}

func getParentAccountDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return h.ParentItem.(accounts.Account), nil
}