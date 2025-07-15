package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/workers"
	"github.com/cloudflare/cloudflare-go/v4/zones"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableCloudflareWorkerRoute(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_worker_route",
		Description: "Routes are basic patterns used to enable or disable workers that match requests.",
		List: &plugin.ListConfig{
			Hydrate:       listWorkerRoutes,
			ParentHydrate: listZones,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "zone_id", Require: plugin.Optional},
			},
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "API item identifier tag."},
			{Name: "zone_name", Type: proto.ColumnType_STRING, Hydrate: getParentZoneDetails, Transform: transform.FromField("Name"), Description: "Specifies the zone name."},
			{Name: "pattern", Type: proto.ColumnType_STRING, Description: "Patterns decide what (if any) script is matched based on the URL of that request."},
			{Name: "script", Type: proto.ColumnType_STRING, Description: "Name of the script to apply when the route is matched. The route is skipped when this is blank/missing."},
			{Name: "zone_id", Type: proto.ColumnType_STRING, Hydrate: getParentZoneDetails, Transform: transform.FromField("ID"), Description: "Specifies the zone identifier."},
		}),
	}
}

func listWorkerRoutes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	zoneDetails := h.Item.(zones.Zone)

	inputZoneId := d.EqualsQualString("zone_id")

	// Only list routes for zones stated in the input query
	if inputZoneId != "" && inputZoneId != zoneDetails.ID {
		return nil, nil
	}

	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_worker_route.listWorkerRoutes", "connect error", err)
		return nil, err
	}

	input := workers.RouteListParams{
		ZoneID: cloudflare.F(zoneDetails.ID),
	}

	iter := conn.Workers.Routes.ListAutoPaging(ctx, input)
	if err := iter.Err(); err != nil {
		logger.Error("cloudflare_worker_route.listWorkerRoutes", "api call error", err)
		return nil, err
	}

	for iter.Next() {
		resource := iter.Current()
		d.StreamListItem(ctx, resource)
	}
	return nil, nil
}

func getParentZoneDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return h.ParentItem.(zones.Zone), nil
}
