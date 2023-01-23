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

func tableOpenStackSecurityGroup(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openstack_security_group",
		Description: "OpenStack Security Group",
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "The ID of the security group we're retrieving.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the instance",
				Transform:   transform.FromField("Name"),
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
				Description: "The ID of the instance's project (aka tenant)",
				Transform:   transform.FromField("TenantID"),
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_STRING,
				Description: "The creation time of the security group",
				Transform:   transform.FromField("CreatedAt").Transform(ToTime),
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_STRING,
				Description: "The update time of the security group",
				Transform:   transform.FromField("UpdatedAt").Transform(ToTime),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: "Tags is a list of security group tags. Tags are arbitrarily defined strings attached to a security group.",
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "security_group_rule_ids",
				Type:        proto.ColumnType_JSON,
				Description: "The id of the security group rules that belong to the current security group.",
				Transform:   transform.FromField("Rules").Transform(extractSecGroupRuleIDs), //.Transform(transform.EnsureStringArray)
			},
			{
				Name:        "security_group_rules",
				Type:        proto.ColumnType_JSON,
				Description: "The security group rules that belong to the current security group.",
				Transform:   transform.FromField("Rules"), //.Transform(transform.EnsureStringArray),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listOpenStackSecurityGroup,
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
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getOpenStackSecurityGroup,
		},
	}
}

//// LIST FUNCTION

func listOpenStackSecurityGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	plugin.Logger(ctx).Debug("list security groups", "query data", utils.ToPrettyJSON(d), "hydrate data", utils.ToPrettyJSON(h))

	client, err := getServiceClient(ctx, d, NetworkV2)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	opts := buildOpenStackSecurityGroupFilter(ctx, d.EqualsQuals)

	allPages, err := groups.List(client, opts).AllPages()
	if err != nil {
		plugin.Logger(ctx).Error("error listing security groups with options", "options", utils.ToPrettyJSON(opts), "error", err)
		return nil, err
	}
	allGroups, err := groups.ExtractGroups(allPages)
	if err != nil {
		plugin.Logger(ctx).Error("error extracting groups", "error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("groups retrieved", "count", len(allGroups))

	for _, group := range allGroups {
		if ctx.Err() != nil {
			plugin.Logger(ctx).Debug("context done, exit")
			break
		}
		group := group
		d.StreamListItem(ctx, &group)
	}
	return nil, nil

}

//// HYDRATE FUNCTIONS

func getOpenStackSecurityGroup(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

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

func buildOpenStackSecurityGroupFilter(ctx context.Context, quals plugin.KeyColumnEqualsQualMap) groups.ListOpts {
	opts := groups.ListOpts{}

	if value, ok := quals["id"]; ok {
		opts.ID = value.GetStringValue()
	}
	if value, ok := quals["name"]; ok {
		opts.Name = value.GetStringValue()
	}
	if value, ok := quals["description"]; ok {
		opts.Description = value.GetStringValue()
	}
	if value, ok := quals["marker"]; ok {
		opts.Marker = value.GetStringValue()
	}
	if value, ok := quals["project_id"]; ok {
		opts.ProjectID = value.GetStringValue()
	}

	// TODO: handle tags
	plugin.Logger(ctx).Debug("returning", "filter", utils.ToPrettyJSON(opts))
	return opts
}

func extractSecGroupRuleIDs(_ context.Context, d *transform.TransformData) (interface{}, error) {
	var values []string
	if d.Value != nil {
		if list, ok := d.Value.([]rules.SecGroupRule); ok {
			for _, sgr := range list {
				values = append(values, sgr.ID)
			}
		}
	}
	return values, nil
}
