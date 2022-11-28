package openstack

import (
	"context"
	"strings"

	"github.com/dihedron/steampipe-plugin-utils/utils"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableOpenStackImage(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openstack_image",
		Description: "OpenStack Disk Image",
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The unique id of the image.",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "Human-readable name for the image. Might not be unique.",
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "Indicates the current status of the image.",
				Transform:   transform.FromField("Status"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: "Tags is a list of image tags. Tags are arbitrarily defined strings attached to an image.",
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "container_format",
				Type:        proto.ColumnType_STRING,
				Description: "The container format of the image; valid values are ami, ari, aki, bare, and ovf.",
				Transform:   transform.FromField("ContainerFormat"),
			},
			{
				Name:        "disk_format",
				Type:        proto.ColumnType_STRING,
				Description: "The disk format of the image; if set, valid values are ami, ari, aki, vhd, vmdk, raw, qcow2, vdi, and iso.",
				Transform:   transform.FromField("DiskFormat"),
			},
			{
				Name:        "min_disk",
				Type:        proto.ColumnType_INT,
				Description: "This is the amount of disk space in GB that is required to boot the image.",
				Transform:   transform.FromField("MinDiskGigabytes").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "image_min_ram",
				Type:        proto.ColumnType_INT,
				Description: "This is the amount of RAM in MB that is required to boot the image.",
				Transform:   transform.FromField("MinRAMMegabytes").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "project_id",
				Type:        proto.ColumnType_STRING,
				Description: "The id of the project the image belongs to.",
				Transform:   transform.FromField("Owner"),
			},
			{
				Name:        "protected",
				Type:        proto.ColumnType_BOOL,
				Description: "Protected is whether the image is deletable or not.",
				Transform:   transform.FromField("Protected"),
			},
			{
				Name:        "visibility",
				Type:        proto.ColumnType_STRING,
				Description: " Visibility defines who can see/use the image.",
				Transform:   transform.FromField("Visibility"),
			},
			{
				Name:        "hidden",
				Type:        proto.ColumnType_BOOL,
				Description: "Hidden is whether the image is listed in default image list or not.",
				Transform:   transform.FromField("Hidden"),
			},
			{
				Name:        "checksum",
				Type:        proto.ColumnType_STRING,
				Description: "Checksum is the checksum of the data that's associated with the image.",
				Transform:   transform.FromField("Checksum"),
			},
			{
				Name:        "size",
				Type:        proto.ColumnType_INT,
				Description: "Size is the size of the data that's associated with the image.",
				Transform:   transform.FromField("SizeBytes").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "metadata",
				Type:        proto.ColumnType_JSON,
				Description: "Metadata is a set of metadata associated with the image.",
				Transform:   transform.FromField("Metadata").Transform(transform.NullIfZeroValue),
			},
			{
				Name:        "properties",
				Type:        proto.ColumnType_JSON,
				Description: "Properties is a set of key-value pairs, if any, that are associated with the image.",
				Transform:   transform.FromField("Properties").Transform(transform.NullIfZeroValue),
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
				Name:        "virtual_size",
				Type:        proto.ColumnType_INT,
				Description: "VirtualSize is the virtual size of the image.",
				Transform:   transform.FromField("VirtualSize").Transform(transform.NullIfZeroValue),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listOpenStackImage,
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
					Name:    "status",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "container_format",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "disk_format",
					Require: plugin.Optional,
				},
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getOpenStackImage,
		},
	}
}

//// LIST FUNCTION

func listOpenStackImage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	plugin.Logger(ctx).Debug("retrieving openstack images list", "query data", utils.ToPrettyJSON(d))

	client, err := getServiceClient(ctx, d, ImageServiceV2)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	opts := buildOpenStackImageFilter(ctx, d.EqualsQuals)

	allPages, err := images.List(client, opts).AllPages()
	if err != nil {
		plugin.Logger(ctx).Error("error listing images with options", "options", utils.ToPrettyJSON(opts), "error", err)
		return nil, err
	}
	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		plugin.Logger(ctx).Error("error extracting images", "error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("images retrieved", "count", len(allImages))

	for _, image := range allImages {
		if ctx.Err() != nil {
			plugin.Logger(ctx).Debug("context done, exit")
			break
		}
		image := image
		d.StreamListItem(ctx, image)
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOpenStackImage(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	id := d.EqualsQuals["id"].GetStringValue()
	plugin.Logger(ctx).Debug("retrieving openstack image", "id", id)

	client, err := getServiceClient(ctx, d, ImageServiceV2)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	result := images.Get(client, id)
	//plugin.Logger(ctx).Debug("request run", "result", utils.ToPrettyJSON(result))

	image := &images.Image{}
	err = result.ExtractInto(image)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving image", "error", err)
		return nil, err
	}

	return image, nil
}
func buildOpenStackImageFilter(ctx context.Context, quals plugin.KeyColumnEqualsQualMap) images.ListOpts {
	opts := images.ListOpts{}
	if value, ok := quals["id"]; ok {
		opts.ID = value.GetStringValue()
	}
	if value, ok := quals["name"]; ok {
		opts.Name = value.GetStringValue()
	}
	if value, ok := quals["status"]; ok {
		switch strings.ToLower(value.GetStringValue()) {
		case "queued":
			opts.Status = images.ImageStatusQueued
		case "saving":
			opts.Status = images.ImageStatusSaving
		case "active":
			opts.Status = images.ImageStatusActive
		case "killed":
			opts.Status = images.ImageStatusKilled
		case "deleted":
			opts.Status = images.ImageStatusDeleted
		case "pending_delete":
			opts.Status = images.ImageStatusPendingDelete
		case "deactivated":
			opts.Status = images.ImageStatusDeactivated
		case "importing":
			opts.Status = images.ImageStatusImporting
		}
	}
	if value, ok := quals["disk_format"]; ok {
		opts.DiskFormat = value.GetStringValue()
	}
	if value, ok := quals["container_format"]; ok {
		opts.ContainerFormat = value.GetStringValue()
	}
	plugin.Logger(ctx).Debug("returning", "filter", utils.ToPrettyJSON(opts))
	return opts
}
