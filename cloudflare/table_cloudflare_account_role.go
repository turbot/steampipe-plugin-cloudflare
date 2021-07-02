package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

type accountRoleInfo = struct {
	ID          string                                      `json:"id"`
	Name        string                                      `json:"name"`
	Description string                                      `json:"description"`
	Permissions map[string]cloudflare.AccountRolePermission `json:"permissions"`
	AccountID   string
}

//// TABLE DEFINITION

func tableCloudflareAccountRole(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_account_role",
		Description: "Cloudflare Account Role",
		List: &plugin.ListConfig{
			Hydrate:       listRoles,
			ParentHydrate: listAccount,
		},
		Get: &plugin.GetConfig{
			Hydrate:           getAccountRole,
			KeyColumns:        plugin.AllColumns([]string{"account_id", "id"}),
			ShouldIgnoreError: isNotFoundError([]string{"HTTP status 403"}),
		},
		Columns: []*plugin.Column{
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

			// steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},
		},
	}
}

//// LIST FUNCTIONS

func listRoles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	accountData := h.Item.(cloudflare.Account)

	resp, err := conn.AccountRoles(ctx, accountData.ID)
	if err != nil {
		return nil, err
	}

	for _, role := range resp {
		d.StreamLeafListItem(ctx, accountRoleInfo{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Permissions: role.Permissions,
			AccountID:   accountData.ID,
		})
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAccountRole(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	accountID := d.KeyColumnQuals["account_id"].GetStringValue()
	id := d.KeyColumnQuals["id"].GetStringValue()

	data, err := conn.AccountRole(ctx, accountID, id)
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
