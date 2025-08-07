package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/logpush"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableCloudflareLogpushJob(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_logpush_job",
		Description: "Cloudflare Logpush job is a configuration that automatically ships log data from a specific zone or account to a chosen external destination in near real‑time batch delivery",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "account_id", Require: plugin.AnyOf},
				{Name: "zone_id", Require: plugin.AnyOf},
			},
			Hydrate: listLogpushJobs,
		},
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "id", Require: plugin.Required},
				{Name: "account_id", Require: plugin.AnyOf},
				{Name: "zone_id", Require: plugin.AnyOf},
			},
			ShouldIgnoreError: isNotFoundError([]string{"Invalid logpush job identifier"}),
			Hydrate:           getLogpushJob,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_INT, Transform: transform.FromField("ID"), Description: "Logpush Job identifier."},
			{Name: "dataset", Type: proto.ColumnType_STRING, Description: "Name of the dataset."},
			{Name: "destination_conf", Type: proto.ColumnType_STRING, Description: "Uniquely identifies a resource (such as an s3 bucket) where data will be pushed."},
			{Name: "enabled", Type: proto.ColumnType_BOOL, Description: "Flag that indicates if the job is enabled."},
			{Name: "error_message", Type: proto.ColumnType_STRING, Description: "If not null, the job is currently failing."},
			{Name: "frequency", Type: proto.ColumnType_STRING, Description: "[Deprecated] The frequency at which Cloudflare sends batches of logs to your destination - Use `max_upload_*` parameters instead "},
			{Name: "kind", Type: proto.ColumnType_STRING, Description: "The kind parameter (optional) is used to differentiate between Logpush and Edge Log Delivery jobs (when supported by the dataset)."},
			{Name: "last_complete", Type: proto.ColumnType_TIMESTAMP, Description: "Records the last time for which logs have been successfully pushed."},
			{Name: "last_error", Type: proto.ColumnType_TIMESTAMP, Description: "Records the last time the job failed."},
			{Name: "logpull_options", Type: proto.ColumnType_STRING, Description: "[Deprecated] It specifies things like requested fields and timestamp formats - Use `output_options` instead."},
			{Name: "max_upload_bytes", Type: proto.ColumnType_DOUBLE, Description: "The maximum uncompressed file size of a batch of logs."},
			{Name: "max_upload_interval_seconds", Type: proto.ColumnType_DOUBLE, Description: "The maximum interval in seconds for log batches."},
			{Name: "max_upload_records", Type: proto.ColumnType_DOUBLE, Description: "The maximum number of log lines per batch."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Optional human readable job name."},

			// Query columns for filtering
			{Name: "account_id", Type: proto.ColumnType_STRING, Transform: transform.FromQual("account_id"), Description: "The account ID to filter logpush jobs."},
			{Name: "zone_id", Type: proto.ColumnType_STRING, Transform: transform.FromQual("zone_id"), Description: "The zone ID to filter logpush jobs."},
		
			// JSON Columns
			{Name: "output_options", Type: proto.ColumnType_JSON, Description: "The structured replacement for `logpull_options`."},
		}),
	}
}

//// LIST FUNCTION

// listLogpushJobs retrieves all logpush jobs for the specified account_id or zone_id.
//
// This function handles both account-level and zone-level logpush jobs:
// - Account-level logpush jobs (account_id)
// - Zone-level logpush jobs (zone_id)
func listLogpushJobs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_logpush_job.listLogpushJobs", "connection_error", err)
		return nil, err
	}

	// Get the qualifiers
	quals := d.EqualsQuals
	accountID := ""
	zoneID := ""

	if quals["account_id"] != nil {
		accountID = quals["account_id"].GetStringValue()
	}
	if quals["zone_id"] != nil {
		zoneID = quals["zone_id"].GetStringValue()
	}

	// Build API parameters based on account or zone context
	input := logpush.JobListParams{}
	if accountID != "" {
		input.AccountID = cloudflare.F(accountID)
	}
	if zoneID != "" {
		input.ZoneID = cloudflare.F(zoneID)
	}

	// Execute paginated API call
	iter := conn.Logpush.Jobs.ListAutoPaging(ctx, input)
	for iter.Next() {
		logpush_job := iter.Current()
		d.StreamListItem(ctx, logpush_job)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	if err := iter.Err(); err != nil {
		logger.Error("cloudflare_logpush_job.listLogpushJobs", "ListAutoPaging error", err)
		return nil, err
	}

	return nil, nil
}

//// GET FUNCTION

// getLogpushJob retrieves a specific logpush job by ID.
//
// This function fetches detailed information about a single logpush job.
// It requires either account_id or zone_id to determine the correct
// context for the logpush job lookup.
//
// Parameters:
// - id: The logpush job identifier (required)
// - account_id OR zone_id: The account or zone context (at least one required)
func getLogpushJob(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_logpush_job.getLogpushJob", "connection_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	logpushJobID := quals["id"].GetInt64Value()
	accountID := quals["account_id"].GetStringValue()
	zoneID := quals["zone_id"].GetStringValue()

	// Build API parameters with appropriate context
	input := logpush.JobGetParams{}
	if accountID != "" {
		input.AccountID = cloudflare.F(accountID)
	}
	if zoneID != "" {
		input.ZoneID = cloudflare.F(zoneID)
	}

	// Execute API call to get the specific logpush job
	ruleset, err := conn.Logpush.Jobs.Get(ctx, logpushJobID, input)
	if err != nil {
		logger.Error("cloudflare_logpush_job.getLogpushJob", "error", err)
		return nil, err
	}

	return ruleset, nil
}
