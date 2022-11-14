package openstack

import (
	"context"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableOpenStackProject(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openstack_project",
		Description: "OpenStack Project (aka Tenant)",
		List: &plugin.ListConfig{
			Hydrate: listOpenStackProject,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			// IgnoreConfig: &plugin.IgnoreConfig{
			// 	ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInstanceID.NotFound", "InvalidInstanceID.Unavailable", "InvalidInstanceID.Malformed"}),
			// },
			Hydrate: getOpenStackProject,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The project (or tenant) id"},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the project (or tenant)"},
			// {Name: "project_id", Type: proto.ColumnType_STRING, Description: "The ID of the instance's project (aka tenant)"},
			// {Name: "user_id", Type: proto.ColumnType_STRING, Description: "The ID of the instance's user"},
			// {Name: "created_at", Type: proto.ColumnType_STRING, Description: "The creation time of the instance"},
			// {Name: "launched_at", Type: proto.ColumnType_STRING, Description: "The launch time of the instance"},
			// {Name: "updated_at", Type: proto.ColumnType_STRING, Description: "The update time of the instance"},
			// {Name: "terminated_at", Type: proto.ColumnType_STRING, Description: "The termintaion time of the instance"},
			// {Name: "host_id", Type: proto.ColumnType_STRING, Description: "The ID of the hypervisor (host) the instance is running on"},
			// {Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the instance"},
			// {Name: "progress", Type: proto.ColumnType_INT, Description: "Progress information about the instance"},
			// AccessIPv4   string                 `json:"accessIPv4"`
			// AccessIPv6   string                 `json:"accessIPv6"`

		},
	}
}

//// LIST FUNCTION

func listOpenStackProject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("retrieving openstack projects list")
	plugin.Logger(ctx).Debug("plugin query data: %s", toPrettyJSON(d))
	plugin.Logger(ctx).Debug("plugin hydrate data %s", toPrettyJSON(h))
	return nil, ErrNotImplemented
}

//// HYDRATE FUNCTIONS

func getOpenStackProject(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	id := d.KeyColumnQuals["id"].GetStringValue()
	plugin.Logger(ctx).Debug("retrieving openstack project", "id", id)

	client, err := getIdentityV3Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("error creating identity v3 client", "error", err)
		return nil, err
	}

	result := projects.Get(client, id)
	var project *projects.Project
	project, err = result.Extract()
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving project", "error", err)
		return nil, err
	}

	return &openstackProject{
		ID:   id,
		Name: project.Name,
	}, nil
}

type openstackProject struct {
	ID   string
	Name string
	// ProjectID    string
	// UserID       string
	// CreatedAt    string
	// LaunchedAt   string
	// UpdatedAt    string
	// TerminatedAt string
	// HostID       string
	// Status       string
	// Progress     int
}
