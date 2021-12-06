package cloudflare

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudflare/cloudflare-go"

	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableCloudflareAccessPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_access_policy",
		Description: "Access Policies define the users or groups who can, or cannot, reach the Application Resource.",
		List: &plugin.ListConfig{
			Hydrate:       listAccessPolicies,
			ParentHydrate: listParentAccessApplications,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "application_id", Require: plugin.Optional},
			},
		},
		GetMatrixItem: BuildAccountmatrix,
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Access plolicy unique API identifier."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the policy. Only used in the UI."},
			{Name: "application_id", Type: proto.ColumnType_STRING, Hydrate: getParentApplicationDetails, Transform: transform.FromField("ID"), Description: "The id of application to which policy belongs."},
			{Name: "application_name", Type: proto.ColumnType_STRING, Hydrate: getParentApplicationDetails, Transform: transform.FromField("Name"), Description: "The name of application to which policy belongs."},
			{Name: "account_id", Type: proto.ColumnType_STRING, Transform: transform.FromMatrixItem(matrixKeyAccount), Description: "The ID of account where application belongs."},

			// Other columns
			{Name: "decision", Type: proto.ColumnType_STRING, Description: "Defines the action Access will take if the policy matches the user. Allowed values: allow, deny, non_identity, bypass"},
			{Name: "precedence", Type: proto.ColumnType_INT, Description: "The unique precedence for policies on a single application."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when access policy was created."},

			// JSON columns
			{Name: "exclude", Type: proto.ColumnType_JSON, Description: "The exclude policy works like a NOT logical operator. The user must not satisfy all of the rules in exclude."},
			{Name: "include", Type: proto.ColumnType_JSON, Description: "The include policy works like an OR logical operator. The user must satisfy one of the rules in includes."},
			{Name: "require", Type: proto.ColumnType_JSON, Description: "The require policy works like a AND logical operator. The user must satisfy all of the rules in require."},
		},
	}
}

func listAccessPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	accountID := d.KeyColumnQualString(matrixKeyAccount)
	appID := h.Item.(cloudflare.AccessApplication).ID
	inputAppID := d.KeyColumnQuals["application_id"].GetStringValue()

	// Avoid getting access policies for other applications id
	// "application_id" mentioned in where clause
	if inputAppID != "" && appID != inputAppID {
		return nil, nil
	}

	conn, err := connect(ctx, d)
	if err != nil {
		logger.Error("listAccessPolicies", "connection error", err)
		return nil, err
	}

	opts := cloudflare.PaginationOptions{
		PerPage: 100,
		Page:    1,
	}

	for {
		items, result_info, err := conn.AccessPolicies(ctx, accountID, appID, opts)
		if err != nil {
			logger.Error("listAccessPolicies", "AccessPolicies api error", err)
			return nil, err
		}
		for _, i := range items {
			d.StreamListItem(ctx, i)
		}

		if result_info.Page >= result_info.TotalPages {
			break
		}
		opts.Page = opts.Page + 1
	}

	return nil, nil
}

func listParentAccessApplications(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	accountID := d.KeyColumnQualString(matrixKeyAccount)

	conn, err := connect(ctx, d)
	if err != nil {
		logger.Error("listAccessApplications", "connection error", err)
		return nil, err
	}

	opts := cloudflare.PaginationOptions{
		PerPage: 100,
		Page:    1,
	}

	for {
		items, result_info, err := conn.AccessApplications(ctx, accountID, opts)
		if err != nil {
			var cloudFlareErr *cloudflare.APIRequestError
			if errors.As(err, &cloudFlareErr) {
				if helpers.StringSliceContains(cloudFlareErr.ErrorMessages(), "Access is not enabled. Visit the Access dashboard at https://dash.cloudflare.com/ and click the 'Enable Access' button.") {
					logger.Warn("listAccessApplications", fmt.Sprintf("AccessApplications api error for account: %s", accountID), err)
					return nil, nil
				}
			}
			logger.Error("listAccessApplications", "AccessApplications api error", err)
			return nil, err
		}
		for _, i := range items {
			d.StreamListItem(ctx, i)
		}

		if result_info.Page >= result_info.TotalPages {
			break
		}
		opts.Page = opts.Page + 1
	}

	return nil, nil
}

func getParentApplicationDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return h.ParentItem.(cloudflare.AccessApplication), nil
}
