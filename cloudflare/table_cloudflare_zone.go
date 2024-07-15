package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableCloudflareZone(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_zone",
		Description: "A Zone is a domain name along with its subdomains and other identities.",
		List: &plugin.ListConfig{
			Hydrate: listZones,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"Invalid zone identifier"}),
			Hydrate:           getZone,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Zone identifier tag."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The domain name."},

			// Other columns
			// TODO - do we need this here {Name: "account", Type: proto.ColumnType_JSON, Description: "TODO"},
			{Name: "betas", Type: proto.ColumnType_JSON, Description: "Beta feature flags associated with the zone."},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the zone was created."},
			{Name: "deactivation_reason", Type: proto.ColumnType_STRING, Description: "TODO"},
			{Name: "development_mode", Type: proto.ColumnType_INT, Description: "The interval (in seconds) from when development mode expires (positive integer) or last expired (negative integer) for the domain. If development mode has never been enabled, this value is 0."},
			{Name: "dnssec", Type: proto.ColumnType_JSON, Hydrate: getZoneDNSSEC, Transform: transform.FromValue(), Description: "DNSSEC settings for the zone."},
			{Name: "host", Type: proto.ColumnType_JSON, Description: "TODO"},
			{Name: "meta", Type: proto.ColumnType_JSON, Description: "Metadata associated with the zone."},
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the zone was last modified."},
			{Name: "name_servers", Type: proto.ColumnType_JSON, Description: "Cloudflare-assigned name servers. This is only populated for zones that use Cloudflare DNS."},
			{Name: "original_dnshost", Type: proto.ColumnType_STRING, Description: "DNS host at the time of switching to Cloudflare."},
			{Name: "original_name_servers", Type: proto.ColumnType_JSON, Description: "Original name servers before moving to Cloudflare."},
			{Name: "original_registrar", Type: proto.ColumnType_STRING, Description: "Registrar for the domain at the time of switching to Cloudflare."},
			{Name: "owner", Type: proto.ColumnType_JSON, Description: "Information about the user or organization that owns the zone."},
			{Name: "paused", Type: proto.ColumnType_BOOL, Description: "Indicates if the zone is only using Cloudflare DNS services. A true value means the zone will not receive security or performance benefits."},
			{Name: "permissions", Type: proto.ColumnType_JSON, Description: "Available permissions on the zone for the current user requesting the item."},
			{Name: "settings", Type: proto.ColumnType_JSON, Hydrate: getZoneSettings, Transform: transform.FromValue().Transform(settingsToStandard), Description: "Simple key value map of zone settings like advanced_ddos = on. Full settings details are in settings_src."},
			//{Name: "settings_src", Type: proto.ColumnType_JSON, Hydrate: getZoneSettings, Transform: transform.FromValue(), Description: "Original source form of zone settings for caching, security and other features of Cloudflare."},
			{Name: "plan", Type: proto.ColumnType_JSON, Description: "Current plan associated with the zone."},
			{Name: "plan_pending", Type: proto.ColumnType_JSON, Description: "Pending plan change associated with the zone."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the zone."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "A full zone implies that DNS is hosted with Cloudflare. A partial zone is typically a partner-hosted zone or a CNAME setup."},
			//{Name: "universal_ssl_settings", Type: proto.ColumnType_JSON, Hydrate: getZoneUniversalSSLSettings, Transform: transform.FromValue(), Description: "Universal SSL settings for a zone."},
			{Name: "vanity_name_servers", Type: proto.ColumnType_JSON, Description: "Custom name servers for the zone."},
			// TODO - It's unclear when this is set {Name: "verification_key", Type: proto.ColumnType_STRING, Description: "TODO"},
		}),
	}
}

func listZones(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connect(ctx, d)
	if err != nil {
		logger.Error("listZones", "connection_error", err)
		return nil, err
	}
	resp, err := conn.ListZonesContext(ctx)
	if err != nil {
		logger.Error("listZones", "ListZonesContext api error", err)
		return nil, err
	}
	for _, i := range resp.Result {
		d.StreamListItem(ctx, i)
	}
	return nil, nil
}

func getZone(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	quals := d.EqualsQuals
	zoneID := quals["id"].GetStringValue()
	item, err := conn.ZoneDetails(ctx, zoneID)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func getZoneSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	zone := h.Item.(cloudflare.Zone)
	item, err := conn.ZoneSettings(ctx, zone.ID)
	if err != nil {
		return nil, err
	}
	return item.Result, nil
}

func settingsToStandard(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	settings := d.HydrateItem.([]cloudflare.ZoneSetting)
	// Convert the settings into a map, which makes them a lot easier to query by name
	settingsMap := map[string]interface{}{}
	for _, i := range settings {
		settingsMap[i.ID] = i.Value
	}
	return settingsMap, nil
}

func getZoneDNSSEC(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	zone := h.Item.(cloudflare.Zone)
	item, err := conn.ZoneDNSSECSetting(ctx, zone.ID)
	if err != nil {
		return nil, err
	}
	return item, nil
}
