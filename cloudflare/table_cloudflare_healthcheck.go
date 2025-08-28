package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/healthchecks"
	"github.com/cloudflare/cloudflare-go/v4/zones"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableCloudflareHealthcheck(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_healthcheck",
		Description: "A Health Check is a service that runs on Cloudflare’s edge network to monitor whether an origin server is online.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "zone_id", Require: plugin.Optional},
			},
			Hydrate: listHealthchecks,
			ParentHydrate: listZones,
		},
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "id", Require: plugin.Required},
				{Name: "zone_id", Require: plugin.Required},
			},
			ShouldIgnoreError: isNotFoundError([]string{"Invalid healthcheck identifier"}),
			Hydrate:           getHealthcheck,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "Healthcheck identifier."},
			{Name: "address", Type: proto.ColumnType_STRING, Description: "The hostname or IP address of the origin server to run health checks on."},
			{Name: "consecutive_fails", Type: proto.ColumnType_INT, Description: "The number of consecutive fails required from a health check before changing the health to unhealthy."},
			{Name: "consecutive_successes", Type: proto.ColumnType_INT, Description: "The number of consecutive successes required from a health check before changing the health to healthy."},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the healthcheck was created."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "A human-readable description of the health check."},
			{Name: "failure_reason", Type: proto.ColumnType_STRING, Description: "The current failure reason if status is unhealthy."},
			{Name: "interval", Type: proto.ColumnType_INT, Description: "The interval between each health check."},
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the healthcheck was last modified."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Healthcheck name."},
			{Name: "retries", Type: proto.ColumnType_INT, Description: "The number of retries to attempt in case of a timeout before marking the origin as unhealthy. Retries are attempted immediately."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The current status of the origin server according to the health check."},
			{Name: "suspended", Type: proto.ColumnType_BOOL, Description: "If suspended, no health checks are sent to the origin."},
			{Name: "timeout", Type: proto.ColumnType_INT, Description: "The timeout (in seconds) before marking the health check as failed."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The protocol to use for the health check. Currently supported protocols are 'HTTP', 'HTTPS' and 'TCP'."},
			
			// Query columns for filtering
			{Name: "zone_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ZoneID"), Description: "The zone ID to filter healthchecks."},
		
			// JSON Columns
			{Name: "tcp_config", Type: proto.ColumnType_JSON, Transform: transform.FromField("TCPConfig"), Description: "Parameters specific to TCP health check."},
			{Name: "http_config", Type: proto.ColumnType_JSON, Transform: transform.FromField("HTTPConfig"), Description: "Parameters specific to an HTTP or HTTPS health check."},
			{Name: "check_regions", Type: proto.ColumnType_JSON, Description: "A list of regions from which to run health checks. Null means Cloudflare will pick a default region."},
		}),
	}
}

type HealthcheckInfo struct {
	ZoneID string
	healthchecks.Healthcheck
}

//// LIST FUNCTION

// listHealthchecks retrieves all healthchecks.
func listHealthchecks(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	zoneDetails := h.Item.(zones.Zone)
	inputZoneId := d.EqualsQualString("zone_id")

	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_healthcheck.listHealthchecks", "connection_error", err)
		return nil, err
	}

	// Only list healthchecks for zones stated in the input query
	if inputZoneId != "" && inputZoneId != zoneDetails.ID {
		return nil, nil
	}

	// Build API parameters
	input := healthchecks.HealthcheckListParams{
		ZoneID: cloudflare.F(zoneDetails.ID),
	}

	// Execute paginated API call
	iter := conn.Healthchecks.ListAutoPaging(ctx, input)
	for iter.Next() {
		current := iter.Current()

		healthcheck := HealthcheckInfo{
			ZoneID:				zoneDetails.ID,
			Healthcheck:		current,
		}
		d.StreamListItem(ctx, healthcheck)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	if err := iter.Err(); err != nil {
		logger.Error("cloudflare_healthcheck.listHealthchecks", "api_error", err)
		return nil, err
	}

	return nil, nil
}

//// GET FUNCTION

// getHealthcheck retrieves a specific healthcheck by ID.
//
// Parameters:
// - id: The healthcheck identifier (required)
// - zone_id: The zone context (required)
func getHealthcheck(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_healthcheck.getHealthcheck", "connection_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	healthcheckID := quals["id"].GetStringValue()
	zoneID := quals["zone_id"].GetStringValue()

	input := healthchecks.HealthcheckGetParams{
		ZoneID: cloudflare.F(zoneID),
	}

	// Execute API call to get the specific healthcheck
	item, err := conn.Healthchecks.Get(ctx, healthcheckID, input)
	if err != nil {
		logger.Error("cloudflare_healthcheck.getHealthcheck", "api_error", err)
		return nil, err
	}

	healthcheck := HealthcheckInfo{
		ZoneID:				zoneID,
		Healthcheck:		*item,
	}

	return healthcheck, nil
}
