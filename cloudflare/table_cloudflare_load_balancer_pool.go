package cloudflare

import (
	"context"
	"strings"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/accounts"
	"github.com/cloudflare/cloudflare-go/v4/load_balancers"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableCloudflareLoadBalancerPool(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_load_balancer_pool",
		Description: "A pool is a group of origin servers, with each origin identified by its IP address or hostname.",
		List: &plugin.ListConfig{
			Hydrate:       listLoadBalancerPools,
			ParentHydrate: listAccount,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "The API item identifier."},
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

			{Name: "account_name", Type: proto.ColumnType_STRING, Hydrate: getLoadBalancerPoolAccountName, Transform: transform.FromValue(), Description: "The name of the account associated with the load balancer pool."},

			// JSON columns
			{Name: "check_regions", Type: proto.ColumnType_JSON, Description: "A list of regions (specified by region code) from which to run health checks."},
			{Name: "load_shedding", Type: proto.ColumnType_JSON, Description: "Setting for controlling load shedding for this pool."},
			{Name: "origins", Type: proto.ColumnType_JSON, Description: "The list of origins within this pool. Traffic directed at this pool is balanced across all currently healthy origins, provided the pool itself is healthy."},
			{Name: "health", Type: proto.ColumnType_JSON, Hydrate: getLoadBalancerPoolHealth, Transform: transform.FromValue(), Description: "The Pool Health details."},
		}),
	}
}

func listLoadBalancerPools(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	account := h.Item.(accounts.Account)

	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_load_balancer_pool.listLoadBalancerPools", "connection_error", err)
		return nil, err
	}
	// Rest api only supports monitor as an input.
	input := load_balancers.PoolListParams{
		AccountID: cloudflare.F(account.ID),
	}

	iter := conn.LoadBalancers.Pools.ListAutoPaging(ctx, input)
	if err := iter.Err(); err != nil {
		logger.Error("cloudflare_load_balancer_pool.listLoadBalancerPools", "api error", err)
		return nil, err
	}
	for iter.Next() {
		resource := iter.Current()
		d.StreamListItem(ctx, resource)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, nil
}

func getLoadBalancerPoolHealth(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		return nil, err
	}
	pool := h.Item.(load_balancers.Pool)
	account := h.ParentItem.(accounts.Account)

	input := load_balancers.PoolHealthGetParams{
		AccountID: cloudflare.F(account.ID),
	}

	pool_health, err := conn.LoadBalancers.Pools.Health.Get(ctx, pool.ID, input)
	if err != nil {
		// This setting might not be available for all zones
		if strings.Contains(err.Error(), "Health info unavailable") {
			return nil, nil
		}
		logger.Error("cloudflare_load_balancer_pool.getLoadBalancerPoolHealth", "load balancer pool health API error", err)
		return nil, nil
	}
	return pool_health, nil
}

func getLoadBalancerPoolAccountName(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
  account := h.ParentItem.(accounts.Account)
  return account.Name, nil
}