package cloudflare

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/accounts"
	"github.com/cloudflare/cloudflare-go/v4/zero_trust"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Transform: transform.FromField("ID"), Description: "Identifier of the access group."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Friendly name of the access group."},
			{Name: "account_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("Account.ID"), Description: "ID of the account, access group belongs."},
			{Name: "account_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("Account.Name"), Description: "Name of the account, access group belongs."},

			// Other columns
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "Timestamp when access group was created."},
			{Name: "updated_at", Type: proto.ColumnType_TIMESTAMP, Description: "TImestamp when access group was last modified."},

			// JSON columns
			{Name: "exclude", Type: proto.ColumnType_JSON, Description: "The exclude policy works like a NOT logical operator. The user must not satisfy all of the rules in exclude."},
			{Name: "include", Type: proto.ColumnType_JSON, Description: "The include policy works like an OR logical operator. The user must satisfy one of the rules in includes."},
			{Name: "require", Type: proto.ColumnType_JSON, Description: "The require policy works like a AND logical operator. The user must satisfy all of the rules in require."},
		}),
	}
}
type AccessGroupInfo struct {
	Account accounts.Account
	zero_trust.AccessGroupListResponse
}

func listAccessGroups(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	account := h.Item.(accounts.Account)

	if accountID := d.EqualsQualString("account_id"); accountID != "" && account.ID != accountID {
		return nil, nil
	}

	if accountName := d.EqualsQualString("account_name"); accountName != "" && account.Name != accountName {
		return nil, nil
	}

	conn, err := connectV4(ctx, d)
	if err != nil {
		logger.Error("cloudflare_access_group.listAccessGroups", "connection error", err)
		return nil, err
	}

	opts := zero_trust.AccessGroupListParams{
		AccountID: cloudflare.String(account.ID),
	}

	iter := conn.ZeroTrust.Access.Groups.ListAutoPaging(ctx, opts)

	if err := iter.Err(); err != nil {
		var apiErr *cloudflare.Error
		if errors.As(err, &apiErr) {
			if strings.Contains(apiErr.Error(), "Access is not enabled. Visit the Access dashboard at https://dash.cloudflare.com/ and click the 'Enable Access' button.") {
				logger.Warn("listAccessGroups", fmt.Sprintf("AccessGroups api error for account: %s", account.ID), err)
				return nil, nil
			}
		}
		logger.Error("cloudflare_access_group.listAccessGroups", "AccessGroups api error", err)
		return nil, err
	}

	for iter.Next() {
		group := iter.Current()
		d.StreamListItem(ctx, AccessGroupInfo{
			account,
			group,
		})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}
