package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/rulesets"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableCloudflareRuleset(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_ruleset",
		Description: "Cloudflare Rulesets provide a powerful framework for configuring rules to process HTTP requests.",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "account_id", Require: plugin.AnyOf},
				{Name: "zone_id", Require: plugin.AnyOf},
			},
			Hydrate: listRulesets,
		},
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "id", Require: plugin.Required},
				{Name: "account_id", Require: plugin.AnyOf},
				{Name: "zone_id", Require: plugin.AnyOf},
			},
			ShouldIgnoreError: isNotFoundError([]string{"Invalid ruleset identifier"}),
			Hydrate:           getRuleset,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "Ruleset identifier."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The human-readable name of the ruleset."},
			{Name: "kind", Type: proto.ColumnType_STRING, Description: "The kind of the ruleset (managed, custom, root, or zone)."},
			{Name: "phase", Type: proto.ColumnType_STRING, Description: "The phase of the ruleset."},

			// Query columns for filtering
			{Name: "account_id", Type: proto.ColumnType_STRING, Transform: transform.FromQual("account_id"), Description: "The account ID to filter rulesets."},
			{Name: "zone_id", Type: proto.ColumnType_STRING, Transform: transform.FromQual("zone_id"), Description: "The zone ID to filter rulesets."},

			// Other columns
			{Name: "description", Type: proto.ColumnType_STRING, Description: "An informative description of the ruleset."},
			{Name: "last_updated", Type: proto.ColumnType_TIMESTAMP, Description: "The timestamp when the ruleset was last updated."},
			{Name: "version", Type: proto.ColumnType_STRING, Description: "The version of the ruleset."},

			// JSON columns
			{Name: "rules", Type: proto.ColumnType_JSON, Description: "The list of rules in the ruleset."},
		}),
	}
}

//// LIST FUNCTION

// listRulesets retrieves all rulesets for the specified account_id or zone_id.
//
// This function handles both account-level and zone-level rulesets:
// - Account-level rulesets (account_id)
// - Zone-level rulesets (zone_id)
func listRulesets(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("listRulesets", "connection_error", err)
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
	input := rulesets.RulesetListParams{}
	if accountID != "" {
		input.AccountID = cloudflare.F(accountID)
	}
	if zoneID != "" {
		input.ZoneID = cloudflare.F(zoneID)
	}

	// Execute paginated API call
	iter := conn.Rulesets.ListAutoPaging(ctx, input)
	for iter.Next() {
		ruleset := iter.Current()
		d.StreamListItem(ctx, ruleset)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	if err := iter.Err(); err != nil {
		logger.Error("listRulesets", "ListAutoPaging error", err)
		return nil, err
	}

	return nil, nil
}

//// GET FUNCTION

// getRuleset retrieves a specific ruleset by ID.
//
// This function fetches detailed information about a single ruleset, including all its rules
// and configuration details. It requires either account_id or zone_id to determine the correct
// context for the ruleset lookup.
//
// Parameters:
// - id: The ruleset identifier (required)
// - account_id OR zone_id: The account or zone context (at least one required)
func getRuleset(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("getRuleset", "connection_error", err)
		return nil, err
	}

	quals := d.EqualsQuals
	rulesetID := quals["id"].GetStringValue()
	accountID := quals["account_id"].GetStringValue()
	zoneID := quals["zone_id"].GetStringValue()

	// Validate required parameters
	if rulesetID == "" {
		return nil, nil
	}

	// Build API parameters with appropriate context
	input := rulesets.RulesetGetParams{}
	if accountID != "" {
		input.AccountID = cloudflare.F(accountID)
	}
	if zoneID != "" {
		input.ZoneID = cloudflare.F(zoneID)
	}

	// Execute API call to get the specific ruleset
	ruleset, err := conn.Rulesets.Get(ctx, rulesetID, input)
	if err != nil {
		logger.Error("getRuleset", "error", err)
		return nil, err
	}

	return &ruleset, nil
}
