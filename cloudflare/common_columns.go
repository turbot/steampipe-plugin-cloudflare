package cloudflare

import (
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func commonColumns(columns []*plugin.Column) []*plugin.Column {
	return append(columns, &plugin.Column{
		Name:        "user_id",
		Hydrate:     getUserId,
		Type:        proto.ColumnType_STRING,
		Description: "ID of the current user.",
		Transform:   transform.FromValue(),
	})
}
