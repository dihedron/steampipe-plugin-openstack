package openstack

import (
	"context"
	"strconv"
	"strings"

	"github.com/dihedron/steampipe-plugin-utils/utils"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
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
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Human-readable name for the volume. Might not be unique.",
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
				Description: "The id of the project the volume belongs to.",
				Transform:   transform.FromField("OsVolTenantAttrTenantID"),
			},
			{
				Name:        "user_id",
				Type:        proto.ColumnType_STRING,
				Description: "The id of the user who created the volume.",
				Transform:   transform.FromField("UserID"),
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "Indicates the current status of the volume.",
				Transform:   transform.FromField("Status"),
			},
			{
				Name:        "replication_status",
				Type:        proto.ColumnType_STRING,
				Description: "Indicates the status of replication of the volume.",
				Transform:   transform.FromField("ReplicationStatus"),
			},
			{
				Name:        "size",
				Type:        proto.ColumnType_INT,
				Description: "Size of the volume in GB.",
				Transform:   transform.FromField("Size"),
			},
			{
				Name:        "availability_zone",
				Type:        proto.ColumnType_STRING,
				Description: "AvailabilityZone is which availability zone the volume is in.",
				Transform:   transform.FromField("AvailabilityZone"),
			},
			{
				Name:        "bootable",
				Type:        proto.ColumnType_BOOL,
				Description: "Indicates whether this is a bootable volume.",
				Transform:   transform.FromField("Bootable"),
			},
			{
				Name:        "encrypted",
				Type:        proto.ColumnType_BOOL,
				Description: "Denotes if the volume is encrypted.",
				Transform:   transform.FromField("Encrypted"),
			},
			{
				Name:        "multiattach",
				Type:        proto.ColumnType_BOOL,
				Description: "denotes if the volume is multi-attach capable.",
				Transform:   transform.FromField("Multiattach"),
			},
			{
				Name:        "consistencygroup_id",
				Type:        proto.ColumnType_STRING,
				Description: "The volume's consistency group id.",
				Transform:   transform.FromField("ConsistencyGroupID"),
			},
			{
				Name:        "volume_type",
				Type:        proto.ColumnType_STRING,
				Description: "The type of volume to create, either SATA or SSD.",
				Transform:   transform.FromField("VolumeType"),
			},
			{
				Name:        "snapshot_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the snapshot from which the volume was created.",
				Transform:   transform.FromField("SnapshotID"),
			},
			{
				Name:        "source_vol_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of another block storage volume from which the current volume was created.",
				Transform:   transform.FromField("SourceVolID"),
			},
			{
				Name:        "backup_id",
				Type:        proto.ColumnType_STRING,
				Description: "The backup ID, from which the volume was restored; this value is available starting from microversion 3.47.",
				Transform:   transform.FromField("BackupID"),
			},
			{
				Name:        "group_id",
				Type:        proto.ColumnType_STRING,
				Description: "The group ID of the volume; this value is available starting from microversion 3.47.",
				Transform:   transform.FromField("GroupID"),
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_STRING,
				Description: "Timestamp when the volume was created.",
				Transform:   transform.FromField("CreatedAt").Transform(ToTime),
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_STRING,
				Description: "The date when this volume was last updated.",
				Transform:   transform.FromField("UpdatedAt").Transform(ToTime),
			},
			{
				Name:        "image_id",
				Type:        proto.ColumnType_STRING,
				Description: "The id of the image from which this volume was created, if any.",
				Transform: transform.FromField("VolumeImageMetadata").Transform(transform.NullIfZeroValue).Transform(func(ctx context.Context, d *transform.TransformData) (any, error) {
					if d.Value != nil {
						if value, ok := d.Value.(map[string]string); ok {
							return value["image_id"], nil
						}
					}
					return nil, nil
				}),
			},
			{
				Name:        "image_name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the image from which this volume was created, if any.",
				Transform: transform.FromField("VolumeImageMetadata").Transform(transform.NullIfZeroValue).Transform(func(ctx context.Context, d *transform.TransformData) (any, error) {
					if d.Value != nil {
						if value, ok := d.Value.(map[string]string); ok {
							return value["image_name"], nil
						}
					}
					return nil, nil
				}),
			},
			{
				Name:        "image_size",
				Type:        proto.ColumnType_INT,
				Description: "The size of the image from which this volume was created, if any.",
				Transform: transform.FromField("VolumeImageMetadata").Transform(transform.NullIfZeroValue).Transform(func(ctx context.Context, d *transform.TransformData) (any, error) {
					if d.Value != nil {
						if value, ok := d.Value.(map[string]string); ok {
							return value["size"], nil
						}
					}
					return nil, nil
				}),
			},
			{
				Name:        "image_architecture",
				Type:        proto.ColumnType_STRING,
				Description: "The architecture of the image from which this volume was created, if any.",
				Transform: transform.FromField("VolumeImageMetadata").Transform(transform.NullIfZeroValue).Transform(func(ctx context.Context, d *transform.TransformData) (any, error) {
					if d.Value != nil {
						if value, ok := d.Value.(map[string]string); ok {
							return value["architecture"], nil
						}
					}
					return nil, nil
				}),
			},
			{
				Name:        "image_checksum",
				Type:        proto.ColumnType_STRING,
				Description: "The checksum of the image from which this volume was created, if any.",
				Transform: transform.FromField("VolumeImageMetadata").Transform(transform.NullIfZeroValue).Transform(func(ctx context.Context, d *transform.TransformData) (any, error) {
					if d.Value != nil {
						if value, ok := d.Value.(map[string]string); ok {
							return value["checksum"], nil
						}
					}
					return nil, nil
				}),
			},
			{
				Name:        "image_container_format",
				Type:        proto.ColumnType_STRING,
				Description: "The container format of the image from which this volume was created, if any.",
				Transform: transform.FromField("VolumeImageMetadata").Transform(transform.NullIfZeroValue).Transform(func(ctx context.Context, d *transform.TransformData) (any, error) {
					if d.Value != nil {
						if value, ok := d.Value.(map[string]string); ok {
							return value["container_format"], nil
						}
					}
					return nil, nil
				}),
			},
			{
				Name:        "image_disk_format",
				Type:        proto.ColumnType_STRING,
				Description: "The disk format of the image from which this volume was created, if any.",
				Transform: transform.FromField("VolumeImageMetadata").Transform(transform.NullIfZeroValue).Transform(func(ctx context.Context, d *transform.TransformData) (any, error) {
					if d.Value != nil {
						if value, ok := d.Value.(map[string]string); ok {
							return value["disk_format"], nil
						}
					}
					return nil, nil
				}),
			},
			{
				Name:        "image_hw_disk_bus",
				Type:        proto.ColumnType_STRING,
				Description: "The hardware disk bus of the image from which this volume was created, if any.",
				Transform: transform.FromField("VolumeImageMetadata").Transform(transform.NullIfZeroValue).Transform(func(ctx context.Context, d *transform.TransformData) (any, error) {
					if d.Value != nil {
						if value, ok := d.Value.(map[string]string); ok {
							return value["hw_disk_bus"], nil
						}
					}
					return nil, nil
				}),
			},
			{
				Name:        "image_hw_qemu_guest_agent",
				Type:        proto.ColumnType_BOOL,
				Description: "Whether the hardware QEMU guest agent is installe",
				Transform: transform.FromField("VolumeImageMetadata").Transform(transform.NullIfZeroValue).Transform(func(ctx context.Context, d *transform.TransformData) (any, error) {
					if d.Value != nil {
						if value, ok := d.Value.(map[string]string); ok {
							return strings.ToLower(value["hw_qemu_guest_agent"]) == "yes", nil
						}
					}
					return nil, nil
				}),
			},
			{
				Name:        "image_hw_rng_model",
				Type:        proto.ColumnType_STRING,
				Description: "The hardware random number generator of the image from which this volume was created, if any.",
				Transform: transform.FromField("VolumeImageMetadata").Transform(transform.NullIfZeroValue).Transform(func(ctx context.Context, d *transform.TransformData) (any, error) {
					if d.Value != nil {
						if value, ok := d.Value.(map[string]string); ok {
							return value["hw_rng_model"], nil
						}
					}
					return nil, nil
				}),
			},
			{
				Name:        "image_hw_scsi_model",
				Type:        proto.ColumnType_STRING,
				Description: "The hardware SCSi model of the image from which this volume was created, if any.",
				Transform: transform.FromField("VolumeImageMetadata").Transform(transform.NullIfZeroValue).Transform(func(ctx context.Context, d *transform.TransformData) (any, error) {
					if d.Value != nil {
						if value, ok := d.Value.(map[string]string); ok {
							return value["hw_scsi_model"], nil
						}
					}
					return nil, nil
				}),
			},
			{
				Name:        "image_min_disk",
				Type:        proto.ColumnType_INT,
				Description: "The minimum disk size (in Gb) of the image from which this volume was created, if any.",
				Transform: transform.FromField("VolumeImageMetadata").Transform(transform.NullIfZeroValue).Transform(func(ctx context.Context, d *transform.TransformData) (any, error) {
					if d.Value != nil {
						if value, ok := d.Value.(map[string]string); ok {
							if value["image_min_disk"] == "" {
								return 0, nil
							}
							return strconv.Atoi(value["image_min_disk"])
						}
					}
					return nil, nil
				}),
			},
			{
				Name:        "image_min_ram",
				Type:        proto.ColumnType_INT,
				Description: "The minimum RAM size of the image from which this volume was created, if any.",
				Transform: transform.FromField("VolumeImageMetadata").Transform(transform.NullIfZeroValue).Transform(func(ctx context.Context, d *transform.TransformData) (any, error) {
					if d.Value != nil {
						if value, ok := d.Value.(map[string]string); ok {
							if value["min_ram"] == "" {
								return 0, nil
							}
							return strconv.Atoi(value["min_ram"])
						}
					}
					return nil, nil
				}),
			},
			{
				Name:        "image_os_distro",
				Type:        proto.ColumnType_STRING,
				Description: "The operating system distribution of the image from which this volume was created, if any.",
				Transform: transform.FromField("VolumeImageMetadata").Transform(transform.NullIfZeroValue).Transform(func(ctx context.Context, d *transform.TransformData) (any, error) {
					if d.Value != nil {
						if value, ok := d.Value.(map[string]string); ok {
							return value["os_distro"], nil
						}
					}
					return nil, nil
				}),
			},
			{
				Name:        "metadata",
				Type:        proto.ColumnType_JSON,
				Description: "The volume metadata.",
				Transform:   transform.FromField("Metadata"),
			},
			{
				Name:        "volume_image_metadata",
				Type:        proto.ColumnType_JSON,
				Description: "The volume image metadata.",
				Transform:   transform.FromField("VolumeImageMetadata"),
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

//// LIST FUNCTION

func listOpenStackVolume(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	plugin.Logger(ctx).Debug("retrieving openstack volumes list", "query data", utils.ToPrettyJSON(d))

	client, err := getServiceClient(ctx, d, BlockStorageV3)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	opts := buildOpenStackVolumeFilter(ctx, d.EqualsQuals)

	allPages, err := volumes.List(client, opts).AllPages()
	if err != nil {
		plugin.Logger(ctx).Error("error listing volumes with options", "options", utils.ToPrettyJSON(opts), "error", err)
		return nil, err
	}
	allVolumes := []*apiVolume{}
	err = volumes.ExtractVolumesInto(allPages, &allVolumes)
	if err != nil {
		plugin.Logger(ctx).Error("error extracting volumes", "error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("volumes retrieved", "count", len(allVolumes))

	for _, volume := range allVolumes {
		if ctx.Err() != nil {
			plugin.Logger(ctx).Debug("context done, exit")
			break
		}
		volume := volume
		d.StreamListItem(ctx, volume)
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOpenStackVolume(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	id := d.EqualsQuals["id"].GetStringValue()
	plugin.Logger(ctx).Debug("retrieving openstack volume", "id", id)

	client, err := getServiceClient(ctx, d, BlockStorageV3)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	result := volumes.Get(client, id)
	//plugin.Logger(ctx).Debug("request run", "result", utils.ToPrettyJSON(result))

	volume := &apiVolume{}
	err = result.ExtractInto(volume)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving volume", "error", err)
		return nil, err
	}

	return volume, nil
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
	plugin.Logger(ctx).Debug("returning", "filter", utils.ToPrettyJSON(opts))
	return opts
}

type apiVolume struct {
	// Unique identifier for the volume.
	ID string `json:"id"`
	// Current status of the volume.
	Status string `json:"status"`
	// Size of the volume in GB.
	Size int `json:"size"`
	// AvailabilityZone is which availability zone the volume is in.
	AvailabilityZone string `json:"availability_zone"`
	// The date when this volume was created.
	CreatedAt Time `json:"created_at"`
	// The date when this volume was last updated
	UpdatedAt Time `json:"updated_at"`
	// Instances onto which the volume is attached.
	Attachments []volumes.Attachment `json:"attachments"`
	// Human-readable display name for the volume.
	Name string `json:"name"`
	// Human-readable description for the volume.
	Description string `json:"description"`
	// The type of volume to create, either SATA or SSD.
	VolumeType string `json:"volume_type"`
	// The ID of the snapshot from which the volume was created
	SnapshotID string `json:"snapshot_id"`
	// The ID of another block storage volume from which the current volume was created
	SourceVolID string `json:"source_volid"`
	// The backup ID, from which the volume was restored
	// This field is supported since 3.47 microversion
	BackupID *string `json:"backup_id"`
	// The group ID; this field is supported since 3.47 microversion
	GroupID *string `json:"group_id"`
	// Arbitrary key-value pairs defined by the user.
	Metadata map[string]string `json:"metadata"`
	// UserID is the id of the user who created the volume.
	UserID string `json:"user_id"`
	// Indicates whether this is a bootable volume.
	Bootable string `json:"bootable"`
	// Encrypted denotes if the volume is encrypted.
	Encrypted bool `json:"encrypted"`
	// ReplicationStatus is the status of replication.
	ReplicationStatus string `json:"replication_status"`
	// ConsistencyGroupID is the consistency group ID.
	ConsistencyGroupID string `json:"consistencygroup_id"`
	// Multiattach denotes if the volume is multi-attach capable.
	Multiattach bool `json:"multiattach"`
	// Image metadata entries, only included for volumes that were created from an image, or from a snapshot of a volume originally created from an image.
	VolumeImageMetadata map[string]string `json:"volume_image_metadata"`
	// The volume migration status
	MigrationStatus string `json:"migration_status"`

	OsVolHostAttrHost         string `json:"os-vol-host-attr:host"`
	OsVolMigStatusAttrMigstat string `json:"os-vol-mig-status-attr:migstat"`
	OsVolMigStatusAttrNameID  string `json:"os-vol-mig-status-attr:name_id"`
	OsVolTenantAttrTenantID   string `json:"os-vol-tenant-attr:tenant_id"`
	ProviderID                string `json:"provider_id"`
	ServiceUUID               string `json:"service_uuid"`
	SharedTargets             bool   `json:"shared_targets"`
}
