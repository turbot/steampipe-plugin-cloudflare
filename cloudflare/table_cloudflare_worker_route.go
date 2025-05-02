package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/workers"
	"github.com/cloudflare/cloudflare-go/v4/zones"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableCloudflareWorkerRoute() *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_worker_route",
		Description: "Routes are basic patterns used to enable or disable workers that match requests.",
		List: &plugin.ListConfig{
			Hydrate:       listWorkerRoutes,
			ParentHydrate: TableCloudflareZone().List.Hydrate,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "zone_id", Require: plugin.Optional},
			},
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "API item identifier tag."},
			{Name: "zone_name", Type: proto.ColumnType_STRING, Hydrate: getParentZoneDetails, Transform: transform.FromField("Name"), Description: "Specifies the zone name."},
			{Name: "pattern", Type: proto.ColumnType_STRING, Description: "Patterns decide what (if any) script is matched based on the URL of that request."},
			{Name: "script", Type: proto.ColumnType_STRING, Description: "Name of the script to apply when the route is matched. The route is skipped when this is blank/missing."},
			{Name: "zone_id", Type: proto.ColumnType_STRING, Hydrate: getParentZoneDetails, Transform: transform.FromField("ID"), Description: "Specifies the zone identifier."},
		}),
	}
}

func listWorkerRoutes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	zoneData := h.Item.(*zones.Zone)

	inputZoneId := d.EqualsQualString("zone_id")

	// Only list routes for zones stated in the input query
	if inputZoneId != "" && inputZoneId != zoneData.ID {
		return nil, nil
	}

	conn, err := connect(ctx, d)
	if err != nil {
		logger.Error("listWorkerRoutes", "connect error", err)
		return nil, err
	}

	// Get the client's options from the connection
	opts := []option.RequestOption{}
	if conn != nil {
		opts = append(opts, conn.Options...)
	}

	routeService := workers.NewRouteService(opts...)
	result, err := routeService.List(ctx, workers.RouteListParams{
		ZoneID: cloudflare.F(zoneData.ID),
	})
	if err != nil {
		logger.Error("listWorkerRoutes", "api call error", err)
		return nil, err
	}

	for _, route := range result.Result {
		d.StreamListItem(ctx, route)
	}
	return nil, nil
}

func getParentZoneDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return h.ParentItem.(*zones.Zone), nil
}
