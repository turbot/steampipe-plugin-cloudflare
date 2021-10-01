package cloudflare

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableCloudflareWorkerScript(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_worker_script",
		Description: "A Zone is a domain name along with its subdomains and other identities.",
		List: &plugin.ListConfig{
			Hydrate: listWorkers,
		},
		// Get: &plugin.GetConfig{
		// 	KeyColumns:        plugin.SingleColumn("id"),
		// 	ShouldIgnoreError: isNotFoundError([]string{"Invalid zone identifier"}),
		// 	Hydrate:           getZone,
		// },
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Zone identifier tag."},
			{Name: "etag", Type: proto.ColumnType_STRING, Description: "Hashed script content, can be used in a If-None-Match header when updating."},
			{Name: "size", Type: proto.ColumnType_INT, Description: "Size of the script, in bytes."},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "The domain name."},
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the script was last modified."},
		},
	}
}

func listWorkers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	resp, err := conn.ListWorkerScripts(ctx)
	if err != nil {
		return nil, err
	}
	for index, resource := range resp.WorkerList {
		logger.Info("listWorkers", index, resource)
		d.StreamListItem(ctx, resource)
	}
	return nil, nil
}

// func getZone(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
// 	conn, err := connect(ctx, d)
// 	if err != nil {
// 		return nil, err
// 	}
// 	quals := d.KeyColumnQuals
// 	zoneID := quals["id"].GetStringValue()
// 	item, err := conn.ZoneDetails(ctx, zoneID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return item, nil
// }

// func getZoneSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	conn, err := connect(ctx, d)
// 	if err != nil {
// 		return nil, err
// 	}
// 	zone := h.Item.(cloudflare.Zone)
// 	item, err := conn.ZoneSettings(ctx, zone.ID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return item.Result, nil
// }

// func settingsToStandard(ctx context.Context, d *transform.TransformData) (interface{}, error) {
// 	settings := d.HydrateItem.([]cloudflare.ZoneSetting)
// 	// Convert the settings into a map, which makes them a lot easier to query by name
// 	settingsMap := map[string]interface{}{}
// 	for _, i := range settings {
// 		settingsMap[i.ID] = i.Value
// 	}
// 	return settingsMap, nil
// }

// func getZoneDNSSEC(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	conn, err := connect(ctx, d)
// 	if err != nil {
// 		return nil, err
// 	}
// 	zone := h.Item.(cloudflare.Zone)
// 	item, err := conn.ZoneDNSSECSetting(ctx, zone.ID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return item, nil
// }

// func getZoneUniversalSSLSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
// 	conn, err := connect(ctx, d)
// 	if err != nil {
// 		return nil, err
// 	}
// 	zone := h.Item.(cloudflare.Zone)
// 	item, err := conn.UniversalSSLSettingDetails(ctx, zone.ID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return item, nil
// }
