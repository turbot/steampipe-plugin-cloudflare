package cloudflare

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/shared"
	"github.com/cloudflare/cloudflare-go/v4/zero_trust"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
		GetMatrixItemFunc: BuildAccountmatrix,
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "Access policy unique API identifier."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the policy. Only used in the UI."},
			{Name: "application_id", Type: proto.ColumnType_STRING, Hydrate: getParentApplicationDetails, Transform: transform.FromField("ID"), Description: "The id of application to which policy belongs."},
			{Name: "application_name", Type: proto.ColumnType_STRING, Hydrate: getParentApplicationDetails, Transform: transform.FromField("Name"), Description: "The name of application to which policy belongs."},
			{Name: "account_id", Type: proto.ColumnType_STRING, Transform: transform.FromMatrixItem(matrixKeyAccount), Description: "The ID of account where application belongs."},

			// Other columns
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when access policy was created."},
			{Name: "decision", Type: proto.ColumnType_STRING, Description: "Defines the action Access will take if the policy matches the user. Allowed values: allow, deny, non_identity, bypass"},
			{Name: "precedence", Type: proto.ColumnType_INT, Description: "The unique precedence for policies on a single application."},
			{Name: "purpose_justification_prompt", Type: proto.ColumnType_STRING, Description: "The text the user will be prompted with when a purpose justification is required."},
			{Name: "purpose_justification_required", Type: proto.ColumnType_BOOL, Description: "Defines whether or not the user is prompted for a justification when this policy is applied."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when access policy was last modified."},

			// JSON columns
			{Name: "approval_groups", Type: proto.ColumnType_JSON, Description: "The list of approval groups that must approve the access request."},
			{Name: "exclude", Type: proto.ColumnType_JSON, Description: "The exclude policy works like a NOT logical operator. The user must not satisfy all of the rules in exclude."},
			{Name: "include", Type: proto.ColumnType_JSON, Description: "The include policy works like an OR logical operator. The user must satisfy one of the rules in includes."},
			{Name: "require", Type: proto.ColumnType_JSON, Description: "The require policy works like a AND logical operator. The user must satisfy all of the rules in require."},
		}),
	}
}

func listAccessPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	accountID := d.EqualsQualString(matrixKeyAccount)
	app := h.Item.(zero_trust.AccessApplicationListResponse)
	inputAppID := d.EqualsQuals["application_id"].GetStringValue()

	// Avoid getting access policies for other applications id
	// "application_id" mentioned in where clause
	if inputAppID != "" && app.ID != inputAppID {
		return nil, nil
	}

	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("listAccessPolicies", "connection error", err)
		return nil, err
	}

	opts := zero_trust.AccessApplicationPolicyListParams{
		AccountID: cloudflare.String(accountID),
	}

	iter := conn.ZeroTrust.Access.Applications.Policies.ListAutoPaging(ctx, app.ID, opts)

	if err := iter.Err(); err != nil {
		logger.Error("listAccessPolicies", "AccessPolicies api error", err)
		return nil, err
	}

	for iter.Next() {
		policy := iter.Current()
		d.StreamListItem(ctx, policy)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func listParentAccessApplications(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	accountID := d.EqualsQualString(matrixKeyAccount)

	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("listParentAccessApplications", "connection error", err)
		return nil, err
	}

	opts := zero_trust.AccessApplicationListParams{
		AccountID: cloudflare.String(accountID),
	}

	iter := conn.ZeroTrust.Access.Applications.ListAutoPaging(ctx, opts)

	if err := iter.Err(); err != nil {
		var cloudFlareErr *shared.ErrorData
		if errors.As(err, &cloudFlareErr) {
			if slices.Contains([]string{cloudFlareErr.Message}, "Access is not enabled. Visit the Access dashboard at https://dash.cloudflare.com/ and click the 'Enable Access' button.") {
				logger.Warn("listParentAccessApplications", fmt.Sprintf("AccessApplications api error for account: %s", accountID), err)
				return nil, nil
			}
		}
		logger.Error("listParentAccessApplications", "AccessApplications api error", err)
		return nil, err
	}

	for iter.Next() {
		application := iter.Current()
		d.StreamListItem(ctx, application)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

func getParentApplicationDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return h.ParentItem.(zero_trust.AccessApplicationListResponse), nil
}
