package cloudflare

import (
	"context"
	"time"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/custom_pages"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

// CUSTOM PAGE STRUCT DEFINITION

type CustomPage struct {
	ID            string    `json:"id"`
	Description   string    `json:"description"`
	State         string    `json:"state"`
	URL           string    `json:"url"`
	ModifiedOn    time.Time `json:"modified_on"`
	CreatedOn     time.Time `json:"created_on"`
	RequiredTokens []string `json:"required_tokens"`
	PreviewTarget string    `json:"preview_target"`
}

//// TABLE DEFINITION

func tableCloudflareCustomPage(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_custom_page",
		Description: "",
		List: &plugin.ListConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "account_id", Require: plugin.AnyOf},
				{Name: "zone_id", Require: plugin.AnyOf},
			},
			Hydrate: listCustomPages,
		},
		Get: &plugin.GetConfig{
			KeyColumns: []*plugin.KeyColumn{
				{Name: "id", Require: plugin.Required},
				{Name: "account_id", Require: plugin.AnyOf},
				{Name: "zone_id", Require: plugin.AnyOf},
			},
			ShouldIgnoreError: isNotFoundError([]string{"Invalid custom page identifier"}),
			Hydrate:           getCustomPage,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("id"), Description: "Custom page identifier."},
			{Name: "description", Type: proto.ColumnType_STRING, Transform: transform.FromField("description"), Description: "Custom page description."},
			{Name: "state", Type: proto.ColumnType_STRING, Transform: transform.FromField("state"), Description: "The custom page state."},
			{Name: "url", Type: proto.ColumnType_STRING, Transform: transform.FromField("url"), Description: "The URL associated with the custom page."},
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("modified_on"), Description: "When the setting was last modified."},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("created_on"), Description: "When the custom page was created."},
			{Name: "preview_target", Type: proto.ColumnType_STRING, Transform: transform.FromField("preview_target"), Description: "Preview action to apply."},

			// Query columns for filtering
			{Name: "account_id", Type: proto.ColumnType_STRING, Transform: transform.FromQual("account_id"), Description: "The account ID to filter custom pages."},
			{Name: "zone_id", Type: proto.ColumnType_STRING, Transform: transform.FromQual("zone_id"), Description: "The zone ID to filter custom pages."},

			// JSON Columns
			{Name: "required_tokens", Type: proto.ColumnType_JSON, Transform: transform.FromField("required_tokens"), Description: "Error tokens are required by the custom page."},
		}),
	}
}

//// LIST FUNCTION

// listCustomPages retrieves all custom pages for the specified account_id or zone_id.
//
// This function handles both account-level and zone-level custom pages:
// - Account-level custom pages (account_id)
// - Zone-level custom pages (zone_id)
func listCustomPages(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_custom_pages.listCustomPages", "connection_error", err)
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
	input := custom_pages.CustomPageListParams{}
	if accountID != "" {
		input.AccountID = cloudflare.F(accountID)
	}
	if zoneID != "" {
		input.ZoneID = cloudflare.F(zoneID)
	}

	// Execute paginated API call
	iter := conn.CustomPages.ListAutoPaging(ctx, input)
	for iter.Next() {
		current := iter.Current()

		d.StreamListItem(ctx, current)

		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	if err := iter.Err(); err != nil {
		logger.Error("cloudflare_custom_pages.listCustomPages", "ListAutoPaging error", err)
		return nil, err
	}

	return nil, nil
}

// GET FUNCTION

// getCustomPage retrieves a specific custom page by ID.
// Parameters:
// - id: The custom page identifier (required)
// - account_id OR zone_id: The account or zone context (at least one required)

func getCustomPage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_custom_pages.getCustomPage", "connection_error", err)
		return nil, err
	}
	quals := d.EqualsQuals
	customPageID := quals["id"].GetStringValue()
	accountID := quals["account_id"].GetStringValue()
	zoneID := quals["zone_id"].GetStringValue()

	// Build API parameters with appropriate context
	input := custom_pages.CustomPageGetParams{}
	if accountID != "" {
		input.AccountID = cloudflare.F(accountID)
	}
	if zoneID != "" {
		input.ZoneID = cloudflare.F(zoneID)
	}

	// Execute API call to get the specific custom page
	customPage, err := conn.CustomPages.Get(ctx, customPageID, input)
	if err != nil {
		logger.Error("cloudflare_custom_pages.getCustomPages", "error", err)
		return nil, err
	}

	return *customPage, nil
}
