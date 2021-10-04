package cloudflare

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableCloudflareWorkersKVNamespace(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_workers_kv_namespace",
		Description: "A Namespace is a collection of key-value pairs stored in Workers KV.",
		List: &plugin.ListConfig{
			Hydrate: listWorkersKVNamespaces,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Namespace identifier."},
			{Name: "title", Type: proto.ColumnType_STRING, Description: "A human-readable string name for a Namespace."},
			// {Name: "keys", Type: proto.ColumnType_JSON, Transform: transform.FromField("Result"), Hydrate: listWorkersKVNamespaceKeys, Description: "A human-readable string name for a Namespace."},
		},
	}
}

func listWorkersKVNamespaces(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connect(ctx, d)
	if err != nil {
		logger.Error("listWorkersKVNamespaces", "connection error", err)
		return nil, err
	}
	// Paging is handled by default inside the API
	resp, err := conn.ListWorkersKVNamespaces(ctx)
	if err != nil {
		logger.Error("listWorkersKVNamespaces", "api error", err)
		return nil, err
	}
	for _, resource := range resp {
		d.StreamListItem(ctx, resource)
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS
// func listWorkersKVNamespaceKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	logger := plugin.Logger(ctx)
// 	namespace := h.Item.(cloudflare.WorkersKVNamespace).ID

// 	conn, err := connect(ctx, d)
// 	if err != nil {
// 		logger.Error("listWorkersKVNamespaceKeys", "connection error", err)
// 		return nil, err
// 	}
// 	TODO check for paging
// 	resp, err := conn.ListWorkersKVs(ctx, namespace)
// 	if err != nil {
// 		logger.Error("listWorkersKVNamespaceKeys", "api error", err)
// 		return nil, err
// 	}

// 	return resp, nil
// }