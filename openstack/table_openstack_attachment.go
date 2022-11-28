package openstack

import (
	"context"

	"github.com/dihedron/steampipe-plugin-utils/utils"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/attachments"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableOpenStackAttachment(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openstack_attachment",
		Description: "OpenStack Disk Volume Attachment",
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique id of the attachment.",
				Transform:   transform.FromField("ID"),
			},

			// ConnectionInfo struct {
			// 	AccessMode       string   `json:"access_mode"`
			// 	AttachmentID     string   `json:"attachment_id"`
			// 	AuthEnabled      bool     `json:"auth_enabled"`
			// 	AuthUsername     string   `json:"auth_username"`
			// 	ClusterName      string   `json:"cluster_name"`
			// 	Discard          bool     `json:"discard"`
			// 	DriverVolumeType string   `json:"driver_volume_type"`
			// 	Encrypted        bool     `json:"encrypted"`
			// 	Hosts            []string `json:"hosts"`
			// 	Keyring          string   `json:"keyring"`
			// 	Name             string   `json:"name"`
			// 	Ports            []string `json:"ports"`
			// 	SecretType       string   `json:"secret_type"`
			// 	SecretUUID       string   `json:"secret_uuid"`
			// 	VolumeID         string   `json:"volume_id"`
			// } `json:"connection_info"`

			{
				Name:        "attached_at",
				Type:        proto.ColumnType_STRING,
				Description: "When the attachment was created.",
				Transform:   transform.FromField("AttachedAt").Transform(ToTime),
			},
			{
				Name:        "detached_at",
				Type:        proto.ColumnType_STRING,
				Description: "When the attachment was destroyed.",
				Transform:   transform.FromField("DetachedAt").Transform(ToTime),
			},
			{
				Name:        "attachment_id",
				Type:        proto.ColumnType_STRING,
				Description: "The identifier of the attachment.",
				Transform:   transform.FromField("AttachmentID"),
			},
			{
				Name:        "volume_id",
				Type:        proto.ColumnType_STRING,
				Description: "The id of the volume.",
				Transform:   transform.FromField("VolumeID"),
			},
			{
				Name:        "instance_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the instance owning the attachment.",
				Transform:   transform.FromField("Instance"),
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "Indicates the current status of the volume.",
				Transform:   transform.FromField("Status"),
			},
			{
				Name:        "attach_mode",
				Type:        proto.ColumnType_STRING,
				Description: "The attach mode of attachment, read-only ('ro') or read-and-write ('rw'), default is 'rw'.",
				Transform:   transform.FromField("AttachMode"),
			},
			{
				Name:        "project_id",
				Type:        proto.ColumnType_STRING,
				Description: "The id of the project the attachment belongs to.",
				Transform:   transform.FromField("ProjectID"), //FromField("OsVolTenantAttrTenantID"),
			},
			{
				Name:        "connection_info",
				Type:        proto.ColumnType_JSON,
				Description: "The connection info used for server to connect the volume.",
				Transform:   transform.FromField("ConnectionInfo"),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listOpenStackAttachment,
			KeyColumns: plugin.KeyColumnSlice{
				&plugin.KeyColumn{
					Name:    "instance_id",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "volume_id",
					Require: plugin.Optional,
				},
				//
				// NOTE: does not seem to work as documented!
				//
				// &plugin.KeyColumn{
				// 	Name:    "status",
				// 	Require: plugin.Optional,
				// },
				&plugin.KeyColumn{
					Name:    "project_id",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			Hydrate: getOpenStackAttachment,
			KeyColumns: plugin.KeyColumnSlice{
				&plugin.KeyColumn{
					Name:    "id",
					Require: plugin.Required,
				},
				&plugin.KeyColumn{
					Name:    "project_id",
					Require: plugin.Optional,
				},
			},
		},
	}
}

//// LIST FUNCTION

