package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/accounts"
	"github.com/cloudflare/cloudflare-go/v4/shared"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type accountRoleInfo = struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Permissions shared.RolePermissions `json:"permissions"`
	AccountID   string
}

//// TABLE DEFINITION

func tableCloudflareAccountRole(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_account_role",
		Description: "A Role defines what permissions a Member of an Account has.",
		List: &plugin.ListConfig{
			Hydrate:       listRoles,
			ParentHydrate: listAccount,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "account_id", Require: plugin.Optional},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate:           getAccountRole,
			KeyColumns:        plugin.AllColumns([]string{"account_id", "id"}),
			ShouldIgnoreError: isNotFoundError([]string{"HTTP status 403"}),
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{
				Name:        "id",
				Description: "Specifies the Role identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Description: "Specifies the name of the role.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "account_id",
				Description: "Specifies the account id where the role is created at.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountID"),
			},

			// Other columns
			{
				Name:        "description",
				Description: "A description of the role.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "permissions",
				Description: "A list of permissions attached with the role.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		}),
	}
}

//// LIST FUNCTIONS

func listRoles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	account := h.Item.(accounts.Account)
	if accountID := d.EqualsQualString("account_id"); accountID != "" && account.ID != accountID {
		return nil, nil
	}

	conn, err := connectV4(ctx, d)
	if err != nil {
		return nil, err
	}

	maxLimit := int32(500)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	iter := conn.Accounts.Roles.ListAutoPaging(ctx, accounts.RoleListParams{
		PerPage:   cloudflare.F(float64(maxLimit)),
		AccountID: cloudflare.F(account.ID),
	})
	if err := iter.Err(); err != nil {
		return nil, err
	}

	for iter.Next() {
		role := iter.Current()
		d.StreamLeafListItem(ctx, accountRoleInfo{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Permissions: role.Permissions,
			AccountID:   account.ID,
		})

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAccountRole(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connectV4(ctx, d)
	if err != nil {
		return nil, err
	}

	accountID := d.EqualsQuals["account_id"].GetStringValue()
	id := d.EqualsQuals["id"].GetStringValue()

	data, err := conn.Accounts.Roles.Get(ctx, id, accounts.RoleGetParams{
		AccountID: cloudflare.F(accountID),
	})
	if err != nil {
		return nil, err
	}
	return accountRoleInfo{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		Permissions: data.Permissions,
		AccountID:   accountID,
	}, nil
}
