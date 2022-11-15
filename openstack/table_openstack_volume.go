package openstack

import (
	"context"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableOpenStackVolume(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openstack_volume",
		Description: "OpenStack Disk Volume",
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique id of the volume.",
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Human-readable name for the volume. Might not be unique.",
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The description of the project (or tenant)",
			},
			{
				Name:        "project_id",
				Type:        proto.ColumnType_STRING,
				Description: "The id of the project the volume belongs to.",
			},
			{
				Name:        "user_id",
				Type:        proto.ColumnType_STRING,
				Description: "The id of the user who created the volume.",
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "Indicates the current status of the volume.",
			},
			{
				Name:        "replication_status",
				Type:        proto.ColumnType_STRING,
				Description: "Indicates the status of replication of the volume.",
			},
			{
				Name:        "size",
				Type:        proto.ColumnType_INT,
				Description: "Size of the volume in GB.",
			},
			{
				Name:        "availability_zone",
				Type:        proto.ColumnType_STRING,
				Description: "AvailabilityZone is which availability zone the volume is in.",
			},
			{
				Name:        "bootable",
				Type:        proto.ColumnType_STRING, // TODO: check if convertible to BOOL
				Description: "Indicates whether this is a bootable volume.",
			},
			{
				Name:        "encrypted",
				Type:        proto.ColumnType_BOOL,
				Description: "Denotes if the volume is encrypted.",
			},
			{
				Name:        "multiattach",
				Type:        proto.ColumnType_BOOL,
				Description: "denotes if the volume is multi-attach capable.",
			},
			{
				Name:        "consistencygroup_id",
				Type:        proto.ColumnType_STRING,
				Description: "The volume's consistency group id.",
			},
			{
				Name:        "volume_type",
				Type:        proto.ColumnType_STRING,
				Description: "The type of volume to create, either SATA or SSD.",
			},
			{
				Name:        "snapshot_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the snapshot from which the volume was created.",
			},
			{
				Name:        "source_vol_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of another block storage volume from which the current volume was created.",
			},
			{
				Name:        "backup_id",
				Type:        proto.ColumnType_STRING,
				Description: "The backup ID, from which the volume was restored; this value is available starting from microversion 3.47.",
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_STRING,
				Description: "Timestamp when the volume was created.",
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_STRING,
				Description: "The date when this volume was last updated.",
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listOpenStackVolume,
			KeyColumns: plugin.KeyColumnSlice{
				// &plugin.KeyColumn{
				// 	Name:    "id",
				// 	Require: plugin.Optional,
				// },
				&plugin.KeyColumn{
					Name:    "name",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "status",
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
			Hydrate:    getOpenStackVolume,
		},
	}
}

// openstackPort is the struct representing the result of the list and hydrate functions.
type openstackVolume struct {
	ID                 string
	Name               string
	Description        string
	UserID             string
	ProjectID          string
	Status             string
	ReplicationStatus  string
	Bootable           string
	Encrypted          bool
	Multiattach        bool
	Size               int
	ConsistencyGroupID string
	AvailabilityZone   string
	VolumeType         string
	SnapshotID         string
	SourceVolID        string
	BackupID           string
	CreatedAt          string
	UpdatedAt          string
	// // Instances onto which the volume is attached.
	// Attachments []Attachment `json:"attachments"`
	// // Arbitrary key-value pairs defined by the user.
	// Metadata map[string]string `json:"metadata"`
	// // Image metadata entries, only included for volumes that were created from an image, or from a snapshot of a volume originally created from an image.
	// VolumeImageMetadata map[string]string `json:"volume_image_metadata"`
}

//// LIST FUNCTION

func listOpenStackVolume(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	plugin.Logger(ctx).Debug("retrieving openstack volumes list", "query data", toPrettyJSON(d))

	client, err := getServiceClient(ctx, d, BlockStorageV3)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	opts := buildOpenStackVolumeFilter(ctx, d.KeyColumnQuals)

	allPages, err := volumes.List(client, opts).AllPages()
	if err != nil {
		plugin.Logger(ctx).Error("error listing volumes with options", "options", toPrettyJSON(opts), "error", err)
		return nil, err
	}
	allVolumes, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		plugin.Logger(ctx).Error("error extracting volumes", "error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("volumes retrieved", "count", len(allVolumes))

	for _, volume := range allVolumes {
		d.StreamListItem(ctx, buildOpenStackVolume(ctx, &volume))
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOpenStackVolume(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	id := d.KeyColumnQuals["id"].GetStringValue()
	plugin.Logger(ctx).Debug("retrieving openstack volume", "id", id)

	client, err := getServiceClient(ctx, d, BlockStorageV3)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	result := volumes.Get(client, id)
	var volume *volumes.Volume
	volume, err = result.Extract()
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving volume", "error", err)
		return nil, err
	}

	return buildOpenStackVolume(ctx, volume), nil
}

func buildOpenStackVolume(ctx context.Context, volume *volumes.Volume) *openstackVolume {
	result := &openstackVolume{
		ID:               volume.ID,
		Name:             volume.Name,
		Description:      volume.Description,
		Status:           volume.Status,
		Bootable:         volume.Bootable,
		Size:             volume.Size,
		AvailabilityZone: volume.AvailabilityZone,
		VolumeType:       volume.VolumeType,
		SnapshotID:       volume.SnapshotID,
		SourceVolID:      volume.SourceVolID,
		CreatedAt:        volume.CreatedAt.String(),
	}
	plugin.Logger(ctx).Debug("returning volume", "volume", toPrettyJSON(result))
	return result
}

func buildOpenStackVolumeFilter(ctx context.Context, quals plugin.KeyColumnEqualsQualMap) volumes.ListOpts {
	opts := volumes.ListOpts{
		AllTenants: true,
	}
	if value, ok := quals["name"]; ok {
		opts.Name = value.GetStringValue()
	}
	if value, ok := quals["status"]; ok {
		opts.Status = value.GetStringValue()
	}
	if value, ok := quals["project_id"]; ok {
		opts.TenantID = value.GetStringValue()
	}
	// TODO: add metadata
	plugin.Logger(ctx).Debug("returning", "filter", toPrettyJSON(opts))
	return opts
}
