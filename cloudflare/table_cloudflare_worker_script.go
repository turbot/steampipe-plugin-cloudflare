package cloudflare

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableCloudflareWorkerScript(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_worker_script",
		Description: "A Worker script is a single script that is executed on matching routes in the Cloudflare edge.",
		List: &plugin.ListConfig{
			Hydrate: listWorkerScripts,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Script identifier."},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the script was created."},
			{Name: "etag", Type: proto.ColumnType_STRING, Description: "Hashed script content, can be used in a If-None-Match header when updating."},
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the script was last modified."},
			// Always coming as 0
			// {Name: "size", Type: proto.ColumnType_INT, Description: "Size of the script, in bytes."},
		},
	}
}

func listWorkerScripts(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connect(ctx, d)
	if err != nil {
		logger.Error("listWorkersScripts", "connection error", err)
		return nil, err
	}
	resp, err := conn.ListWorkerScripts(ctx)
	if err != nil {
		logger.Error("listWorkersScripts", "api error", err)
		return nil, err
	}
	for _, resource := range resp.WorkerList {
		d.StreamListItem(ctx, resource)
	}
	return nil, nil
}
