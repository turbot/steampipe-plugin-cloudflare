package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/zero_trust"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableCloudflareAccessPolicy(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_access_policy",
		Description: "Access Policies define the users or groups who can, or cannot, reach the Application Resource.",
		List: &plugin.ListConfig{
			Hydrate:       listAccessPolicies,
			ParentHydrate: listParentAccessApplications,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Access policy identifier."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Access policy name."},
			{Name: "application_id", Type: proto.ColumnType_STRING, Description: "Access application identifier."},
			{Name: "decision", Type: proto.ColumnType_STRING, Description: "Access policy decision (allow, deny, non_identity, bypass)."},
			{Name: "precedence", Type: proto.ColumnType_INT, Description: "Access policy precedence."},

			// Other columns
			{Name: "approval_group", Type: proto.ColumnType_JSON, Description: "Access policy approval group."},
			{Name: "approval_required", Type: proto.ColumnType_BOOL, Description: "Access policy approval required."},
			{Name: "exclude", Type: proto.ColumnType_JSON, Description: "Access policy exclude rules."},
			{Name: "include", Type: proto.ColumnType_JSON, Description: "Access policy include rules."},
			{Name: "purpose_justification_prompt", Type: proto.ColumnType_STRING, Description: "Access policy purpose justification prompt."},
			{Name: "purpose_justification_required", Type: proto.ColumnType_BOOL, Description: "Access policy purpose justification required."},
			{Name: "require", Type: proto.ColumnType_JSON, Description: "Access policy require rules."},
		}),
	}
}

func listAccessPolicies(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	accountID := d.EqualsQualString(matrixKeyAccount)
	app := h.Item.(*zero_trust.AccessApplicationListResponse)
	appID := app.ID
	inputAppID := d.EqualsQuals["application_id"].GetStringValue()

	// If an application_id is specified in the where clause, skip the rows that don't match
	if inputAppID != "" && inputAppID != appID {
		return nil, nil
	}

	conn, err := connect(ctx, d)
	if err != nil {
		logger.Error("listAccessPolicies", "connection_error", err)
		return nil, err
	}

	// Get the client's options from the connection
	opts := []option.RequestOption{}
	if conn != nil {
		opts = append(opts, conn.Options...)
	}

	service := zero_trust.NewAccessService(opts...)
	iter := service.Applications.Policies.ListAutoPaging(ctx, appID, zero_trust.AccessApplicationPolicyListParams{
		AccountID: cloudflare.F(accountID),
	})

	for iter.Next() {
		policy := iter.Current()
		d.StreamListItem(ctx, policy)
	}
	if err := iter.Err(); err != nil {
		logger.Error("listAccessPolicies", "api error", err)
		return nil, err
	}

	return nil, nil
}

func listParentAccessApplications(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	conn, err := connect(ctx, d)
	if err != nil {
		logger.Error("listParentAccessApplications", "connection_error", err)
		return nil, err
	}

	// Get the client's options from the connection
	opts := []option.RequestOption{}
	if conn != nil {
		opts = append(opts, conn.Options...)
	}

	service := zero_trust.NewAccessService(opts...)
	iter := service.Applications.ListAutoPaging(ctx, zero_trust.AccessApplicationListParams{
		AccountID: cloudflare.F(d.EqualsQualString(matrixKeyAccount)),
	})

	for iter.Next() {
		app := iter.Current()
		d.StreamListItem(ctx, app)
	}
	if err := iter.Err(); err != nil {
		logger.Error("listParentAccessApplications", "api error", err)
		return nil, err
	}

	return nil, nil
}

func getParentApplicationDetails(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	return h.ParentItem.(*zero_trust.AccessApplicationListResponse), nil
}
