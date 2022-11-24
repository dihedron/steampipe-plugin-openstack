package openstack

import (
	"context"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
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
				Name:        "created_at",
				Type:        proto.ColumnType_STRING,
				Description: "Timestamp when the volume was created.",
				Transform:   TransformFromTimeField("CreatedAt"),
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_STRING,
				Description: "The date when this volume was last updated.",
				Transform:   TransformFromTimeField("UpdatedAt"),
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
	allVolumes := []*apiVolume{}
	err = volumes.ExtractVolumesInto(allPages, &allVolumes)
	if err != nil {
		plugin.Logger(ctx).Error("error extracting volumes", "error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("volumes retrieved", "count", len(allVolumes))

	for _, volume := range allVolumes {
		volume := volume
		d.StreamListItem(ctx, volume)
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
	//plugin.Logger(ctx).Debug("request run", "result", toPrettyJSON(result))

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
	plugin.Logger(ctx).Debug("returning", "filter", toPrettyJSON(opts))
	return opts
}

type apiVolume struct {
	Attachments        []interface{} `json:"attachments"`
	AvailabilityZone   string        `json:"availability_zone"`
	Bootable           string        `json:"bootable"`
	ConsistencyGroupID string        `json:"consistencygroup_id"`
	CreatedAt          Time          `json:"created_at"`
	Description        string        `json:"description"`
	Encrypted          bool          `json:"encrypted"`
	GroupID            string        `json:"group_id"`
	BackupID           string        `json:"backup_id"`
	ID                 string        `json:"id"`
	Links              []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"links"`
	Metadata struct {
		Readonly string `json:"readonly"`
	} `json:"metadata"`
	MigrationStatus           string `json:"migration_status"`
	Multiattach               bool   `json:"multiattach"`
	Name                      string `json:"name"`
	OsVolHostAttrHost         string `json:"os-vol-host-attr:host"`
	OsVolMigStatusAttrMigstat string `json:"os-vol-mig-status-attr:migstat"`
	OsVolMigStatusAttrNameID  string `json:"os-vol-mig-status-attr:name_id"`
	OsVolTenantAttrTenantID   string `json:"os-vol-tenant-attr:tenant_id"`
	ProviderID                string `json:"provider_id"`
	ReplicationStatus         string `json:"replication_status"`
	ServiceUUID               string `json:"service_uuid"`
	SharedTargets             bool   `json:"shared_targets"`
	Size                      int    `json:"size"`
	SnapshotID                string `json:"snapshot_id"`
	SourceVolID               string `json:"source_volid"`
	Status                    string `json:"status"`
	UpdatedAt                 Time   `json:"updated_at"`
	UserID                    string `json:"user_id"`
	VolumeType                string `json:"volume_type"`
}
