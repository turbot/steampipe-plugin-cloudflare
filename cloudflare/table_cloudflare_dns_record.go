package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/zones"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

func tableCloudflareDNSRecord(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_dns_record",
		Description: "DNS records for a zone.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "zone_id", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
			},
			Hydrate: listDNSRecord,
			ParentHydrate: listZones,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"zone_id", "id"}),
			ShouldIgnoreError: isNotFoundError([]string{"HTTP status 404"}),
			Hydrate:           getDNSRecord,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "zone_id", Type: proto.ColumnType_STRING, Description: "Zone where the record is defined.", Transform: transform.FromField("ZoneID")},
			{Name: "zone_name", Type: proto.ColumnType_STRING, Description: "[Deprecated] Name of the zone where the record is defined."},
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "ID of the record."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Type of the record (e.g. A, MX, CNAME)."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Domain name for the record (e.g. steampipe.io)."},
			{Name: "content", Type: proto.ColumnType_STRING, Description: "Content or value of the record. Changes by type, including IP address for A records and domain for CNAME records."},
			
			// Property ttl is NULL without the transform FromField function
			{Name: "ttl", Type: proto.ColumnType_DOUBLE,Transform: transform.FromField("TTL"), Description: "Time to live in seconds of the record."},

			// Other columns
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the record was created."},
			{Name: "locked", Type: proto.ColumnType_BOOL, Description: "[DEPRECATED] True if the record is locked."},
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the record was last modified."},
			{Name: "priority", Type: proto.ColumnType_INT, Description: "Priority for this record, primarily used for MX records."},
			{Name: "proxiable", Type: proto.ColumnType_BOOL, Description: "True if the record is eligible for Cloudflare's origin protection."},
			{Name: "proxied", Type: proto.ColumnType_BOOL, Description: "True if the record has Cloudflare's origin protection."},

			// JSON columns
			{Name: "data", Type: proto.ColumnType_JSON, Description: "Map of attributes that constitute the record value. Primarily used for LOC and SRV record types."},
			{Name: "meta", Type: proto.ColumnType_JSON, Description: "Cloudflare metadata for this record."},
		}),
	}
}

type RecordInfo struct {
	ZoneID string
	dns.RecordResponse
}

func listDNSRecord(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	zoneDetails := h.Item.(zones.Zone)
	inputZoneId := d.EqualsQualString("zone_id")

	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_dns_record.listDNSRecord", "connection error", err)
		return nil, err
	}

	// Only list zones stated in the input query
	if inputZoneId != "" && inputZoneId != zoneDetails.ID {
		return nil, nil
	}

	maxLimit := int32(500)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := dns.RecordListParams{
		ZoneID:  cloudflare.F(zoneDetails.ID),
		PerPage: cloudflare.F(float64(maxLimit)),
	}

	iter := conn.DNS.Records.ListAutoPaging(ctx, input)
	for iter.Next() {
		current := iter.Current()

		record := RecordInfo{
			ZoneID:				zoneDetails.ID,
			RecordResponse:		current,
		}
		d.StreamListItem(ctx, record)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	if err := iter.Err(); err != nil {
		logger.Error("cloudflare_dns_record.listDNSRecord", "DNSRecords api error", err)
		return nil, err
	}

	return nil, nil
}

func getDNSRecord(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_dns_record.getDNSRecord", "connection error", err)
		return nil, err
	}
	quals := d.EqualsQuals
	zoneID := quals["zone_id"].GetStringValue()
	id := quals["id"].GetStringValue()

	input := dns.RecordGetParams{
		ZoneID: cloudflare.F(zoneID),
	}

	item, err := conn.DNS.Records.Get(ctx, id, input)
	if err != nil {
		logger.Error("cloudflare_dns_record.getDNSRecord", "DNSRecord api error", err)
		return nil, err
	}

	record := RecordInfo{
		ZoneID:				zoneID,
		RecordResponse:		*item,
	}

	return record, nil
}
