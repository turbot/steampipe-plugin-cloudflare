package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableCloudflareAccessGroup(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_access_group",
		Description: "Access Groups allows to define a set of users to which an application policy can be applied.",
		List: &plugin.ListConfig{
			ParentHydrate: listAccount,
			Hydrate:       listAccessGroups,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "account_id", Require: plugin.Optional},
				{Name: "account_name", Require: plugin.Optional},
			},
		},
		// Get Config - Currently SDK is not supporting get call
		Columns: []*plugin.Column{
			// Top fields
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Identifier of the access group."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Friendly name of the access group."},
			{Name: "account_id", Type: proto.ColumnType_STRING, Hydrate: getAccountDetails, Transform: transform.FromField("ID"), Description: "ID of the account, access group belongs."},
			{Name: "account_name", Type: proto.ColumnType_STRING, Hydrate: getAccountDetails, Transform: transform.FromField("Name"), Description: "Name of the account, access group belongs."},

			// Other fields
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when access group was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "TImestamp when access group was last modified."},

			// JSON fields
			{Name: "exclude", Type: proto.ColumnType_JSON, Description: "The exclude policy works like a NOT logical operator. The user must not satisfy all of the rules in exclude."},
			{Name: "include", Type: proto.ColumnType_JSON, Description: "The include policy works like an OR logical operator. The user must satisfy one of the rules in includes."},
			{Name: "require", Type: proto.ColumnType_JSON, Description: "The require policy works like a AND logical operator. The user must satisfy all of the rules in require."},
		},
	}
}

func listAccessGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	account := h.Item.(cloudflare.Account)

	if account_id := d.KeyColumnQualString("account_id"); account_id != "" && account.ID != account_id {
		return nil, nil
	}

	if account_name := d.KeyColumnQualString("account_name"); account_name != "" && account.Name != account_name {
		return nil, nil
	}

	conn, err := connect(ctx, d)
	if err != nil {
		logger.Error("listAccessGroups", "connection error", err)
		return nil, err
	}

	opts := cloudflare.PaginationOptions{
		PerPage: 100,
		Page:    1,
	}

	for {
		items, result_info, err := conn.AccessGroups(ctx, account.ID, opts)
		if err != nil {
			logger.Error("listAccessGroups", "AccessGroups api error", err)
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
