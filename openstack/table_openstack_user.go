package openstack

import (
	"context"

	"github.com/dihedron/steampipe-plugin-utils/utils"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableOpenStackUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openstack_user",
		Description: "OpenStack Users",
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique id of the user.",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the user.",
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The description of the user.",
				Transform:   transform.FromField("Description"),
			},
			{
				Name:        "default_project_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the default project of the user.",
				Transform:   transform.FromField("DefaultProjectID"),
			},
			{
				Name:        "domain_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the domain the user belongs to.",
				Transform:   transform.FromField("DomainID"),
			},
			{
				Name:        "enabled",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether or not the user is enabled.",
				Transform:   transform.FromField("Enabled"),
			},
			{
				Name:        "password_expires_at",
				Type:        proto.ColumnType_STRING,
				Description: "The timestamp when the user's password expires.",
				Transform:   transform.FromField("PasswordExpiresAt").Transform(ToTime),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listOpenStackUser,
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
					Name:    "domain_id",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "enabled",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "password_expires_at",
					Require: plugin.Optional,
				},
				// TODO: add tags support
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getOpenStackUser,
		},
	}
}

//// LIST FUNCTION

func listOpenStackUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	plugin.Logger(ctx).Debug("retrieving openstack users list", "query data", utils.ToPrettyJSON(d))

	client, err := getServiceClient(ctx, d, IdentityV3)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	opts := buildOpenStackUserFilter(ctx, d.EqualsQuals)

	allPages, err := users.List(client, opts).AllPages()
	if err != nil {
		plugin.Logger(ctx).Error("error listing users with options", "options", utils.ToPrettyJSON(opts), "error", err)
		return nil, err
	}
	allUsers, err := users.ExtractUsers(allPages)
	if err != nil {
		plugin.Logger(ctx).Error("error extracting users", "error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("users retrieved", "count", len(allUsers))

	for _, user := range allUsers {
		if ctx.Err() != nil {
			plugin.Logger(ctx).Debug("context done, exit")
			break
		}
		user := user
		d.StreamListItem(ctx, &user)
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOpenStackUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	id := d.EqualsQuals["id"].GetStringValue()
	plugin.Logger(ctx).Debug("retrieving openstack user", "id", id)

	client, err := getServiceClient(ctx, d, IdentityV3)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	result := users.Get(client, id)
	var user *users.User
	user, err = result.Extract()
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving user", "error", err)
		return nil, err
	}

	return user, nil
}

func buildOpenStackUserFilter(ctx context.Context, quals plugin.KeyColumnEqualsQualMap) users.ListOpts {
	opts := users.ListOpts{}
	if value, ok := quals["id"]; ok {
		opts.UniqueID = value.GetStringValue()
	}
	if value, ok := quals["name"]; ok {
		opts.Name = value.GetStringValue()
	}
	if value, ok := quals["domain_id"]; ok {
		opts.DomainID = value.GetStringValue()
	}
	if value, ok := quals["enabled"]; ok {
		opts.Enabled = utils.PointerTo(value.GetBoolValue())
	}
	if value, ok := quals["password_expires_at"]; ok {
		opts.PasswordExpiresAt = value.GetStringValue()
	}
	plugin.Logger(ctx).Debug("returning", "filter", utils.ToPrettyJSON(opts))
	return opts
}
