package cloudflare

import (
	"context"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type accountMemberInfo = struct {
	ID        string                              `json:"id"`
	Code      string                              `json:"code"`
	User      cloudflare.AccountMemberUserDetails `json:"user"`
	Status    string                              `json:"status"`
	Roles     []cloudflare.AccountRole            `json:"roles"`
	AccountID string
}

//// TABLE DEFINITION

func tableCloudflareAccountMember(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_account_member",
		Description: "Cloudflare Account Member",
		List: &plugin.ListConfig{
			Hydrate:       listAccountMembers,
			ParentHydrate: listAccount,
		},
		Get: &plugin.GetConfig{
			Hydrate:           getAccountMember,
			KeyColumns:        plugin.AllColumns([]string{"account_id", "id"}),
			ShouldIgnoreError: isNotFoundError([]string{"HTTP status 403", "HTTP status 404"}),
		},
		Columns: commonColumns([]*plugin.Column{
			{
				Name:        "user_email",
				Description: "Specifies the user email.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("User.Email"),
			},
			{
				Name:        "id",
				Description: "Specifies the account membership identifier.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "status",
				Description: "A member's status in the account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "account_id",
				Description: "Specifies the account id, the member is associated with.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("AccountID"),
			},
			{
				Name:        "code",
				Description: "The unique activation code for the account membership.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user",
				Description: "A set of information about the user.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "roles",
				Description: "A list of permissions that a Member of an Account has.",
				Type:        proto.ColumnType_JSON,
			},

			// steampipe standard columns
			{
				Name:        "title",
				Description: "Title of the resource.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.From(accountMemberTitle),
			},
		}),
	}
}

//// LIST FUNCTIONS

func listAccountMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	accountData := h.Item.(cloudflare.Account)

	var pages cloudflare.PaginationOptions
	for {
		response, pageData, err := conn.AccountMembers(ctx, accountData.ID, pages)
		if err != nil {
			return nil, err
		}

		for _, i := range response {
			d.StreamListItem(ctx, accountMemberInfo{i.ID, i.Code, i.User, i.Status, i.Roles, accountData.ID})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		if pageData.Page == pageData.TotalPages {
			break
		}
		pages.Page = pageData.Page + 1
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAccountMember(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	accountID := d.EqualsQuals["account_id"].GetStringValue()
	id := d.EqualsQuals["id"].GetStringValue()

	// empty check
	if accountID == "" || id == "" {
		return nil, nil
	}

	data, err := conn.AccountMember(ctx, accountID, id)
	if err != nil {
		return nil, err
	}

	return accountMemberInfo{data.ID, data.Code, data.User, data.Status, data.Roles, accountID}, nil
}

//// TRANSFORM FUNCTIONS

func accountMemberTitle(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	data := d.HydrateItem.(accountMemberInfo)

	if len(data.User.FirstName) > 0 && len(data.User.LastName) > 0 {
		return data.User.FirstName + " " + data.User.LastName, nil
	}
	return strings.Split(data.User.Email, "@")[0], nil
}
