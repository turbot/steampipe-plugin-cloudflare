package cloudflare

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudflare/cloudflare-go"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableCloudflareAccessApplication(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_access_application",
		Description: "Access Applications are used to restrict access to a whole application using an authorisation gateway managed by Cloudflare.",
		List: &plugin.ListConfig{
			ParentHydrate: listAccount,
			Hydrate:       listAccessApplications,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "account_id", Require: plugin.Optional},
				{Name: "account_name", Require: plugin.Optional},
			},
		},
		// Get Config - Currently SDK is not supporting get call
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Application API uuid."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Friendly name of the access application."},
			{Name: "account_id", Type: proto.ColumnType_STRING, Hydrate: getAccountDetails, Transform: transform.FromField("ID"), Description: "ID of the account, access application belongs."},
			{Name: "account_name", Type: proto.ColumnType_STRING, Hydrate: getAccountDetails, Transform: transform.FromField("Name"), Description: "Name of the account, access application belongs."},
			{Name: "domain", Type: proto.ColumnType_STRING, Description: "The domain and path that access will block."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the application was created."},

			// Other columns
			{Name: "aud", Type: proto.ColumnType_STRING, Description: "Audience tag."},
			{Name: "auto_redirect_to_identity", Type: proto.ColumnType_BOOL, Description: "Option to skip identity provider selection if only one is configured in allowed_idps. Defaults to false (disabled)."},
			{Name: "custom_deny_message", Type: proto.ColumnType_STRING, Description: "Option that returns a custom error message when a user is denied access to the application."},
			{Name: "custom_deny_url", Type: proto.ColumnType_STRING, Description: "Option that redirects to a custom URL when a user is denied access to the application."},
			{Name: "enable_binding_cookie", Type: proto.ColumnType_BOOL, Description: "Option to provide increased security against compromised authorization tokens and CSRF attacks by requiring an additional \"binding\" cookie on requests. Defaults to false."},
			{Name: "session_duration", Type: proto.ColumnType_STRING, Description: "How often a user will be forced to re-authorise. Must be in the format \"48h\" or \"2h45m\". Valid time units are ns, us (or µs), ms, s, m, h. Defaults to 24h."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when the application was last modified."},

			// JSON columns
			{Name: "allowed_idps", Type: proto.ColumnType_JSON, Description: "The identity providers selected for the application."},
			{Name: "cors_headers", Type: proto.ColumnType_JSON, Description: "CORS configuration for the access application. See below for reference structure."},
		}),
	}
}

func listAccessApplications(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	account := h.Item.(cloudflare.Account)

	if accountID := d.EqualsQualString("account_id"); accountID != "" && account.ID != accountID {
		return nil, nil
	}

	if accountName := d.EqualsQualString("account_name"); accountName != "" && account.Name != accountName {
		return nil, nil
	}

	conn, err := connect(ctx, d)
	if err != nil {
		logger.Error("listAccessApplications", "connection error", err)
		return nil, err
	}

	opts := cloudflare.PaginationOptions{
		PerPage: 100,
		Page:    1,
	}

	type ListPageResponse struct {
		Applications []cloudflare.AccessApplication
		resp         cloudflare.ResultInfo
	}

	limit := d.QueryContext.Limit
	if limit != nil {
		if *limit < int64(opts.PerPage) {
			opts.PerPage = int(*limit)
		}
	}

	listPage := func(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
		applications, resp, err := conn.AccessApplications(ctx, account.ID, opts)
		return ListPageResponse{
			Applications: applications,
			resp:         resp,
		}, err
	}

	for {
		listPageResponse, err := plugin.RetryHydrate(ctx, d, h, listPage, &plugin.RetryConfig{ShouldRetryError: shouldRetryError})
		if err != nil {
			var cloudFlareErr *cloudflare.APIRequestError
			if errors.As(err, &cloudFlareErr) {
				if helpers.StringSliceContains(cloudFlareErr.ErrorMessages(), "Access is not enabled. Visit the Access dashboard at https://dash.cloudflare.com/ and click the 'Enable Access' button.") {
					logger.Warn("listAccessApplications", fmt.Sprintf("AccessApplications api error for account: %s", account.ID), err)
					return nil, nil
				}
			}
			logger.Error("listAccessApplications", "AccessApplications api error", err)
			return nil, err
		}

		listResponse := listPageResponse.(ListPageResponse)
		apps := listResponse.Applications
		resp := listResponse.resp
		for _, i := range apps {
			d.StreamListItem(ctx, i)
		}

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}

		if resp.Page >= resp.TotalPages {
			break
		}
		opts.Page = opts.Page + 1
	}

	return nil, nil
}

func getAccountDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	account := h.ParentItem.(cloudflare.Account)
	return account, nil
}
