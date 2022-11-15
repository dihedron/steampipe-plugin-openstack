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
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "Indicates the current status of the volume.",
			},
			{
				Name:        "bootable",
				Type:        proto.ColumnType_STRING, // TODO: check if convertible to BOOL
				Description: "Indicates whether this is a bootable volume.",
			},
			{
				Name:        "size",
				Type:        proto.ColumnType_INT,
				Description: "Size of the volume in GB.",
			},
			{
				Name:        "availability_zone",
				Type:        proto.ColumnType_STRING,
				Description: "Availability zone of the volume; this parameter is no longer used.",
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
				Name:        "created_at",
				Type:        proto.ColumnType_STRING,
				Description: "Timestamp when the volume was created.",
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listOpenStackVolume,
			KeyColumns: plugin.KeyColumnSlice{
				&plugin.KeyColumn{
					Name:    "id",
					Require: plugin.Optional,
				},
				// &plugin.KeyColumn{
				// 	Name:    "name",
				// 	Require: plugin.Optional,
				// },
				// &plugin.KeyColumn{
				// 	Name:    "status",
				// 	Require: plugin.Optional,
				// },
				// &plugin.KeyColumn{
				// 	Name:    "description",
				// 	Require: plugin.Optional,
				// },
				// &plugin.KeyColumn{
				// 	Name:    "admin_state_up",
				// 	Require: plugin.Optional,
				// },
				// &plugin.KeyColumn{
				// 	Name:    "network_id",
				// 	Require: plugin.Optional,
				// },
				// &plugin.KeyColumn{
				// 	Name:    "project_id",
				// 	Require: plugin.Optional,
				// },
				// &plugin.KeyColumn{
				// 	Name:    "device_owner",
				// 	Require: plugin.Optional,
				// },
				// &plugin.KeyColumn{
				// 	Name:    "device_id",
				// 	Require: plugin.Optional,
				// },
				// &plugin.KeyColumn{
				// 	Name:    "mac_address",
				// 	Require: plugin.Optional,
				// },
				// TODO: add tags support
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
	ID               string
	Name             string
	Description      string
	Status           string
	Bootable         string
	Size             int
	AvailabilityZone string
	VolumeType       string
	SnapshotID       string
	SourceVolID      string
	CreatedAt        string
	// // Instances onto which the volume is attached.
	// Attachments []map[string]interface{} `json:"attachments"`
	// // Arbitrary key-value pairs defined by the user.
	// Metadata map[string]string `json:"metadata"`
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
	plugin.Logger(ctx).Debug("returning", "filter", toPrettyJSON(opts))
	return opts
}
