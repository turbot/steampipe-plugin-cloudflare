package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go"

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
			ParentHydrate: listAccessApplications,
		},
		Columns: []*plugin.Column{
			//Top fields
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Access plolicy unique API identifier."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the policy. Only used in the UI."},
			{Name: "application_id", Type: proto.ColumnType_STRING, Hydrate: getParentApplicationDetails, Transform: transform.FromField("ID"), Description: "The id of application to which policy belongs."},
			{Name: "application_name", Type: proto.ColumnType_STRING, Hydrate: getParentApplicationDetails, Transform: transform.FromField("Name"), Description: "The name of application to which policy belongs."},
			{Name: "decision", Type: proto.ColumnType_STRING, Description: "Defines the action Access will take if the policy matches the user."},
			{Name: "precedence", Type: proto.ColumnType_INT, Description: "The unique precedence for policies on a single application."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when access policy was created."},

			// Other fields
			{Name: "purpose_justification_prompt", Type: proto.ColumnType_STRING, Description: "The text the user will be prompted with when a purpose justification is required."},
			{Name: "purpose_justification_required", Type: proto.ColumnType_BOOL, Description: "Defines whether or not the user is prompted for a justification when this policy is applied."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when access policy was last modified."},

			// Json fields
			{Name: "approval_groups", Type: proto.ColumnType_JSON, Description: "The list of approval groups that must approve the access request."},
			{Name: "include", Type: proto.ColumnType_JSON, Description: "The include policy works like an OR logical operator. The user must satisfy one of the rules in includes."},
			{Name: "exclude", Type: proto.ColumnType_JSON, Description: "The exclude policy works like a NOT logical operator. The user must not satisfy all of the rules in exclude."},
			{Name: "require", Type: proto.ColumnType_JSON, Description: "The require policy works like a AND logical operator. The user must satisfy all of the rules in require."},
		},
	}
}

func listAccessPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	appID := h.Item.(cloudflare.AccessApplication).ID
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
		items, result_info, err := conn.AccessPolicies(ctx, conn.AccountID, appID, opts)
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

func getAccessPolicy(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	appID := h.Item.(cloudflare.AccessApplication).ID
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
		items, result_info, err := conn.AccessPolicies(ctx, conn.AccountID, appID, opts)
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

func getParentApplicationDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return h.ParentItem.(cloudflare.AccessApplication), nil
}
