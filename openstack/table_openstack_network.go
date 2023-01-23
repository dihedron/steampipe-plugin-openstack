package openstack

import (
	"context"

	"github.com/dihedron/steampipe-plugin-utils/utils"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableOpenStackNetwork(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openstack_network",
		Description: "OpenStack Network",
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique id of the network.",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Human-readable name for the network.",
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The description of the project (or tenant)",
				Transform:   transform.FromField("Description"),
			},
			{
				Name:        "project_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the project owning this network.",
				Transform:   transform.FromField("ProjectID"),
			},
			{
				Name:        "admin_state_up",
				Type:        proto.ColumnType_BOOL,
				Description: "The administrative state of the network, which is up (true) or down (false).",
				Transform:   transform.FromField("AdminStateUp"),
			},
			{
				Name:        "availability_zone_hints",
				Type:        proto.ColumnType_STRING,
				Description: "The availability zone candidate for the network.",
				Transform:   transform.FromField("AvailabilityZoneHints").Transform(transform.EnsureStringArray),
			},
			// {
			// 	Name:        "availability_zones",
			// 	Type:        proto.ColumnType_STRING,
			// 	Description: "The availability zone for the network.",
			// 	Transform:   transform.FromField("AvailabilityZones").Transform(transform.EnsureStringArray),
			// },
			// {
			// 	Name:        "dns_domain",
			// 	Type:        proto.ColumnType_STRING,
			// 	Description: "A valid DNS domain.",
			// 	Transform:   transform.FromField("DNSDomain"),
			// },
			// {
			// 	Name:        "ipv4_address_scope",
			// 	Type:        proto.ColumnType_STRING,
			// 	Description: "The ID of the IPv4 address scope that the network is associated with.",
			// 	Transform:   transform.FromField("IPv4AddressScope"),
			// },
			// {
			// 	Name:        "ipv6_address_scope",
			// 	Type:        proto.ColumnType_STRING,
			// 	Description: "The ID of the IPv6 address scope that the network is associated with.",
			// 	Transform:   transform.FromField("IPv6AddressScope"),
			// },
			// {
			// 	Name:        "l2_adjacency",
			// 	Type:        proto.ColumnType_BOOL,
			// 	Description: "Indicates whether L2 connectivity is available throughout the network.",
			// 	Transform:   transform.FromField("L2Adjacency"),
			// },
			// {
			// 	Name:        "mtu",
			// 	Type:        proto.ColumnType_INT,
			// 	Description: "The maximum transmission unit (MTU) value to address fragmentation. Minimum value is 68 for IPv4, and 1280 for IPv6.",
			// 	Transform:   transform.FromField("MTU"),
			// },
			// {
			// 	Name:        "port_security_enabled",
			// 	Type:        proto.ColumnType_BOOL,
			// 	Description: "The port default security status of the network. Valid values are enabled (true) and disabled (false).",
			// 	Transform:   transform.FromField("PortSecurityEnabled"),
			// },
			// {
			// 	Name:        "provider_network_type",
			// 	Type:        proto.ColumnType_STRING,
			// 	Description: "The type of physical network that this network is mapped to. For example, flat, vlan, vxlan, or gre.",
			// 	Transform:   transform.FromField("ProviderNetworkType"),
			// },
			// {
			// 	Name:        "provider_physical_network",
			// 	Type:        proto.ColumnType_STRING,
			// 	Description: "The physical network where this network/segment is implemented.",
			// 	Transform:   transform.FromField("ProviderPhysicalNetwork"),
			// },
			// {
			// 	Name:        "provider_segmentation_id",
			// 	Type:        proto.ColumnType_INT,
			// 	Description: "The ID of the isolated segment on the physical network. The network_type attribute defines the segmentation model.",
			// 	Transform:   transform.FromField("ProviderSegmentationID"),
			// },
			{
				Name:        "qos_policy_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the QoS policy associated with the network.",
				Transform:   transform.FromField("QoSPolicyID"),
			},
			{
				Name:        "revision_number",
				Type:        proto.ColumnType_INT,
				Description: "The revision number of the resource, optionally set via extensions/standard-attr-revisions.",
				Transform:   transform.FromField("RevisionNumber"),
			},
			// {
			// 	Name:        "router_external",
			// 	Type:        proto.ColumnType_BOOL,
			// 	Description: "Defines whether the network may be used for creation of floating IPs. Only networks with this flag may be an external gateway for routers.",
			// 	Transform:   transform.FromField("RouterExternal"),
			// },
			// {
			// 	Name:        "segments",
			// 	Type:        proto.ColumnType_STRING,
			// 	Description: "A list of provider segment objects.",
			// 	Transform:   transform.FromField("Status").Transform(transform.EnsureStringArray),
			// },
			{
				Name:        "shared",
				Type:        proto.ColumnType_BOOL,
				Description: "	Indicates whether this network is shared across all tenants.",
				Transform:   transform.FromField("Shared"),
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "The network status. Values are ACTIVE, DOWN, BUILD or ERROR.",
				Transform:   transform.FromField("Status"),
			},
			{
				Name:        "subnets",
				Type:        proto.ColumnType_STRING,
				Description: "The associated subnets.",
				Transform:   transform.FromField("Subnets").Transform(transform.EnsureStringArray),
			},
			{
				Name:        "vlan_transparent",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates the VLAN transparency mode of the network, which is VLAN transparent (true) or not VLAN transparent (false).",
				Transform:   transform.FromField("VLANTransparent"),
			},
			{
				Name:        "is_default",
				Type:        proto.ColumnType_STRING,
				Description: "The network is default pool or not.",
				Transform:   transform.FromField("IsDefault"),
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_STRING,
				Description: "Timestamp when the port was created.",
				Transform:   transform.FromField("CreatedAt").Transform(ToTime),
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_STRING,
				Description: "Timestamp when the port was last updated.",
				Transform:   transform.FromField("UpdatedAt").Transform(ToTime),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: "Tags is a list of security group tags. Tags are arbitrarily defined strings attached to a security group.",
				Transform:   transform.FromField("Tags"),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listOpenStackNetwork,
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
					Name:    "description",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "project_id",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "status",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "shared",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "admin_state_up",
					Require: plugin.Optional,
				},
				// TODO: add tags support
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getOpenStackNetwork,
		},
	}
}

