package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/option"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

func tableCloudflareDNSRecord() *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_dns_record",
		Description: "DNS records for a zone.",
		List: &plugin.ListConfig{
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "zone_id", Require: plugin.Required, CacheMatch: query_cache.CacheMatchExact},
			},
			Hydrate: listDNSRecord,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.AllColumns([]string{"zone_id", "id"}),
			ShouldIgnoreError: isNotFoundError([]string{"HTTP status 404"}),
			Hydrate:           getDNSRecord,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "zone_id", Type: proto.ColumnType_STRING, Description: "Zone where the record is defined.", Transform: transform.FromQual("zone_id")},
			{Name: "zone_name", Type: proto.ColumnType_STRING, Description: "[Deprecated] Name of the zone where the record is defined."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "ID of the record."},
			{Name: "type", Type: proto.ColumnType_STRING, Description: "Type of the record (e.g. A, MX, CNAME)."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Domain name for the record (e.g. steampipe.io)."},
			{Name: "content", Type: proto.ColumnType_STRING, Description: "Content or value of the record. Changes by type, including IP address for A records and domain for CNAME records."},
			{Name: "ttl", Type: proto.ColumnType_INT, Description: "Time to live in seconds of the record."},

			// Other columns
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the record was created."},
			{Name: "locked", Type: proto.ColumnType_BOOL, Description: "True if the record is locked."},
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

func listDNSRecord(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get the client's options from the connection
	opts := []option.RequestOption{}
	if conn != nil {
		opts = append(opts, conn.Options...)
	}

	quals := d.EqualsQuals
	zoneID := quals["zone_id"].GetStringValue()

	// Empty check
	if zoneID == "" {
		return nil, nil
	}

	dnsService := dns.NewDNSService(opts...)
	recordService := dnsService.Records
	iter := recordService.ListAutoPaging(ctx, dns.RecordListParams{
		ZoneID: cloudflare.F(zoneID),
	})

	for iter.Next() {
		record := iter.Current()
		d.StreamListItem(ctx, record)
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

func getDNSRecord(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get the client's options from the connection
	opts := []option.RequestOption{}
	if conn != nil {
		opts = append(opts, conn.Options...)
	}

	quals := d.EqualsQuals
	zoneID := quals["zone_id"].GetStringValue()
	id := quals["id"].GetStringValue()

	dnsService := dns.NewDNSService(opts...)
	recordService := dnsService.Records
	record, err := recordService.Get(ctx, id, dns.RecordGetParams{
		ZoneID: cloudflare.F(zoneID),
	})
	if err != nil {
		return nil, err
	}
	return record, nil
}
