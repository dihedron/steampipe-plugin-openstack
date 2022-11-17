package openstack

import (
	"context"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableOpenStackProject(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openstack_project",
		Description: "OpenStack Project (aka Tenant)",
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique id of the project (or tenant).",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the project (or tenant).",
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The description of the project (or tenant)",
				Transform:   transform.FromField("Description"),
			},
			{
				Name:        "is_domain",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether the project is a domain.",
				Transform:   transform.FromField("IsDomain"),
			},
			{
				Name:        "domain_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the domain the project belongs to.",
				Transform:   transform.FromField("DomainID"),
			},
			{
				Name:        "enabled",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether or not the project is enabled.",
				Transform:   transform.FromField("Enabled"),
			},
			{
				Name:        "parent_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the parent project.",
				Transform:   transform.FromField("ParentID"),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listOpenStackProject,
			KeyColumns: plugin.KeyColumnSlice{
				&plugin.KeyColumn{
					Name:    "name",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "is_domain",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "domain_id",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "enabled",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "parent_id",
					Require: plugin.Optional,
				},
				// TODO: add tags support
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getOpenStackProject,
		},
	}
}

//// LIST FUNCTION

func listOpenStackProject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	plugin.Logger(ctx).Debug("retrieving openstack projects list", "query data", toPrettyJSON(d))

	client, err := getServiceClient(ctx, d, IdentityV3)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	opts := buildOpenStackProjectFilter(ctx, d.KeyColumnQuals)

	allPages, err := projects.List(client, opts).AllPages()
	if err != nil {
		plugin.Logger(ctx).Error("error listing projects with options", "options", toPrettyJSON(opts), "error", err)
		return nil, err
	}
	allProjects, err := projects.ExtractProjects(allPages)
	if err != nil {
		plugin.Logger(ctx).Error("error extracting projects", "error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("projects retrieved", "count", len(allProjects))

	for _, project := range allProjects {
		d.StreamListItem(ctx, &project)
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOpenStackProject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	id := d.KeyColumnQuals["id"].GetStringValue()
	plugin.Logger(ctx).Debug("retrieving openstack project", "id", id)

	client, err := getServiceClient(ctx, d, IdentityV3)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	result := projects.Get(client, id)
	var project *projects.Project
	project, err = result.Extract()
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving project", "error", err)
		return nil, err
	}

	return project, nil
}

func buildOpenStackProjectFilter(ctx context.Context, quals plugin.KeyColumnEqualsQualMap) projects.ListOpts {
	opts := projects.ListOpts{}
	if value, ok := quals["name"]; ok {
		opts.Name = value.GetStringValue()
	}
	if value, ok := quals["is_domain"]; ok {
		opts.IsDomain = pointerTo(value.GetBoolValue())
	}
	if value, ok := quals["domain_id"]; ok {
		opts.DomainID = value.GetStringValue()
	}
	if value, ok := quals["enabled"]; ok {
		opts.Enabled = pointerTo(value.GetBoolValue())
	}
	if value, ok := quals["parent_id"]; ok {
		opts.ParentID = value.GetStringValue()
	}
	plugin.Logger(ctx).Debug("returning", "filter", toPrettyJSON(opts))
	return opts
}
