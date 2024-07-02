package cloudflare

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableCloudflareLoadBalancerMonitor(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_load_balancer_monitor",
		Description: "A monitor issues health checks at regular intervals to evaluate the health of an origin pool.",
		List: &plugin.ListConfig{
			Hydrate: listLoadBalancerMonitors,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Load balancer monitor API item identifier."},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the load balancer monitor was created."},
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the load balancer monitor was last modified."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The protocol to use for the healthcheck. Currently supported protocols are \"HTTP\", \"HTTPS\" and \"TCP\". Default: \"http\"."},

			// Other columns
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Monitor description."},
			{Name: "method", Type: proto.ColumnType_STRING, Description: "The method to use for the health check. Valid values are any valid HTTP verb if type is \"http\" or \"https\", or connection_established if type is \"tcp\". Default: \"GET\" if type is \"http\" or \"https\", or \"connection_established\" if type is \"tcp\" ."},
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The endpoint path to health check against. Default: \"/\". Only valid if type is \"http\" or \"https\"."},
			{Name: "header", Type: proto.ColumnType_JSON, Description: "The HTTP request headers to send in the health check. It is recommended you set a Host header by default. The User-Agent header cannot be overridden. Fields documented below. Only valid if type is \"http\" or \"https\"."},
			{Name: "timeout", Type: proto.ColumnType_INT, Description: "The timeout (in seconds) before marking the health check as failed. Default: 5."},
			{Name: "retries", Type: proto.ColumnType_INT, Description: "The number of retries to attempt in case of a timeout before marking the origin as unhealthy. Retries are attempted immediately. Default: 2."},
			{Name: "interval", Type: proto.ColumnType_INT, Description: "The interval between each health check. Shorter intervals may improve failover time, but will increase load on the origins as we check from multiple locations. Default: 60."},
			{Name: "port", Type: proto.ColumnType_INT, Description: "The port number to use for the healthcheck, required when creating a TCP monitor. Valid values are in the range 0-65535"},
			{Name: "expected_body", Type: proto.ColumnType_STRING, Description: "A case-insensitive sub-string to look for in the response body. If this string is not found, the origin will be marked as unhealthy. Only valid if type is \"http\" or \"https\". Default: \"\"."},
			{Name: "expected_codes", Type: proto.ColumnType_STRING, Description: "The expected HTTP response code or code range of the health check. Eg 2xx. Only valid and required if type is \"http\" or \"https\"."},
			{Name: "follow_redirects", Type: proto.ColumnType_BOOL, Description: "Follow redirects if returned by the origin. Only valid if type is \"http\" or \"https\"."},
			{Name: "allow_insecure", Type: proto.ColumnType_BOOL, Description: "Do not validate the certificate when monitor use HTTPS. Only valid if type is \"http\" or \"https\"."},
			{Name: "probe_zone", Type: proto.ColumnType_STRING, Description: "Assign this monitor to emulate the specified zone while probing. Only valid if type is \"http\" or \"https\"."},
		}),
	}
}

func listLoadBalancerMonitors(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	conn, err := connect(ctx, d)
	if err != nil {
		logger.Error("listLoadBalancers", "connection_error", err)
		return nil, err
	}
	// Paging not supported by rest api
	loadBalancersPools, err := conn.ListLoadBalancerMonitors(ctx)
	if err != nil {
		logger.Error("ListLoadBalancers", "api error", err)
		return nil, err
	}
	for _, resource := range loadBalancersPools {
		d.StreamListItem(ctx, resource)
	}
	return nil, nil
}
