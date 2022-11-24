package openstack

import (
	"context"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableOpenStackPort(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openstack_port",
		Description: "OpenStack Network Port",
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique id of the port.",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Human-readable name for the port. Might not be unique.",
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The description of the project (or tenant)",
				Transform:   transform.FromField("Description"),
			},
			{
				Name:        "network_id",
				Type:        proto.ColumnType_STRING,
				Description: "Network that this port is associated with.",
				Transform:   transform.FromField("NetworkID"),
			},
			{
				Name:        "admin_state_up",
				Type:        proto.ColumnType_BOOL,
				Description: "Administrative state of port. If false (down), port does not forward packets.",
				Transform:   transform.FromField("AdminStateUp"),
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "Indicates whether network is currently operational. Possible values include `ACTIVE', `DOWN', `BUILD', or `ERROR'. Plug-ins might define additional values.",
				Transform:   transform.FromField("Status"),
			},
			{
				Name:        "mac_address",
				Type:        proto.ColumnType_STRING,
				Description: "The MAC address associated with this port.",
				Transform:   transform.FromField("MACAddress"),
			},
			{
				Name:        "project_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the project owning this port.",
				Transform:   transform.FromField("ProjectID"),
			},
			{
				Name:        "device_owner",
				Type:        proto.ColumnType_STRING,
				Description: "Identifies the entity (e.g.: dhcp agent) using this port.",
				Transform:   transform.FromField("DeviceOwner"),
			},
			{
				Name:        "device_id",
				Type:        proto.ColumnType_STRING,
				Description: "Identifies the device (e.g., virtual server) using this port.",
				Transform:   transform.FromField("DeviceID"),
			},
			{
				Name:        "revision_number",
				Type:        proto.ColumnType_INT,
				Description: "RevisionNumber optionally set via extensions/standard-attr-revisions.",
				Transform:   transform.FromField("RevisionNumber"),
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_STRING,
				Description: "Timestamp when the port was created.",
				Transform:   TransformFromTimeField("CreatedAt"),
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_STRING,
				Description: "Timestamp when the port was last updated.",
				Transform:   TransformFromTimeField("UpdatedAt"),
			},
			{
				Name:        "security_group_ids",
				Type:        proto.ColumnType_STRING,
				Description: "The IDs of the security groups that apply to the current port.",
				// 	Hydrate:     getEc2InstanceARN,
				Transform: transform.FromField("SecurityGroups"), //.Transform(transform.EnsureStringArray),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listOpenStackPort,
			KeyColumns: plugin.KeyColumnSlice{
				&plugin.KeyColumn{
					Name:    "id",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "name",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "status",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "description",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "admin_state_up",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "network_id",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "project_id",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "device_owner",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "device_id",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "mac_address",
					Require: plugin.Optional,
				},
				// TODO: add tags support
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getOpenStackPort,
		},
	}
}

//// LIST FUNCTION

func listOpenStackPort(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	plugin.Logger(ctx).Debug("retrieving openstack ports list", "query data", toPrettyJSON(d))

	client, err := getServiceClient(ctx, d, NetworkV2)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	opts := buildOpenStackPortFilter(ctx, d.KeyColumnQuals)

	allPages, err := ports.List(client, opts).AllPages()
	if err != nil {
		plugin.Logger(ctx).Error("error listing ports with options", "options", toPrettyJSON(opts), "error", err)
		return nil, err
	}
	allPorts, err := ports.ExtractPorts(allPages)
	if err != nil {
		plugin.Logger(ctx).Error("error extracting ports", "error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("ports retrieved", "count", len(allPorts))

	for _, port := range allPorts {
		//plugin.Logger(ctx).Error("INFO", "port_id", port.ID, "security_group_ids", toJSON(opts), "error", err)
		port := port
		d.StreamListItem(ctx, &port)
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOpenStackPort(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	id := d.KeyColumnQuals["id"].GetStringValue()
	plugin.Logger(ctx).Debug("retrieving openstack port", "id", id)

	client, err := getServiceClient(ctx, d, NetworkV2)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	result := ports.Get(client, id)
	var port *ports.Port
	port, err = result.Extract()
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving port", "error", err)
		return nil, err
	}

	return port, nil
}

func buildOpenStackPortFilter(ctx context.Context, quals plugin.KeyColumnEqualsQualMap) ports.ListOpts {
	opts := ports.ListOpts{}

	if value, ok := quals["id"]; ok {
		opts.ID = value.GetStringValue()
	}
	if value, ok := quals["name"]; ok {
		opts.Name = value.GetStringValue()
	}
	if value, ok := quals["status"]; ok {
		opts.Status = value.GetStringValue()
	}
	if value, ok := quals["description"]; ok {
		opts.Description = value.GetStringValue()
	}
	if value, ok := quals["admin_state_up"]; ok {
		opts.AdminStateUp = pointerTo(value.GetBoolValue())
	}
	if value, ok := quals["network_id"]; ok {
		opts.NetworkID = value.GetStringValue()
	}
	if value, ok := quals["project_id"]; ok {
		opts.ProjectID = value.GetStringValue()
	}
	if value, ok := quals["device_owner"]; ok {
		opts.DeviceOwner = value.GetStringValue()
	}
	if value, ok := quals["device_id"]; ok {
		opts.DeviceID = value.GetStringValue()
	}
	if value, ok := quals["mac_address"]; ok {
		opts.MACAddress = value.GetStringValue()
	}
	plugin.Logger(ctx).Debug("returning", "filter", toPrettyJSON(opts))
	return opts
}
