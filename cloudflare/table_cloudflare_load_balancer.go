package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableCloudflareLoadBalancer(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_load_balancer",
		Description: "Cloudflare Load balancers allows to distribute traffic across servers, which reduces server strain and latency and improves the experience for end users.",
		List: &plugin.ListConfig{
			Hydrate:       listLoadBalancers,
			ParentHydrate: listZones,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "API item identifier."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The DNS hostname to associate with your Load Balancer."},
			{Name: "zone_name", Type: proto.ColumnType_STRING, Hydrate: getParentZoneDetails, Transform: transform.FromField("Name"), Description: "The zone name to which load balancer belongs."},
			{Name: "zone_id", Type: proto.ColumnType_STRING, Hydrate: getParentZoneDetails, Transform: transform.FromField("ID"), Description: "The zone ID to which load balancer belongs."},
			{Name: "ttl", Type: proto.ColumnType_INT, Description: "Time to live (TTL) of the DNS entry for the IP address returned by this load balancer. This only applies to gray-clouded (unproxied) load balancers."},
			{Name: "enabled", Type: proto.ColumnType_BOOL, Description: "Whether to enable (the default) this load balancer."},

			// Other columns
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "Load balancer creation time."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Load balancer description."},
			{Name: "fallback_pool", Type: proto.ColumnType_STRING, Description: "The pool ID to use when all other pools are detected as unhealthy."},
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Description: "Last modification time."},
			{Name: "proxied", Type: proto.ColumnType_BOOL, Description: "Whether the hostname gets Cloudflare's origin protection. Defaults to false."},
			{Name: "session_affinity", Type: proto.ColumnType_STRING, Description: "The session_affinity specifies the type of session affinity the loadbalancer should use unless specified as \"none\" or \"\"(default). The supported types are \"cookie\" and \"ip_cookie\"."},
			{Name: "session_affinity_ttl", Type: proto.ColumnType_INT, Description: "Time, in seconds, until this load balancers session affinity cookie expires after being created."},
			{Name: "steering_policy", Type: proto.ColumnType_STRING, Description: "Determine which method the load balancer uses to determine the fastest route to your origin. Valid values are: \"off\", \"geo\", \"dynamic_latency\", \"random\", \"proximity\" or \"\". Default is \"\"."},

			// JSON columns
			{Name: "default_pools", Type: proto.ColumnType_JSON, Description: "A list of pool IDs ordered by their failover priority. Pools defined here are used by default, or when region_pools are not configured for a given region."},
			{Name: "pop_pools", Type: proto.ColumnType_JSON, Description: "A mapping of Cloudflare PoP identifiers to a list of pool IDs (ordered by their failover priority) for the PoP (datacenter). This feature is only available to enterprise customers."},
			{Name: "region_pools", Type: proto.ColumnType_JSON, Description: "A mapping of region/country codes to a list of pool IDs (ordered by their failover priority) for the given region. Any regions not explicitly defined will fall back to using default_pools."},
			{Name: "session_affinity_attributes", Type: proto.ColumnType_JSON, Description: "session affinity cookie attributes."},
		}),
	}
}

func listLoadBalancers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	zoneID := h.Item.(cloudflare.Zone).ID

	conn, err := connect(ctx, d)
	if err != nil {
		logger.Error("listLoadBalancers", "connection_error", err)
		return nil, err
	}
	loadBalancers, err := conn.ListLoadBalancers(ctx, zoneID)
	if err != nil {
		logger.Error("ListLoadBalancers", "api error", err)
		return nil, err
	}
	for _, resource := range loadBalancers {
		d.StreamListItem(ctx, resource)
	}
	return nil, nil
}