func listOpenStackAttachment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	plugin.Logger(ctx).Debug("retrieving openstack attachment list", "query data", utils.ToPrettyJSON(d))

	client, err := getServiceClient(ctx, d, BlockStorageV3)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	// the OpenStack Cinder v2 API required that the project_id be specified in
	// the request path; this can be cumbersome when working with SQL, so if the
	// user did NOT specify the project_id filter, we get a list of all project
	// IDs and then loop over them all, one by one. Therefore, the filter function
	// will NOT handle the project_id filter because we set it ourselves.
	projectIDs := []string{}

	opts := buildOpenStackAttachmentFilter(ctx, d.EqualsQuals)
	if projectID, ok := d.EqualsQuals["id"]; ok {
		projectIDs = append(projectIDs, projectID.GetStringValue())
	} else {

		client, err := getServiceClient(ctx, d, IdentityV3)
		if err != nil {
			plugin.Logger(ctx).Error("error retrieving client", "error", err)
			return nil, err
		}

		allPages, err := projects.List(client, &projects.ListOpts{}).AllPages()
		if err != nil {
			plugin.Logger(ctx).Error("error listing projects", "error", err)
			return nil, err
		}
		allProjects, err := projects.ExtractProjects(allPages)
		if err != nil {
			plugin.Logger(ctx).Error("error extracting projects", "error", err)
			return nil, err
		}
		plugin.Logger(ctx).Debug("projects retrieved", "count", len(allProjects))

		for _, project := range allProjects {
			if ctx.Err() != nil {
				plugin.Logger(ctx).Debug("context done, exit")
				break
			}
			project := project
			projectIDs = append(projectIDs, project.ID)
		}
	}

	for _, projectID := range projectIDs {
		opts.ProjectID = projectID

		allPages, err := attachments.List(client, opts).AllPages()
		if err != nil {
			plugin.Logger(ctx).Error("error listing attachments with options", "options", utils.ToPrettyJSON(opts), "error", err)
			return nil, err
		}
		allAttachments := []*apiAttachment{}
		err = attachments.ExtractAttachmentsInto(allPages, &allAttachments)
		if err != nil {
			plugin.Logger(ctx).Error("error extracting attachment", "error", err)
			return nil, err
		}
		plugin.Logger(ctx).Debug("attachment retrieved", "count", len(allAttachments))

		for _, attachment := range allAttachments {
			if ctx.Err() != nil {
				plugin.Logger(ctx).Debug("context done, exit")
				break
			}
			attachment := attachment
			attachment.ProjectID = projectID
			d.StreamListItem(ctx, attachment)
		}
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOpenStackAttachment(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	id := d.EqualsQuals["id"].GetStringValue()

	plugin.Logger(ctx).Debug("retrieving openstack attachment", "id", id)

	client, err := getServiceClient(ctx, d, BlockStorageV3)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	result := attachments.Get(client, id)
	//plugin.Logger(ctx).Debug("request run", "result", utils.ToPrettyJSON(result))

	attachment := &apiAttachment{}
	err = result.ExtractInto(attachment)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving attachment", "error", err)
		return nil, err
	}

	if value, ok := d.EqualsQuals["project_id"]; ok {
		attachment.ProjectID = value.GetStringValue()
	}

	return attachment, nil
}
func buildOpenStackAttachmentFilter(ctx context.Context, quals plugin.KeyColumnEqualsQualMap) attachments.ListOpts {
	opts := attachments.ListOpts{
		AllTenants: true,
	}
	if value, ok := quals["instance_id"]; ok {
		opts.InstanceID = value.GetStringValue()
	}
	if value, ok := quals["volume_id"]; ok {
		opts.VolumeID = value.GetStringValue()
	}
	//
	// NOTE: does not seem to work as documented!
	//
	// if value, ok := quals["status"]; ok {
	// 	opts.Status = value.GetStringValue()
	// }
	if value, ok := quals["project_id"]; ok {
		opts.ProjectID = value.GetStringValue()
	}
	plugin.Logger(ctx).Debug("returning", "filter", utils.ToPrettyJSON(opts))
	return opts
}

type apiAttachment struct {
	ID             string `json:"id"`
	AttachedAt     Time   `json:"attached_at"`
	DetachedAt     Time   `json:"detached_at"`
	AttachmentID   string `json:"attachment_id"`
	VolumeID       string `json:"volume_id"`
	Instance       string `json:"instance"`
	Status         string `json:"status"`
	AttachMode     string `json:"attach_mode"`
	ProjectID      string `json:"-"`
	ConnectionInfo struct {
		AccessMode       string   `json:"access_mode"`
		AttachmentID     string   `json:"attachment_id"`
		AuthEnabled      bool     `json:"auth_enabled"`
		AuthUsername     string   `json:"auth_username"`
		ClusterName      string   `json:"cluster_name"`
		Discard          bool     `json:"discard"`
		DriverVolumeType string   `json:"driver_volume_type"`
		Encrypted        bool     `json:"encrypted"`
		Hosts            []string `json:"hosts"`
		Keyring          string   `json:"keyring"`
		Name             string   `json:"name"`
		Ports            []string `json:"ports"`
		SecretType       string   `json:"secret_type"`
		SecretUUID       string   `json:"secret_uuid"`
		VolumeID         string   `json:"volume_id"`
	} `json:"connection_info"`
}
