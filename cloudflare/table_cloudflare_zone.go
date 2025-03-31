package cloudflare

import (
	"context"
	"strings"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/zones"
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
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "Zone identifier tag."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The domain name."},

			// Other columns
			{Name: "account", Type: proto.ColumnType_JSON, Description: "Account information for the zone."},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the zone was created."},
			{Name: "development_mode", Type: proto.ColumnType_INT, Description: "The interval (in seconds) from when development mode expires (positive integer) or last expired (negative integer) for the domain. If development mode has never been enabled, this value is 0."},
			{Name: "dnssec", Type: proto.ColumnType_JSON, Hydrate: getZoneDNSSEC, Transform: transform.FromValue(), Description: "DNSSEC settings for the zone."},
			{Name: "meta", Type: proto.ColumnType_JSON, Description: "Metadata associated with the zone."},
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the zone was last modified."},
			{Name: "name_servers", Type: proto.ColumnType_JSON, Description: "Cloudflare-assigned name servers. This is only populated for zones that use Cloudflare DNS."},
			{Name: "original_dnshost", Type: proto.ColumnType_STRING, Description: "DNS host at the time of switching to Cloudflare."},
			{Name: "original_name_servers", Type: proto.ColumnType_JSON, Description: "Original name servers before moving to Cloudflare."},
			{Name: "original_registrar", Type: proto.ColumnType_STRING, Description: "Registrar for the domain at the time of switching to Cloudflare."},
			{Name: "owner", Type: proto.ColumnType_JSON, Description: "Information about the user or organization that owns the zone."},
			{Name: "paused", Type: proto.ColumnType_BOOL, Description: "Indicates if the zone is only using Cloudflare DNS services. A true value means the zone will not receive security or performance benefits."},
			// {Name: "permissions", Type: proto.ColumnType_JSON, Description: "Available permissions on the zone for the current user requesting the item."},
			{Name: "settings", Type: proto.ColumnType_JSON, Hydrate: getZoneSettings, Transform: transform.FromValue().Transform(settingsToStandard), Description: "Simple key value map of zone settings like advanced_ddos = on."},
			{Name: "plan", Type: proto.ColumnType_JSON, Hydrate: getZonePlan, Transform: transform.FromValue(), Description: "Current plan associated with the zone."},
			// {Name: "plan_pending", Type: proto.ColumnType_JSON, Description: "Pending plan change associated with the zone."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the zone."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "A full zone implies that DNS is hosted with Cloudflare. A partial zone is typically a partner-hosted zone or a CNAME setup."},
			{Name: "vanity_name_servers", Type: proto.ColumnType_JSON, Description: "Custom name servers for the zone."},
		}),
	}
}

func listZones(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("listZones", "connection_error", err)
		return nil, err
	}

	maxLimit := int32(500)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := zones.ZoneListParams{
		PerPage: cloudflare.F(float64(maxLimit)),
	}

	iter := conn.Zones.ListAutoPaging(ctx, input)
	for iter.Next() {
		zone := iter.Current()
		d.StreamListItem(ctx, zone)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	if err := iter.Err(); err != nil {
		logger.Error("listZones", "ListAutoPaging error", err)
		return nil, err
	}

	return nil, nil
}

func getZone(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connectV4(ctx, d)
	if err != nil {
		return nil, err
	}
	quals := d.EqualsQuals
	zoneID := quals["id"].GetStringValue()

	input := zones.ZoneGetParams{
		ZoneID: cloudflare.F(zoneID),
	}

	zone, err := conn.Zones.Get(ctx, input)
	if err != nil {
		return nil, err
	}
	return &zone, nil
}

func getZoneSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connectV4(ctx, d)
	if err != nil {
		return nil, err
	}

	zone := h.Item.(zones.Zone)
	// valid setting "aegis"?
	settingIds := []string{"0rtt", "advanced_ddos", "always_online", "always_use_https", "automatic_https_rewrites", "brotli", "browser_cache_ttl", "browser_check", "cache_level", "challenge_ttl", "ciphers", "cname_flattening", "development_mode", "early_hints", "edge_cache_ttl", "email_obfuscation", "h2_prioritization", "hotlink_protection", "http2", "http3", "image_resizing", "ip_geolocation", "ipv6", "max_upload", "min_tls_version", "mirage", "nel", "opportunistic_encryption", "opportunistic_onion", "orange_to_orange", "origin_error_page_pass_thru", "origin_h2_max_streams", "origin_max_http_version", "polish", "prefetch_preload", "privacy_pass", "proxy_read_timeout", "pseudo_ipv4", "replace_insecure_js", "response_buffering", "rocket_loader", "automatic_platform_optimization", "security_header", "security_level", "server_side_exclude", "sha1_support", "sort_query_string_for_cache", "ssl", "ssl_recommender", "tls_1_2_only", "tls_1_3", "tls_client_auth", "true_client_ip_header", "waf", "webp", "websockets"}
	settings := []zones.SettingGetResponse{}

	for _, settingId := range settingIds {
		item, err := conn.Zones.Settings.Get(ctx, settingId, zones.SettingGetParams{ZoneID: cloudflare.F(zone.ID)})
		if err != nil {
			// All of the setting could not be retrieved
			if strings.Contains(err.Error(), "Access denied") {
				return nil, nil
			}
			return nil, err
		}
		settings = append(settings, *item)
	}

	return settings, nil
}

func settingsToStandard(_ context.Context, d *transform.TransformData) (interface{}, error) {
	settings := d.HydrateItem.([]zones.SettingGetResponse)
	settingsMap := map[string]interface{}{}
	for _, i := range settings {
		settingsMap[string(i.ID)] = i.Value
	}
	return settingsMap, nil
}

func getZoneDNSSEC(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connectV4(ctx, d)
	if err != nil {
		return nil, err
	}
	zone := h.Item.(zones.Zone)
	input := dns.DNSSECGetParams{
		ZoneID: cloudflare.F(zone.ID),
	}
	dnssec, err := conn.DNS.DNSSEC.Get(ctx, input)
	if err != nil {
		return nil, err
	}
	return dnssec, nil
}

func getZonePlan(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connectV4(ctx, d)
	if err != nil {
		return nil, err
	}
	zone := h.Item.(zones.Zone)
	input := zones.PlanListParams{
		ZoneID: cloudflare.F(zone.ID),
	}
	var plans []zones.AvailableRatePlan
	iter := conn.Zones.Plans.ListAutoPaging(ctx, input)
	for iter.Next() {
		plan := iter.Current()
		plans = append(plans, plan)
	}
	if err != nil {
		return nil, err
	}
	return plans, nil
}
