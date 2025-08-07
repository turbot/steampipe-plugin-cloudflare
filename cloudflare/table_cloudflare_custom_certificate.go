package cloudflare

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/custom_certificates"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableCloudflareCustomCertificate(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_custom_certificate",
		Description: "Custom certificates are meant for Business and Enterprise customers who want to use their own SSL certificates.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "zone_id", Require: plugin.Required},
			},
			Hydrate: listCustomCertificates,
		},
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "id", Require: plugin.Required},
				{Name: "zone_id", Require: plugin.Required},
			},
			ShouldIgnoreError: isNotFoundError([]string{"Invalid custom certificate identifier"}),
			Hydrate:           getCustomCertificate,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "Custom certificate identifier."},
			{Name: "bundle_method", Type: proto.ColumnType_STRING, Description: "The method used to build the SSL certificate chain."},
			{Name: "expires_on", Type: proto.ColumnType_TIMESTAMP,  Description: "When the certificate from the authority expires."},
			{Name: "hosts", Type: proto.ColumnType_STRING, Description: "The domain names covered by the custom certificate."},
			{Name: "issuer", Type: proto.ColumnType_STRING, Description: "The certificate authority that issued the certificate."},
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the certificate was last modified."},
			{Name: "priority", Type: proto.ColumnType_STRING, Description: "The order/priority in which the certificate will be used in a request."},
			{Name: "signature", Type: proto.ColumnType_STRING, Description: "The type of hash used for the certificate."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the zone's custom SSL."},
			{Name: "uploaded_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the certificate was uploaded to Cloudflare."},
			{Name: "geo_restrictions", Type: proto.ColumnType_STRING,Transform: transform.From(transformGeoRestrictions), Description: "Specify the region where your private key can be held locally for optimal TLS performance."},
			{Name: "policy", Type: proto.ColumnType_STRING, Description: "Specify the policy that determines the region where your private key will be held locally."},
			
			// Query columns for filtering
			{Name: "zone_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ZoneID"), Description: "The zone ID to filter custom certificates."},
		
			// JSON Columns
			{Name: "keyless_server", Type: proto.ColumnType_JSON, Description: "Keyless certificate details."},
		}),
	}
}

//// LIST FUNCTION

// listCustomCertificates retrieves all custom certificates for the specified zone_id.
func listCustomCertificates(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_custom_certificate.listCustomCertificates", "connection_error", err)
		return nil, err
	}

	// Get the qualifiers
	quals := d.EqualsQuals
	zoneID := quals["zone_id"].GetStringValue()

	// Empty check
	if zoneID == "" {
		return nil, nil
	}

	// Build API parameters
	input := custom_certificates.CustomCertificateListParams{
		ZoneID: cloudflare.F(zoneID),
	}

	// Execute paginated API call
	iter := conn.CustomCertificates.ListAutoPaging(ctx, input)
	for iter.Next() {
		customCertificate := iter.Current()
		d.StreamListItem(ctx, customCertificate)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	if err := iter.Err(); err != nil {
		// Custom certificates are not available for all plan levels.
		if strings.Contains(err.Error(), "Plan level does not allow custom certificates") {
			return nil, nil
		}
		logger.Error("cloudflare_custom_certificate.listCustomCertificates", "ListAutoPaging error", err)
		return nil, err
	}

	return nil, nil
}

//// GET FUNCTION

// getCustomCertificate retrieves a specific custom certificate by ID.
//
// Parameters:
// - id: The custom certificate identifier (required)
// - zone_id: The zone context (required)
func getCustomCertificate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_custom_certificate.getCustomCertificate", "connection_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	customCertificateID := quals["id"].GetStringValue()
	zoneID := quals["zone_id"].GetStringValue()

	input := custom_certificates.CustomCertificateGetParams{
		ZoneID: cloudflare.F(zoneID),
	}

	// Execute API call to get the specific custom certificate
	customCertificate, err := conn.CustomCertificates.Get(ctx, customCertificateID, input)
	if err != nil {
		logger.Error("cloudflare_custom_certificate.getCustomCertificate", "error", err)
		return nil, err
	}

	return customCertificate, nil
}

//// TRANSFORM FUNCTIONS

func transformGeoRestrictions(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	logger := plugin.Logger(ctx)
    switch customCertificate := d.HydrateItem.(type) {
    case custom_certificates.CustomCertificate:
        return string(customCertificate.GeoRestrictions.Label), nil
    case *custom_certificates.CustomCertificate:
        if customCertificate == nil {
            return nil, nil
        }
        return string(customCertificate.GeoRestrictions.Label), nil
    default:
		err := fmt.Errorf("unexpected type: %T", d.HydrateItem)
		logger.Error("cloudflare_custom_certificate.transformGeoRestrictions", "customCertificate unexpected type", err)
        return nil, err
    }
}
