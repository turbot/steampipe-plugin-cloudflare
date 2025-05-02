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

type accountMemberInfo struct {
	ID        string                `json:"id"`
	Status    shared.MemberStatus   `json:"status"`
	User      shared.MemberUser     `json:"user"`
	Roles     []shared.Role         `json:"roles"`
	Policies  []shared.MemberPolicy `json:"policies"`
	AccountID string
}

//// TABLE DEFINITION

func tableCloudflareAccountMember() *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_account_member",
		Description: "Account members are users that have been invited to access your account.",
		List: &plugin.ListConfig{
			ParentHydrate: listAccountForParent,
			Hydrate:       listAccountMembers,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"account_id", "id"}),
			Hydrate:    getAccountMember,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the member."},
			{Name: "account_id", Type: proto.ColumnType_STRING, Description: "The account ID associated with the member."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the member (accepted/pending)."},
			{Name: "user_email", Type: proto.ColumnType_STRING, Description: "The email address of the member.", Transform: transform.FromField("User.Email")},
			{Name: "user_first_name", Type: proto.ColumnType_STRING, Description: "The first name of the member.", Transform: transform.FromField("User.FirstName")},
			{Name: "user_last_name", Type: proto.ColumnType_STRING, Description: "The last name of the member.", Transform: transform.FromField("User.LastName")},
			{Name: "user_id", Type: proto.ColumnType_STRING, Description: "The unique identifier of the member's user.", Transform: transform.FromField("User.ID")},
			{Name: "roles", Type: proto.ColumnType_JSON, Description: "The roles assigned to the member."},
			{Name: "policies", Type: proto.ColumnType_JSON, Description: "The access policies assigned to the member."},
		},
	}
}

//// LIST FUNCTIONS

func listAccountMembers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accountData := h.Item.(accounts.Account)

	memberService := accounts.NewMemberService()
	result, err := memberService.List(ctx, accounts.MemberListParams{
		AccountID: cloudflare.F(accountData.ID),
	})
	if err != nil {
		plugin.Logger(ctx).Error("cloudflare_account_member.listAccountMembers", "list_error", err)
		return nil, err
	}

	for _, member := range result.Result {
		d.StreamListItem(ctx, accountMemberInfo{
			ID:        member.ID,
			Status:    member.Status,
			User:      member.User,
			Roles:     member.Roles,
			Policies:  member.Policies,
			AccountID: accountData.ID,
		})
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getAccountMember(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	accountID := d.EqualsQuals["account_id"].GetStringValue()
	id := d.EqualsQuals["id"].GetStringValue()

	memberService := accounts.NewMemberService()
	member, err := memberService.Get(ctx, id, accounts.MemberGetParams{
		AccountID: cloudflare.F(accountID),
	})
	if err != nil {
		plugin.Logger(ctx).Error("cloudflare_account_member.getAccountMember", "get_error", err)
		return nil, err
	}

	return accountMemberInfo{
		ID:        member.ID,
		Status:    member.Status,
		User:      member.User,
		Roles:     member.Roles,
		Policies:  member.Policies,
		AccountID: accountID,
	}, nil
}
