package cloudflare

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableCloudflareLoadBalancerPool(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_load_balancer_pool",
		Description: "A pool is a group of origin servers, with each origin identified by its IP address or hostname.",
		List: &plugin.ListConfig{
			Hydrate:       listLoadBalancerPools,
			ParentHydrate: listZones,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The API item identifier."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "A short name for the pool."},
			{Name: "enabled", Type: proto.ColumnType_BOOL, Description: "Status of this pool. Disabled pools will not receive traffic and are excluded from health checks."},
			{Name: "monitor", Type: proto.ColumnType_STRING, Description: "The ID of the Monitor to use for health checking origins within this pool."},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the pool was created."},

			// Other columns
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description for the pool."},
			{Name: "latitude", Type: proto.ColumnType_STRING, Description: "The latitude this pool is physically located at; used for proximity steering. Values should be between -90 and 90."},
			{Name: "longitude", Type: proto.ColumnType_STRING, Description: "The longitude this pool is physically located at; used for proximity steering. Values should be between -180 and 180."},
			{Name: "minimum_origins", Type: proto.ColumnType_INT, Description: "The minimum number of origins that must be healthy for this pool to serve traffic. If the number of healthy origins falls below this number, the pool will be marked unhealthy and we will failover to the next available pool. Default: 1."},
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the pool was last modified."},
			{Name: "notification_email", Type: proto.ColumnType_STRING, Description: "The email address to send health status notifications to. This can be an individual mailbox or a mailing list. Multiple emails can be supplied as a comma delimited list."},

			// JSON columns
			{Name: "check_regions", Type: proto.ColumnType_JSON, Description: "A list of regions (specified by region code) from which to run health checks."},
			{Name: "load_shedding", Type: proto.ColumnType_JSON, Description: "Setting for controlling load shedding for this pool."},
			{Name: "origins", Type: proto.ColumnType_JSON, Description: "The list of origins within this pool. Traffic directed at this pool is balanced across all currently healthy origins, provided the pool itself is healthy."},
		}),
	}
}

func listLoadBalancerPools(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	conn, err := connect(ctx, d)
	if err != nil {
		logger.Error("listLoadBalancers", "connection_error", err)
		return nil, err
	}
	// Rest api only supports monitor as an input.
	loadBalancersPools, err := conn.ListLoadBalancerPools(ctx)
	if err != nil {
		logger.Error("ListLoadBalancers", "api error", err)
		return nil, err
	}
	for _, resource := range loadBalancersPools {
		d.StreamListItem(ctx, resource)
	}
	return nil, nil
}
