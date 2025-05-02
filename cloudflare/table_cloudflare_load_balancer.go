package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/load_balancers"
	"github.com/cloudflare/cloudflare-go/v4/zones"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableCloudflareLoadBalancer() *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_load_balancer",
		Description: "Load balancers distribute traffic across multiple servers to optimize performance and reliability.",
		List: &plugin.ListConfig{
			ParentHydrate: listZones,
			Hydrate:       listLoadBalancers,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"zone_id", "id"}),
			Hydrate:    getLoadBalancer,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the load balancer."},
			{Name: "zone_id", Type: proto.ColumnType_STRING, Description: "The zone ID where the load balancer is configured."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The DNS hostname to associate with your load balancer."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "A description of the load balancer."},
			{Name: "enabled", Type: proto.ColumnType_BOOL, Description: "Whether the load balancer is enabled."},
			{Name: "ttl", Type: proto.ColumnType_INT, Description: "The time to live (TTL) of the DNS entry for the IP address returned by this load balancer."},
			{Name: "fallback_pool", Type: proto.ColumnType_STRING, Description: "The pool ID to use when all other pools are detected as unhealthy."},
			{Name: "default_pools", Type: proto.ColumnType_JSON, Description: "A list of pool IDs ordered by their failover priority."},
			{Name: "proxied", Type: proto.ColumnType_BOOL, Description: "Whether the hostname gets Cloudflare's origin protection."},
			{Name: "steering_policy", Type: proto.ColumnType_STRING, Description: "The steering policy for the load balancer."},
			{Name: "session_affinity", Type: proto.ColumnType_STRING, Description: "The method the load balancer uses to determine the route for new sessions."},
			{Name: "session_affinity_ttl", Type: proto.ColumnType_INT, Description: "Time, in seconds, until this load balancer's session affinity cookie expires after being created."},
			{Name: "session_affinity_attributes", Type: proto.ColumnType_JSON, Description: "Configures attributes for session affinity."},
			{Name: "rules", Type: proto.ColumnType_JSON, Description: "A list of rules for this load balancer."},
			{Name: "random_steering", Type: proto.ColumnType_JSON, Description: "Configures pool weights when using the random steering policy."},
			{Name: "adaptive_routing", Type: proto.ColumnType_JSON, Description: "Configures adaptive routing."},
			{Name: "location_strategy", Type: proto.ColumnType_JSON, Description: "Configures location strategy."},
		},
	}
}

func listLoadBalancers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	zoneData := h.Item.(zones.Zone)

	lbService := load_balancers.NewLoadBalancerService()
	result, err := lbService.List(ctx, load_balancers.LoadBalancerListParams{
		ZoneID: cloudflare.F(zoneData.ID),
	})
	if err != nil {
		plugin.Logger(ctx).Error("cloudflare_load_balancer.listLoadBalancers", "list_error", err)
		return nil, err
	}

	for _, lb := range result.Result {
		d.StreamListItem(ctx, lb)
	}

	return nil, nil
}

func getLoadBalancer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	zoneID := d.EqualsQuals["zone_id"].GetStringValue()
	id := d.EqualsQuals["id"].GetStringValue()

	lbService := load_balancers.NewLoadBalancerService()
	lb, err := lbService.Get(ctx, id, load_balancers.LoadBalancerGetParams{
		ZoneID: cloudflare.F(zoneID),
	})
	if err != nil {
		plugin.Logger(ctx).Error("cloudflare_load_balancer.getLoadBalancer", "get_error", err)
		return nil, err
	}

	return lb, nil
}
