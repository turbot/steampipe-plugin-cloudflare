package cloudflare

import (
	"context"
	"strings"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/zones"
	"github.com/cloudflare/cloudflare-go/v4/leaked_credential_checks"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/cloudflare/cloudflare-go/v4/cache"
	"github.com/cloudflare/cloudflare-go/v4/argo"
	"github.com/cloudflare/cloudflare-go/v4/bot_management"
	"github.com/cloudflare/cloudflare-go/v4/security_txt"
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
			{Name: "permissions", Type: proto.ColumnType_JSON, Description: "Available permissions on the zone for the current user requesting the item.", Transform: transform.FromP(getExtraFieldPermissionsFromAPIresponse, "permissions")},
			{Name: "settings", Type: proto.ColumnType_JSON, Description: "[DEPRECATED] Simple key value map of zone settings like advanced_ddos = on. Use cloudflare_zone_setting table instead."},
			{Name: "plan", Type: proto.ColumnType_JSON, Hydrate: getZonePlan, Transform: transform.FromValue(), Description: "Current plan associated with the zone."},
			{Name: "plan_pending", Type: proto.ColumnType_JSON, Description: "[DEPRECATED] Pending plan change associated with the zone."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the zone."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "A full zone implies that DNS is hosted with Cloudflare. A partial zone is typically a partner-hosted zone or a CNAME setup."},
			{Name: "vanity_name_servers", Type: proto.ColumnType_JSON, Description: "Custom name servers for the zone."},
			{Name: "smart_tiered_cache", Type: proto.ColumnType_JSON, Hydrate: getSmartTieredCache, Transform: transform.FromValue(), Description: "Smart Tiered Cache settings for the zone."},
			{Name: "regional_tiered_cache", Type: proto.ColumnType_JSON, Hydrate: getRegionalTieredCache, Transform: transform.FromValue(), Description: "Regional Tiered Cache settings for the zone."},
			{Name: "argo_tiered_caching", Type: proto.ColumnType_JSON, Hydrate: getArgoTieredCaching, Transform: transform.FromValue(), Description: "Argo Tiered Caching settings for the zone."},
			{Name: "argo_smart_routing", Type: proto.ColumnType_JSON, Hydrate: getArgoSmartRouting, Transform: transform.FromValue(), Description: "Argo Smart Routing settings for the zone."},
			{Name: "bot_management", Type: proto.ColumnType_JSON, Hydrate: getBotManagement, Transform: transform.FromValue(), Description: "Bot management settings for the zone."},
			{Name: "security_txt", Type: proto.ColumnType_JSON, Hydrate: getSecurityTXT, Transform: transform.FromValue(), Description: "Security.txt configuration for the zone."},
			{Name: "leaked_credential_check", Type: proto.ColumnType_JSON, Hydrate: getLeakedCredentialCheck, Transform: transform.FromValue(), Description: "Leaked Credential Check configuration for the zone."},
		}),
	}
}

func listZones(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_zone.listZones", "connection_error", err)
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
		logger.Error("cloudflare_zone.listZones", "ListAutoPaging error", err)
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

func getSmartTieredCache(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connectV4(ctx, d)
	if err != nil {
		return nil, err
	}
	zone := h.Item.(zones.Zone)
	input := cache.SmartTieredCacheGetParams{
		ZoneID: cloudflare.F(zone.ID),
	}

	smartTieredCache, err := conn.Cache.SmartTieredCache.Get(ctx, input)
	if err != nil {
		return nil, err
	}
	return smartTieredCache, nil
}

func getRegionalTieredCache(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		return nil, err
	}
	zone := h.Item.(zones.Zone)
	input := cache.RegionalTieredCacheGetParams{
		ZoneID: cloudflare.F(zone.ID),
	}

	regionalTieredCache, err := conn.Cache.RegionalTieredCache.Get(ctx, input)
	if err != nil {
		// This setting might not be available for all zones
		if strings.Contains(err.Error(), "setting is not available") {
			return nil, nil
		}
		logger.Error("cloudflare_zone_setting.getRegionalTieredCache", "Regional tiered cache api error", err)
		return nil, nil
	}
	return regionalTieredCache, nil
}

func getArgoSmartRouting(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		return nil, err
	}
	zone := h.Item.(zones.Zone)
	input := argo.SmartRoutingGetParams{
		ZoneID: cloudflare.F(zone.ID),
	}

	argoSmartRouting, err := conn.Argo.SmartRouting.Get(ctx, input)
	if err != nil {
		// This setting might not be available for all zones
		if strings.Contains(err.Error(), "The request is not authorized to access this setting") {
			return nil, nil
		}
		logger.Error("cloudflare_zone_setting.getArgoSmartRouting", "Argo smart routing api error", err)
		return nil, nil
	}
	return argoSmartRouting, nil
}

func getArgoTieredCaching(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connectV4(ctx, d)
	if err != nil {
		return nil, err
	}
	zone := h.Item.(zones.Zone)
	input := argo.TieredCachingGetParams{
		ZoneID: cloudflare.F(zone.ID),
	}

	argoTieredCaching, err := conn.Argo.TieredCaching.Get(ctx, input)
	if err != nil {
		return nil, err
	}
	return argoTieredCaching, nil
}

func getBotManagement(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
    logger := plugin.Logger(ctx)
    conn, err := connectV4(ctx, d)
    if err != nil {
        return nil, err
    }

    zone := h.Item.(zones.Zone)
    params := bot_management.BotManagementGetParams{
        ZoneID: cloudflare.F(zone.ID),
    }

    resp, err := conn.BotManagement.Get(ctx, params)
    if err != nil {
        logger.Error("cloudflare_bot_management.getBotManagement", "API error", err)
        return nil, err
    }

	// The BotManagementGetResponse.AsUnion method is designed to return one of the following types:
	// BotFightModeConfiguration, SuperBotFightModeDefinitelyConfiguration, SuperBotFightModeLikelyConfiguration, or SubscriptionConfiguration,
	// depending on the subscription type for the zone's bot management settings.
	// However, due to a bug or incomplete implementation in the SDK, the method returns an incorrect or incomplete type that does not align 
	// with the expected schema. 
	// As a workaround, we directly return the raw JSON response, allowing the caller to manually unmarshal the data as needed.
    return resp.JSON.RawJSON(), nil
}

func getSecurityTXT(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connectV4(ctx, d)
	if err != nil {
		return nil, err
	}
	zone := h.Item.(zones.Zone)
	input := security_txt.SecurityTXTGetParams{
		ZoneID: cloudflare.F(zone.ID),
	}

	security_txt, err := conn.SecurityTXT.Get(ctx, input)
	if err != nil {
		return nil, err
	}
	return security_txt, nil
}

func getLeakedCredentialCheck(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connectV4(ctx, d)
	if err != nil {
		return nil, err
	}
	zone := h.Item.(zones.Zone)
	input := leaked_credential_checks.LeakedCredentialCheckGetParams{
		ZoneID: cloudflare.F(zone.ID),
	}

	leaked_credential_check, err := conn.LeakedCredentialChecks.Get(ctx, input)
	if err != nil {
		return nil, err
	}
	return leaked_credential_check, nil
}

//// TRANSFORM FUNCTIONS

func getExtraFieldPermissionsFromAPIresponse(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	response := d.HydrateItem.(zones.Zone)
	param := d.Param.(string)

	extraFields, err := toMap(response.JSON.RawJSON())
	if err != nil {
		return nil, err
	}

	return extraFields[param], nil
}
