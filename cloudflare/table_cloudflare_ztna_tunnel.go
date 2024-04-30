package cloudflare

import (
	"context"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type TunnelAndConfig struct {
	ID             string                         `json:"id,omitempty"`
	Name           string                         `json:"name,omitempty"`
	Secret         string                         `json:"tunnel_secret,omitempty"`
	CreatedAt      *time.Time                     `json:"created_at,omitempty"`
	DeletedAt      *time.Time                     `json:"deleted_at,omitempty"`
	Connections    []cloudflare.TunnelConnection  `json:"connections,omitempty"`
	ConnsActiveAt  *time.Time                     `json:"conns_active_at,omitempty"`
	ConnInactiveAt *time.Time                     `json:"conns_inactive_at,omitempty"`
	TunnelType     string                         `json:"tun_type,omitempty"`
	Status         string                         `json:"status,omitempty"`
	RemoteConfig   bool                           `json:"remote_config,omitempty"`
	Configuration  cloudflare.TunnelConfiguration `json:"configuration"`
	Deleted        bool                           `json:"deleted,omitempty"`
	// ShowConfig     bool                           `json:"show_config,omitempty"`
}

func tableCloudflareZtnaTunnel(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "cloudflare_ztna_tunnel",
		Description: "List of ZTNA tunnels.",
		List: &plugin.ListConfig{
			Hydrate:       listTunnels,
			ParentHydrate: listAccount,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "deleted", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			// Top columns
			{Name: "id", Type: proto.ColumnType_STRING, Description: "ID of the tunnel."},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "Name of the tunnel."},
			{Name: "created_at", Type: proto.ColumnType_TIMESTAMP, Description: "When the tunnel was created."},
			{Name: "deleted_at", Type: proto.ColumnType_TIMESTAMP, Description: "When the tunnel was deleted."},
			{Name: "connections", Type: proto.ColumnType_JSON, Description: "Connections of the tunnel."},
			{Name: "conns_active_at", Type: proto.ColumnType_TIMESTAMP, Description: "When the tunnel was last active."},
			{Name: "conns_inactive_at", Type: proto.ColumnType_TIMESTAMP, Description: "When the tunnel was last inactive."},
			{Name: "tun_type", Type: proto.ColumnType_STRING, Description: "Type of the tunnel."},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "Status of the tunnel."},
			{Name: "remote_config", Type: proto.ColumnType_BOOL, Description: "Remote configuration of the tunnel."},
			{Name: "deleted", Type: proto.ColumnType_BOOL, Transform: transform.FromField("DeletedAt").Transform(deletedStatus), Description: "A boolean value indicating if the tunnel is deleted or not."},

			// TODO: Not sure if this should be part of this table, or a new table.
			// {Name: "configuration", Type: proto.ColumnType_JSON, Transform: transform.FromField("ShowConfig").Transform(getConfig), Description: "Configuration of the tunnel."},
			// {Name: "show_config", Type: proto.ColumnType_BOOL, Transform: transform.FromField("ShowConfig").Transform(showConfig), Description: "A boolean value that will hydrate the config of each tunnel."},
		},
	}
}

func listTunnels(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	account := h.Item.(cloudflare.Account)
	conn, err := connect(ctx, d)
	if err != nil {
		return nil, err
	}

	// The API will get both active and deleted tunnels by default. Defining `IsDeleted` value will return one or the other.
	// For now we're just going to retrieve all since it's a light call.
	// showDeleted := nill
	// if d.EqualsQuals["deleted"].GetBoolValue() == true {
	// 	showDeleted = true
	// }
	opts := cloudflare.TunnelListParams{
		// IsDeleted: &showDeleted,
	}

	resp, _, err := conn.ListTunnels(ctx, cloudflare.AccountIdentifier(account.ID), opts)
	if err != nil {
		return nil, err
	}

	// Not sure if this should be part of this table, or part of a new table
	showConfig := false
	// if d.EqualsQuals["show_config"].GetBoolValue() == true {
	// 	showConfig = true
	// }

	for _, t := range resp {

		configuration := cloudflare.TunnelConfigurationResult{}
		if showConfig == true {
			configuration, err = conn.GetTunnelConfiguration(ctx, cloudflare.AccountIdentifier(account.ID), t.ID)
			if err != nil {
				continue
			}
		}

		d.StreamListItem(ctx, TunnelAndConfig{
			ID:             t.ID,
			Name:           t.Name,
			CreatedAt:      t.CreatedAt,
			DeletedAt:      t.DeletedAt,
			Connections:    t.Connections,
			ConnsActiveAt:  t.ConnsActiveAt,
			ConnInactiveAt: t.ConnInactiveAt,
			TunnelType:     t.TunnelType,
			Status:         t.Status,
			RemoteConfig:   t.RemoteConfig,
			Configuration:  configuration.Config,
		})
	}

	return nil, nil
}

func deletedStatus(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	// If deleted_at is not nil, it has been deleted
	return d.Value.(*time.Time) != nil, nil
}
