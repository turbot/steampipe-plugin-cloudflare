package cloudflare

import (
	"context"
	"strings"
	"sync"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/zones"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableCloudflareZoneSetting(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_zone_setting",
		Description: "Individual zone settings that control various Cloudflare features for a zone.",
		List: &plugin.ListConfig{
			Hydrate:       listZoneSettings,
			ParentHydrate: listZones,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "zone_id", Require: plugin.Optional},
				{Name: "id", Require: plugin.Optional},
			},
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "The ID of the zone setting."},
			{Name: "zone_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ZoneID"), Description: "Zone identifier."},

			// The value field from the API response may have string or JSON property.
			// for couple of settings like security_header, automatic_platform_optimization, ciphers, etc... we will have JSON value
			// rest are string, integer type.
			// So we have set the data type as string
			{Name: "value", Type: proto.ColumnType_STRING, Description: "The current value of the zone setting."},

			// Other columns
			{Name: "editable", Type: proto.ColumnType_BOOL, Transform: transform.From(transformEditable), Description: "Whether the setting is editable."},
			{Name: "enabled", Type: proto.ColumnType_BOOL, Description: "SSL-recommender enrollment setting."},
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the setting was last modified."},
		}),
	}
}

type ZoneSettingInfo struct {
	ZoneID string
	zones.SettingGetResponse
}

func listZoneSettings(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	zone := h.Item.(zones.Zone)

	// Only list settings for zones stated in the input query
	inputZoneId := d.EqualsQualString("zone_id")
	if inputZoneId != "" && inputZoneId != zone.ID {
		return nil, nil
	}

	// Check if specific setting is requested
	inputSettingId := d.EqualsQualString("id")

	// Available zone setting IDs
	settingIds := []string{
		"aegis", "0rtt", "advanced_ddos", "always_online", "always_use_https", "automatic_https_rewrites",
		"brotli", "browser_cache_ttl", "browser_check", "cache_level", "challenge_ttl",
		"ciphers", "cname_flattening", "development_mode", "early_hints", "edge_cache_ttl",
		"email_obfuscation", "h2_prioritization", "hotlink_protection", "http2", "http3",
		"image_resizing", "ip_geolocation", "ipv6", "max_upload", "min_tls_version",
		"mirage", "nel", "opportunistic_encryption", "opportunistic_onion", "orange_to_orange",
		"origin_error_page_pass_thru", "origin_h2_max_streams", "origin_max_http_version",
		"polish", "prefetch_preload", "privacy_pass", "proxy_read_timeout", "pseudo_ipv4",
		"replace_insecure_js", "response_buffering", "rocket_loader", "automatic_platform_optimization",
		"security_header", "security_level", "server_side_exclude", "sha1_support",
		"sort_query_string_for_cache", "ssl", "ssl_recommender", "tls_1_2_only", "tls_1_3",
		"tls_client_auth", "true_client_ip_header", "waf", "webp", "websockets",
	}

	// If specific setting ID is requested, only query that one
	if inputSettingId != "" {
		settingIds = []string{inputSettingId}
	}

	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_zone_setting.listZoneSettings", "connection error", err)
		return nil, err
	}

	// Process settings in batches of 10 with parallel API calls
	// This is to avoid rate limiting by Cloudflare and optimize the query timing
	batchSize := 10
	var allSettings []ZoneSettingInfo

	for i := 0; i < len(settingIds); i += batchSize {
		end := i + batchSize
		if end > len(settingIds) {
			end = len(settingIds)
		}

		batch := settingIds[i:end]
		var wg sync.WaitGroup
		var mu sync.Mutex
		var batchSettings []ZoneSettingInfo

		for _, settingId := range batch {
			wg.Add(1)
			go func(sid string) {
				defer wg.Done()

				item, err := conn.Zones.Settings.Get(ctx, sid, zones.SettingGetParams{ZoneID: cloudflare.F(zone.ID)})
				if err != nil {
					// Some settings might not be available for all zones
					if strings.Contains(err.Error(), "Undefined zone setting") {
						return
					}
					if strings.Contains(err.Error(), "Access denied") {
						return
					}
					logger.Error("cloudflare_zone_setting.listZoneSettings", "ZoneSetting api error", err)
					return
				}

				setting := ZoneSettingInfo{zone.ID, *item}

				mu.Lock()
				batchSettings = append(batchSettings, setting)
				mu.Unlock()
			}(settingId)
		}

		wg.Wait()
		allSettings = append(allSettings, batchSettings...)

		// Check if context is cancelled or limit reached
		if d.RowsRemaining(ctx) == 0 {
			break
		}
	}

	// Stream all collected settings
	for _, setting := range allSettings {
		d.StreamLeafListItem(ctx, setting)

		// Check if context is cancelled or limit reached
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// TRANSFORM FUNCTIONS

func transformEditable(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	editable := d.HydrateItem.(ZoneSettingInfo).Editable
	if editable == zones.SettingGetResponseEditableTrue {
		return true, nil
	}
	return false, nil
}
