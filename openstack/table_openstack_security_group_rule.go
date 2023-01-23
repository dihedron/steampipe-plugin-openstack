package openstack

import (
	"context"

	"github.com/dihedron/steampipe-plugin-utils/utils"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableOpenStackSecurityGroupRule(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openstack_security_group_rule",
		Description: "OpenStack Security Group Rule",

		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The ID of the security group rule we're retrieving.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The description of the instance",
				Transform:   transform.FromField("Description"),
			},
			{
				Name:        "project_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the project.",
				Transform:   transform.FromField("TenantID"),
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_STRING,
				Description: "Time at which the security group rule has been created (in UTC ISO8601 format).",
				Transform:   transform.FromField("CreatedAt").Transform(ToTime),
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_STRING,
				Description: "Time at which the security group rule has been updated (in UTC ISO8601 format).",
				Transform:   transform.FromField("UpdatedAt").Transform(ToTime),
			},
			{
				Name:        "remote_group_id",
				Type:        proto.ColumnType_STRING,
				Description: "The remote group UUID to associate with this security group rule.",
				Transform:   transform.FromField("RemoteGroupID"),
			},
			{
				Name:        "direction",
				Type:        proto.ColumnType_STRING,
				Description: "Ingress or egress, which is the direction in which the security group rule is applied.",
				Transform:   transform.FromField("Direction"),
			},
			{
				Name:        "protocol",
				Type:        proto.ColumnType_STRING,
				Description: "The IP protocol, either as a string (e.g. 'dhcp') or as the port (e.g. '53').",
				Transform:   transform.FromField("Protocol"),
			},
			{
				Name:        "ethertype",
				Type:        proto.ColumnType_STRING,
				Description: "Must be IPv4 or IPv6, and addresses represented in CIDR must match the ingress or egress rules.",
				Transform:   transform.FromField("EtherType"),
			},
			{
				Name:        "port_range_min",
				Type:        proto.ColumnType_INT,
				Description: "The minimum port number in the range that is matched by the security group rule.",
				Transform:   transform.FromField("PortRangeMin"),
			},
			{
				Name:        "port_range_max",
				Type:        proto.ColumnType_INT,
				Description: "The maximum port number in the range that is matched by the security group rule.",
				Transform:   transform.FromField("PortRangeMax"),
			},
			{
				Name:        "security_group_id",
				Type:        proto.ColumnType_STRING,
				Description: "The security group ID to associate with this security group rule.",
				Transform:   transform.FromField("SecurityGroupID"),
			},
			{
				Name:        "remote_ip_prefix",
				Type:        proto.ColumnType_STRING,
				Description: "The remote IP prefix that is matched by this security group rule.",
				Transform:   transform.FromField("RemoteIPPrefix"),
			},
			{
				Name:        "revision",
				Type:        proto.ColumnType_INT,
				Description: "The revision number of the resource.",
				Transform:   transform.FromField("Revision"),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listOpenStackSecurityGroupRule,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getOpenStackSecurityGroupRule,
		},
	}
}

//// LIST FUNCTION

func listOpenStackSecurityGroupRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	plugin.Logger(ctx).Debug("list security groups rules", "query data", utils.ToPrettyJSON(d), "hydrate data", utils.ToPrettyJSON(h))

	client, err := getServiceClient(ctx, d, NetworkV2)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	opts := buildOpenStackSecurityGroupRuleFilter(ctx, d.EqualsQuals)

	allPages, err := rules.List(client, opts).AllPages()
	if err != nil {
		plugin.Logger(ctx).Error("error listing security group rules with options", "options", utils.ToPrettyJSON(opts), "error", err)
		return nil, err
	}
	allRules, err := rules.ExtractRules(allPages)
	if err != nil {
		plugin.Logger(ctx).Error("error extracting rules", "error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("rules retrieved", "count", len(allRules))

	for _, rule := range allRules {
		if ctx.Err() != nil {
			plugin.Logger(ctx).Debug("context done, exit")
			break
		}
		rule := rule

		d.StreamListItem(ctx, &rule)
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOpenStackSecurityGroupRule(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	id := d.EqualsQuals["id"].GetStringValue()
	plugin.Logger(ctx).Debug("retrieving openstack security group", "id", id)

	client, err := getServiceClient(ctx, d, NetworkV2)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	result := groups.Get(client, id)
	var group *groups.SecGroup
	group, err = result.Extract()
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving security group", "error", err)
		return nil, err
	}

	return group, nil
}

func buildOpenStackSecurityGroupRuleFilter(ctx context.Context, quals plugin.KeyColumnEqualsQualMap) rules.ListOpts {
	opts := rules.ListOpts{}

	//  (Optional)	query	string
	// revision_number (Optional)	query	integer

	if value, ok := quals["id"]; ok {
		// Filter the list result by the ID of the resource.
		opts.ID = value.GetStringValue()
	}
	if value, ok := quals["description"]; ok {
		// Filter the list result by the human-readable description of the resource.
		opts.Description = value.GetStringValue()
	}
	if value, ok := quals["direction"]; ok {
		// Filter the security group rule list result by the direction in which
		// the security group rule is applied, which is ingress or egress.
		opts.Direction = value.GetStringValue()
	}
	if value, ok := quals["ethertype"]; ok {
		// Filter the security group rule list result by the ethertype of network traffic.
		// The value must be IPv4 or IPv6.
		opts.EtherType = value.GetStringValue()
	}
	if value, ok := quals["protocol"]; ok {
		// Filter the security group rule list result by the IP protocol.
		opts.Protocol = value.GetStringValue()
	}
	if value, ok := quals["project_id"]; ok {
		// Filter the list result by the ID of the project that owns the resource.
		opts.ProjectID = value.GetStringValue()
	}
	if value, ok := quals["tenant_id"]; ok {
		// Filter the list result by the ID of the project that owns the resource.
		opts.TenantID = value.GetStringValue()
	}
	if value, ok := quals["remote_group_id"]; ok {
		// Filter the security group rule list result by the ID of the remote group
		// that associates with this security group rule.
		opts.RemoteGroupID = value.GetStringValue()
	}
	if value, ok := quals["port_range_max"]; ok {
		// Filter the security group rule list result by the maximum port number in
		// the range that is matched by the security group rule.
		opts.PortRangeMax = int(value.GetInt64Value())
	}
	if value, ok := quals["port_range_min"]; ok {
		// Filter the security group rule list result by the minimum port number in
		// the range that is matched by the security group rule.
		opts.PortRangeMin = int(value.GetInt64Value())
	}
	if value, ok := quals["security_group_id"]; ok {
		// Filter the security group rule list result by the ID of the security group
		// that associates with this security group rule.
		opts.SecGroupID = value.GetStringValue()
	}
	if value, ok := quals["remote_ip_prefix"]; ok {
		// Filter the list result by the remote IP prefix that is matched by this security
		// group rule.
		opts.RemoteIPPrefix = value.GetStringValue()
	}

	// TODO: handle tags
	plugin.Logger(ctx).Debug("returning", "filter", utils.ToPrettyJSON(opts))
	return opts
}

// func extractSecGroupRuleIDs(_ context.Context, d *transform.TransformData) (interface{}, error) {
// 	var values []string
// 	if d.Value != nil {
// 		if list, ok := d.Value.([]rules.SecGroupRule); ok {
// 			for _, sgr := range list {
// 				values = append(values, sgr.ID)
// 			}
// 		}
// 	}
// 	return values, nil
// }