//// LIST FUNCTION

func listOpenStackNetwork(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	plugin.Logger(ctx).Debug("retrieving openstack networks list", "query data", utils.ToPrettyJSON(d))

	client, err := getServiceClient(ctx, d, NetworkV2)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	opts := buildOpenStackNetworkFilter(ctx, d.EqualsQuals)

	allPages, err := networks.List(client, opts).AllPages()
	if err != nil {
		plugin.Logger(ctx).Error("error listing networks with options", "options", utils.ToPrettyJSON(opts), "error", err)
		return nil, err
	}
	allNetworks, err := networks.ExtractNetworks(allPages)
	if err != nil {
		plugin.Logger(ctx).Error("error extracting networks", "error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("networks retrieved", "count", len(allNetworks))

	for _, network := range allNetworks {
		if ctx.Err() != nil {
			plugin.Logger(ctx).Debug("context done, exit")
			break
		}
		network := network
		d.StreamListItem(ctx, &network)
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOpenStackNetwork(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	id := d.EqualsQuals["id"].GetStringValue()
	plugin.Logger(ctx).Debug("retrieving openstack network", "id", id)

	client, err := getServiceClient(ctx, d, NetworkV2)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	result := networks.Get(client, id)
	var network *networks.Network
	network, err = result.Extract()
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving network", "error", err)
		return nil, err
	}

	return network, nil
}

func buildOpenStackNetworkFilter(ctx context.Context, quals plugin.KeyColumnEqualsQualMap) networks.ListOpts {
	opts := networks.ListOpts{}

	if value, ok := quals["id"]; ok {
		opts.ID = value.GetStringValue()
	}
	if value, ok := quals["project_id"]; ok {
		opts.ProjectID = value.GetStringValue()
	}
	if value, ok := quals["name"]; ok {
		opts.Name = value.GetStringValue()
	}
	if value, ok := quals["description"]; ok {
		opts.Description = value.GetStringValue()
	}
	if value, ok := quals["status"]; ok {
		opts.Status = value.GetStringValue()
	}
	if value, ok := quals["admin_state_up"]; ok {
		opts.AdminStateUp = utils.PointerTo(value.GetBoolValue())
	}
	if value, ok := quals["shared"]; ok {
		opts.Shared = utils.PointerTo(value.GetBoolValue())
	}

	plugin.Logger(ctx).Debug("returning", "filter", utils.ToPrettyJSON(opts))
	return opts
}
