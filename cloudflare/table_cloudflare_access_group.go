package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func tableCloudflareAccessGroup(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_access_group",
		Description: "Access Groups allows to define a set of users to which an application policy can be applied.",
		List: &plugin.ListConfig{
			Hydrate: listAccessGroups,
		},
		// Get Config - Currently SDK is not supporting get call
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "Identifier of the access group."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Friendly name of the access group."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when access group was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "TImestamp when access group was last modified."},
			{Name: "exclude", Type: proto.ColumnType_JSON, Description: "The exclude policy works like a NOT logical operator. The user must not satisfy all of the rules in exclude."},
			{Name: "include", Type: proto.ColumnType_JSON, Description: "The include policy works like an OR logical operator. The user must satisfy one of the rules in includes."},
			{Name: "require", Type: proto.ColumnType_JSON, Description: "The require policy works like a AND logical operator. The user must satisfy all of the rules in require."},
		},
	}
}

func listAccessGroups(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	opts := cloudflare.PaginationOptions{
		PerPage: 100,
		Page:    1,
	}

	for {
		items, result_info, err := conn.AccessGroups(ctx, conn.AccountID, opts)
		if err != nil {
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
