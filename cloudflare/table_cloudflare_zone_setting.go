package cloudflare

import (
	"context"
	"encoding/json"
	"fmt"

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
		Get: &plugin.GetConfig{
			Hydrate:    getZoneSetting,
			KeyColumns: plugin.AllColumns([]string{"zone_id", "id"}),
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
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Transform: transform.From(transformModifiedOn), Description: "When the setting was last modified."},
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

	allSettings, err := ListAllZoneSettings(ctx, d, zone.ID)
	if err != nil {
		logger.Error("cloudflare_zone_setting.listZoneSettings", "error listing zone settings", err)
		return nil, err
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

func getZoneSetting(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	// Get zone_id and setting id from key columns
	zoneID := d.EqualsQualString("zone_id")
	settingID := d.EqualsQualString("id")

	// Get authenticated client using connectV4
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_zone_setting.getZoneSetting", "connection error", err)
		return nil, err
	}

	// Use the SDK function to get individual zone setting
	item, err := conn.Zones.Settings.Get(ctx, settingID, zones.SettingGetParams{ZoneID: cloudflare.F(zoneID)})
	if err != nil {
		logger.Error("cloudflare_zone_setting.getZoneSetting", "API error", err)
		return nil, err
	}

	setting := ZoneSettingInfo{
		ZoneID:             zoneID,
		SettingGetResponse: *item,
	}

	return setting, nil
}

//// TRANSFORM FUNCTIONS

func transformEditable(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	editable := d.HydrateItem.(ZoneSettingInfo).Editable
	if editable == zones.SettingGetResponseEditableTrue {
		return true, nil
	}
	return false, nil
}

func transformModifiedOn(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	modifiedOn := d.HydrateItem.(ZoneSettingInfo).ModifiedOn

	// Check if the time is zero/null (default Go time value or Unix epoch)
	if modifiedOn.IsZero() || modifiedOn.Unix() <= 0 {
		return nil, nil
	}

	// Check for the specific problematic timestamp pattern (year 0001)
	if modifiedOn.Year() <= 1 {
		return nil, nil
	}

	return modifiedOn, nil
}

//// API call

// ListAllZoneSettings makes a GET request to https://api.cloudflare.com/client/v4/zones/$ZONE_ID/settings
// Note: The main API documentation (https://developers.cloudflare.com/api/resources/zones/subresources/settings/methods/list/) lists this endpoint as DEPRECATED.
// However, according to the deprecation reference (https://developers.cloudflare.com/fundamentals/api/reference/deprecations/#2025-06-08), only the "cname_flattening" setting is affected; the endpoint itself is not fully deprecated.
func ListAllZoneSettings(ctx context.Context, d *plugin.QueryData, zoneID string) ([]ZoneSettingInfo, error) {
	logger := plugin.Logger(ctx)

	// Get authenticated client using connectV4
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_zone_setting.ListAllZoneSettings", "connection error", err)
		return nil, err
	}

	// Build the path for the API call
	path := fmt.Sprintf("zones/%s/settings", zoneID)

	// Define response structure
	var apiResponse struct {
		Success bool                     `json:"success"`
		Errors  []map[string]interface{} `json:"errors"`
		Result  []json.RawMessage        `json:"result"`
	}

	// Make the GET request using the client
	err = conn.Get(ctx, path, nil, &apiResponse)
	if err != nil {
		logger.Error("cloudflare_zone_setting.ListZoneSettings", "API request error", err)
		return nil, fmt.Errorf("failed to get zone settings: %w", err)
	}

	if !apiResponse.Success {
		return nil, fmt.Errorf("API request failed: %v", apiResponse.Errors)
	}

	// Process each setting from the result
	var settings []ZoneSettingInfo
	for _, rawSetting := range apiResponse.Result {
		// Parse each setting as a zones.SettingGetResponse
		var settingResponse zones.SettingGetResponse
		if err := json.Unmarshal(rawSetting, &settingResponse); err != nil {
			logger.Error("cloudflare_zone_setting.ListZoneSettings", "failed to parse setting", err, "raw_data", string(rawSetting))
			return nil, fmt.Errorf("failed to parse setting: %w", err)
		}

		settings = append(settings, ZoneSettingInfo{zoneID, settingResponse})
	}

	return settings, nil
}
