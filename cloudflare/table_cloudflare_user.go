package cloudflare

import (
	"context"

	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/user"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

func tableCloudflareUser(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_user",
		Description: "Information about the current user making the request.",
		List: &plugin.ListConfig{
			Hydrate: listUser,
		},
		Columns: commonColumns([]*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "ID of the user."},
			{Name: "email", Type: proto.ColumnType_STRING, Description: "Email of the user."},
			{Name: "username", Type: proto.ColumnType_STRING, Description: "Username (actually often in ID style) of the user."},

			// Other columns
			{Name: "telephone", Type: proto.ColumnType_STRING, Description: "Telephone number of the user."},
			{Name: "first_name", Type: proto.ColumnType_STRING, Description: "First name of the user."},
			{Name: "last_name", Type: proto.ColumnType_STRING, Description: "Last name of the user."},
			{Name: "country", Type: proto.ColumnType_STRING, Description: "Country of the user."},
			{Name: "zipcode", Type: proto.ColumnType_STRING, Description: "Zipcode of the user."},
			{Name: "created_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the user was created."},
			{Name: "modified_on", Type: proto.ColumnType_TIMESTAMP, Description: "When the user was last modified."},
			{Name: "api_key", Type: proto.ColumnType_STRING, Description: "API Key for the user."},
			{Name: "two_factor_authentication_enabled", Type: proto.ColumnType_BOOL, Description: "True if two factor authentication is enabled for this user."},

			// JSON columns
			{Name: "betas", Type: proto.ColumnType_JSON, Description: "Beta feature flags associated with the user."},
			{Name: "organizations", Type: proto.ColumnType_JSON, Description: "Organizations the user is a member of."},
		}),
	}
}

func listUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get the client's options from the connection
	opts := []option.RequestOption{}
	if conn != nil {
		opts = append(opts, conn.Options...)
	}

	userService := user.NewUserService(opts...)
	userDetails, err := userService.Get(ctx)
	if err != nil {
		plugin.Logger(ctx).Error("cloudflare_user.listUser", "api_error", err)
		return nil, err
	}

	d.StreamListItem(ctx, userDetails)
	return nil, nil
}
