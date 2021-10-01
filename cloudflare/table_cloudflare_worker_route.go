package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableCloudflareWorkerRoute(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_worker_route",
		Description: "A Zone is a domain name along with its subdomains and other identities.",
		List: &plugin.ListConfig{
			Hydrate:       listWorkerRoutes,
			ParentHydrate: listZones,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "API item identifier tag."},
			{Name: "pattern", Type: proto.ColumnType_STRING, Description: "Patterns decide what (if any) script is matched based on the URL of that request."},
			{Name: "script", Type: proto.ColumnType_STRING, Description: "Size of the script, in bytes."},
			{Name: "zone_id", Description: "Specifies the zone identifier.", Type: proto.ColumnType_STRING, Hydrate: getParentZoneDetails, Transform: transform.FromField("ID")},
			{Name: "zone_name", Description: "Specifies the zone name.", Type: proto.ColumnType_STRING, Hydrate: getParentZoneDetails, Transform: transform.FromField("Name")},
		},
	}
}

func listWorkerRoutes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connect(ctx, d)
	if err != nil {
		logger.Error("listWorkerRoutes", "ListWorkerRoutes connect error", err)
		return nil, err
	}

	zoneDetails := h.Item.(cloudflare.Zone)
	resp, err := conn.ListWorkerRoutes(ctx, zoneDetails.ID)
	if err != nil {
		logger.Error("listWorkerRoutes", "ListWorkerRoutes api call error", err)
		return nil, err
	}
	for _, resource := range resp.Routes {
		d.StreamListItem(ctx, resource)
	}
	return nil, nil
}

func getParentZoneDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return h.ParentItem.(cloudflare.Zone), nil
}
