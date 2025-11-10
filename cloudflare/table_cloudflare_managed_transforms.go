package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/managed_transforms"
	"github.com/cloudflare/cloudflare-go/v4/zones"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableCloudflareManagedTransforms(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_managed_transforms",
		Description: "Managed Transforms allow you to perform common adjustments to HTTP request and response headers with the click of a button.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "zone_id", Require: plugin.Optional},
			},
			Hydrate: listManagedTransforms,
			ParentHydrate: listZones,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "The human-readable identifier of the Managed Transform."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "The type of Managed Transform: 'request_header' or 'response_header'."},
			{Name: "enabled", Type: proto.ColumnType_BOOL, Description: "Whether the Managed Transform is enabled."},
			{Name: "has_conflict", Type: proto.ColumnType_BOOL, Description: "Whether the Managed Transform conflicts with the currently-enabled Managed Transforms."},
			
			// Other columns
			{Name: "zone_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ZoneID"), Description: "The zone ID."},
			
			// JSON Columns
			{Name: "conflicts_with", Type: proto.ColumnType_JSON, Description: "The Managed Transforms that this Managed Transform conflicts with."},
		}),
	}
}

type ManagedTransformInfo struct {
	ID            string
	Type          string
	Enabled       bool
	HasConflict   bool
	ConflictsWith []string
	ZoneID        string
}

//// LIST FUNCTION

// listManagedTransforms retrieves all managed transforms.
func listManagedTransforms(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	zoneDetails := h.Item.(zones.Zone)
	inputZoneId := d.EqualsQualString("zone_id")

	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_managed_transforms.listManagedTransforms", "connection_error", err)
		return nil, err
	}

	// Only list managed transforms for zones stated in the input query
	if inputZoneId != "" && inputZoneId != zoneDetails.ID {
		return nil, nil
	}

	// Build API parameters
	input := managed_transforms.ManagedTransformListParams{
		ZoneID: cloudflare.F(zoneDetails.ID),
	}

	// Execute API call (this is a GET operation that returns both lists)
	response, err := conn.ManagedTransforms.List(ctx, input)
	if err != nil {
		logger.Error("cloudflare_managed_transforms.listManagedTransforms", "api_error", err)
		return nil, err
	}

	// Process managed request headers
	for _, requestHeader := range response.ManagedRequestHeaders {
		transform := ManagedTransformInfo{
			ID:            requestHeader.ID,
			Type:          "request_header",
			Enabled:       requestHeader.Enabled,
			HasConflict:   requestHeader.HasConflict,
			ConflictsWith: requestHeader.ConflictsWith,
			ZoneID:        zoneDetails.ID,
		}
		d.StreamListItem(ctx, transform)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	// Process managed response headers
	for _, responseHeader := range response.ManagedResponseHeaders {
		transform := ManagedTransformInfo{
			ID:            responseHeader.ID,
			Type:          "response_header",
			Enabled:       responseHeader.Enabled,
			HasConflict:   responseHeader.HasConflict,
			ConflictsWith: responseHeader.ConflictsWith,
			ZoneID:        zoneDetails.ID,
		}
		d.StreamListItem(ctx, transform)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
