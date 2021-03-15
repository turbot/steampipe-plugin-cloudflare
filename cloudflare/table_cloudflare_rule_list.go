package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableCloudflareRuleList(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_rule_list",
		Description: "Lists of IPs and CIDRs to enable more powerful expressions in filters.",
		List: &plugin.ListConfig{
			ParentHydrate: listAccount,
			Hydrate:       listRuleList,
		},
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("id"),
			ShouldIgnoreError: isNotFoundError([]string{"HTTP status 404"}),
			Hydrate:           getRuleList,
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "account_id", Type: proto.ColumnType_STRING, Hydrate: getParentAccount, Transform: transform.FromField("ID"), Description: "Account ID containing the Rule List."},
			{Name: "id", Type: proto.ColumnType_STRING, Description: "ID of the Rule List."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the Rule List, used in filter expressions."},
			// Other columns
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the Rule List was created."},
			{Name: "description", Type: proto.ColumnType_STRING, Description: "Description of the Rule List."},
			{Name: "items", Type: proto.ColumnType_JSON, Hydrate: getRuleListItems, Transform: transform.FromValue(), Description: "IPs and CIDRs in the list."},
			{Name: "kind", Type: proto.ColumnType_STRING, Description: "Kind of the Rule List (e.g. ip)."},
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the Rule List was last modified."},
			{Name: "num_items", Type: proto.ColumnType_INT, Description: "Number of items in the Rule List."},
			{Name: "num_referencing_filters", Type: proto.ColumnType_INT, Description: "Number of referencing filters to the Rule List."},
		},
	}
}

func listRuleList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	account := h.Item.(cloudflare.Account)
	conn.AccountID = account.ID
	items, err := conn.ListIPLists(ctx)
	if err != nil {
		return nil, err
	}
	for _, i := range items {
		d.StreamListItem(ctx, i)
	}
	return nil, nil
}

func getRuleList(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	account := h.Item.(cloudflare.Account)
	conn.AccountID = account.ID
	quals := d.KeyColumnQuals
	id := quals["id"].GetStringValue()
	item, err := conn.GetIPList(ctx, id)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func getParentAccount(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	account := h.ParentItem.(cloudflare.Account)
	return account, nil
}

func getRuleListItems(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}
	account := h.ParentItem.(cloudflare.Account)
	conn.AccountID = account.ID
	ruleList := h.Item.(cloudflare.IPList)
	items, err := conn.ListIPListItems(ctx, ruleList.ID)
	if err != nil {
		return nil, err
	}
	return items, nil
}
