package openstack

import (
	"context"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
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
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the user.",
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The description of the user.",
			},
			{
				Name:        "default_project_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the default project of the user.",
			},
			{
				Name:        "domain_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the domain the user belongs to.",
			},
			{
				Name:        "enabled",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether or not the user is enabled.",
			},
			{
				Name:        "password_expires_at",
				Type:        proto.ColumnType_STRING,
				Description: "The timestamp when the user's password expires.",
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

// openstackUser is the struct representing the result of the list and hydrate functions.
type openstackUser struct {
	ID                string
	Name              string
	Description       string
	DomainID          string
	DefaultProjectID  string
	Enabled           bool
	PasswordExpiresAt string
	// // Extra is a collection of miscellaneous key/values.
	// Extra map[string]interface{} `json:"-"`
	// // Options are a set of defined options of the user.
	// Options map[string]interface{} `json:"options"`
}

//// LIST FUNCTION

func listOpenStackUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	plugin.Logger(ctx).Debug("retrieving openstack users list", "query data", toPrettyJSON(d))

	client, err := getServiceClient(ctx, d, IdentityV3)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	opts := buildOpenStackUserFilter(ctx, d.KeyColumnQuals)

	allPages, err := users.List(client, opts).AllPages()
	if err != nil {
		plugin.Logger(ctx).Error("error listing users with options", "options", toPrettyJSON(opts), "error", err)
		return nil, err
	}
	allProjects, err := users.ExtractUsers(allPages)
	if err != nil {
		plugin.Logger(ctx).Error("error extracting users", "error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("users retrieved", "count", len(allProjects))

	for _, project := range allProjects {
		d.StreamListItem(ctx, buildOpenStackUser(ctx, &project))
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOpenStackUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	id := d.KeyColumnQuals["id"].GetStringValue()
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

	return buildOpenStackUser(ctx, user), nil
}

func buildOpenStackUser(ctx context.Context, user *users.User) *openstackUser {
	result := &openstackUser{
		ID:          user.ID,
		Name:        user.Name,
		Description: user.Description,
		// IsDomain:    project.IsDomain,
		// DomainID:    project.DomainID,
		// Enabled:     project.Enabled,
		// ParentID:    project.ParentID,
	}
	plugin.Logger(ctx).Debug("returning user", "user", toPrettyJSON(result))
	return result
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
		opts.Enabled = pointerTo(value.GetBoolValue())
	}
	if value, ok := quals["password_expires_at"]; ok {
		opts.PasswordExpiresAt = value.GetStringValue()
	}
	plugin.Logger(ctx).Debug("returning", "filter", toPrettyJSON(opts))
	return opts
}
