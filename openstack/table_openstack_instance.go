package openstack

import (
	"context"

	"github.com/dihedron/steampipe-plugin-utils/utils"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableOpenStackInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openstack_instance",
		Description: "OpenStack Virtual Machine Instance",
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Type:        proto.ColumnType_STRING,
				Description: "The instance id",
				Transform:   transform.FromField("ID"),
			},
			{
				Name:        "name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the instance",
				Transform:   transform.FromField("Name"),
			},
			{
				Name:        "description",
				Type:        proto.ColumnType_STRING,
				Description: "The description of the instance",
				Transform:   transform.FromField("Description"),
			},
			{
				Name:        "project_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the instance's project (aka tenant)",
				Transform:   transform.FromField("TenantID"),
			},
			{
				Name:        "user_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the instance's user",
				Transform:   transform.FromField("UserID"),
			},
			{
				Name:        "created_at",
				Type:        proto.ColumnType_STRING,
				Description: "The creation time of the instance",
				Transform:   transform.FromField("CreatedAt").Transform(ToTime),
			},
			{
				Name:        "launched_at",
				Type:        proto.ColumnType_STRING,
				Description: "The launch time of the instance",
				Transform:   transform.FromField("LaunchedAt").Transform(ToTime),
			},
			{
				Name:        "updated_at",
				Type:        proto.ColumnType_STRING,
				Description: "The update time of the instance",
				Transform:   transform.FromField("UpdatedAt").Transform(ToTime),
			},
			{
				Name:        "terminated_at",
				Type:        proto.ColumnType_STRING,
				Description: "The termination time of the instance",
				Transform:   transform.FromField("TerminatedAt").Transform(ToTime),
			},
			{
				Name:        "host_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the hypervisor (host) the instance is running on",
				Transform:   transform.FromField("HostID"),
			},
			{
				Name:        "addresses",
				Type:        proto.ColumnType_JSON,
				Description: "The IP address of the Instance",
				Transform: transform.FromField("Addresses").Transform(func(ctx context.Context, d *transform.TransformData) (any, error) {
					var results []map[string]string
					if value, ok := d.Value.(map[string][]struct {
						MACAddress string `json:"OS-EXT-IPS-MAC:mac_addr"`
						IPType     string `json:"OS-EXT-IPS:type"`
						IPAddress  string `json:"addr"`
						IPVersion  int    `json:"version"`
					}); ok {
						results = make([]map[string]string, 0, len(value)*2)
						for k, v := range value {
							ip := make(map[string]string, len(v))
							for _, a := range v {
								ip["Network"] = k
								ip["IPAddress"] = a.IPAddress
								ip["MACAddress"] = a.MACAddress
								results = append(results, ip)
							}
						}
						return results, nil
					}
					return results, nil
				}),
			},
			{
				Name:        "host_name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the host the instance is running on",
				Transform:   transform.FromField("HostName"),
			},
			{
				Name:        "availability_zone",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the hypervisor (host) the instance is running on",
				Transform:   transform.FromField("AvailabilityZone"),
			},
			{
				Name:        "status",
				Type:        proto.ColumnType_STRING,
				Description: "The status of the instance",
				Transform:   transform.FromField("Status"),
			},
			{
				Name:        "progress",
				Type:        proto.ColumnType_INT,
				Description: "Progress information about the instance.",
				Transform:   transform.FromField("Progress"),
			},
			{
				Name:        "hypervisor_hostname",
				Type:        proto.ColumnType_STRING,
				Description: "The hostname of the hypervisor on which the instance is running.",
				Transform:   transform.FromField("HypervisorHostname"),
			},
			{
				Name:        "tags",
				Type:        proto.ColumnType_JSON,
				Description: "The tags associated with the virtual machine.",
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "user_data",
				Type:        proto.ColumnType_STRING,
				Description: "The user data associated with the virtual machine.",
				Transform:   transform.FromField("UserData"),
			},
			{
				Name:        "disk_config",
				Type:        proto.ColumnType_STRING,
				Description: "The instance disc configuration.",
				Transform:   transform.FromField("DiskConfig"),
			},
			{
				Name:        "instance_name",
				Type:        proto.ColumnType_STRING,
				Description: "The instance name.",
				Transform:   transform.FromField("InstanceName"),
			},
			{
				Name:        "kernel_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the virtual machine kernel.",
				Transform:   transform.FromField("KernelID"),
			},
			{
				Name:        "launch_index",
				Type:        proto.ColumnType_STRING,
				Description: "The instance launch index.",
				Transform:   transform.FromField("LaunchIndex"),
			},
			{
				Name:        "ram_disk_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the instance.",
				Transform:   transform.FromField("RAMDiskID"),
			},
			{
				Name:        "reservation_id",
				Type:        proto.ColumnType_STRING,
				Description: "The instance reservation ID.",
				Transform:   transform.FromField("ReservationID"),
			},
			{
				Name:        "root_device_name",
				Type:        proto.ColumnType_STRING,
				Description: "The instance root device name.",
				Transform:   transform.FromField("RootDeviceName"),
			},
			{
				Name:        "power_state_id",
				Type:        proto.ColumnType_INT,
				Description: "The instance power state (as an integer).",
				Transform:   transform.FromField("PowerState"),
			},
			{
				Name:        "power_state_name",
				Type:        proto.ColumnType_STRING,
				Description: "The instance power state as a string.",
				Transform: transform.FromField("PowerState").Transform(func(ctx context.Context, d *transform.TransformData) (interface{}, error) {
					switch d.Value.(int) {
					case 0:
						return "NOSTATE", nil
					case 1:
						return "RUNNING", nil
					case 3:
						return "PAUSED", nil
					case 4:
						return "SHUTDOWN", nil
					case 6:
						return "CRASHED", nil
					case 7:
						return "SUSPENDED", nil
					}
					return "", nil
				}),
			},
			{
				Name:        "config_drive",
				Type:        proto.ColumnType_STRING,
				Description: "The instance config drive.",
				Transform:   transform.FromField("ConfigDrive"),
			},
			{
				Name:        "flavor_name",
				Type:        proto.ColumnType_STRING,
				Description: "The original name of the flavor used to start the instance.",
				Transform:   transform.FromField("Flavor.OriginalName"),
			},
			{
				Name:        "flavor_vcpus",
				Type:        proto.ColumnType_INT,
				Description: "The number of virtual CPUs in the flavor used to start the instance.",
				Transform:   transform.FromField("Flavor.VCPUs"),
			},
			{
				Name:        "flavor_vgpus",
				Type:        proto.ColumnType_INT,
				Description: "The number of virtual GPUs in the flavor used to start the instance.",
				Transform:   transform.FromField("Flavor.ExtraSpecs.VGPUs").Transform(transform.NullIfZeroValue).Transform(transform.ToInt),
			},
			{
				Name:        "flavor_cores",
				Type:        proto.ColumnType_INT,
				Description: "The number of virtual CPU cores in the flavor used to start the instance.",
				Transform:   transform.FromField("Flavor.ExtraSpecs.CPUCores").Transform(transform.NullIfZeroValue).Transform(transform.ToInt),
			},
			{
				Name:        "flavor_sockets",
				Type:        proto.ColumnType_INT,
				Description: "The number of CPU sockets in the flavor used to start the instance.",
				Transform:   transform.FromField("Flavor.ExtraSpecs.CPUSockets").Transform(transform.NullIfZeroValue).Transform(transform.ToInt),
			},
			{
				Name:        "flavor_ram",
				Type:        proto.ColumnType_INT,
				Description: "The amount of RAM in the flavor used to start the instance.",
				Transform:   transform.FromField("Flavor.RAM"),
			},
			{
				Name:        "flavor_disk",
				Type:        proto.ColumnType_INT,
				Description: "The size of the disk in the flavor used to start the instance.",
				Transform:   transform.FromField("Flavor.Disk"),
			},
			{
				Name:        "flavor_swap",
				Type:        proto.ColumnType_INT,
				Description: "The size of the swap disk in the flavor used to start the instance.",
				Transform:   transform.FromField("Flavor.Swap"),
			},
			{
				Name:        "flavor_ephemeral",
				Type:        proto.ColumnType_INT,
				Description: "The size of the ephemeral disk in the flavor used to start the instance.",
				Transform:   transform.FromField("Flavor.Ephemeral"),
			},
			{
				Name:        "flavor_rng_allowed",
				Type:        proto.ColumnType_BOOL,
				Description: "Whether the RNG is allowed on the flavor used to start the instance.",
				Transform:   transform.FromField("Flavor.ExtraSpecs.RNGAllowed").Transform(transform.NullIfZeroValue).Transform(transform.ToBool),
			},
			{
				Name:        "flavor_watchdog_action",
				Type:        proto.ColumnType_STRING,
				Description: "The action to take when the Nova watchdog detects the instance is not responding.",
				Transform:   transform.FromField("Flavor.ExtraSpecs.WatchdogAction"),
			},
			{
				Name:        "image_id",
				Type:        proto.ColumnType_STRING,
				Description: "The Glance image used to start the instance.",
				Transform: transform.FromField("Image").Transform(func(ctx context.Context, d *transform.TransformData) (any, error) {
					if d.Value != nil {
						if value, ok := d.Value.(map[string]any); ok {
							return value["id"], nil
						}
					}
					return nil, nil
				}),
			},
			{
				Name:        "attached_volume_ids",
				Type:        proto.ColumnType_JSON,
				Description: "The volumes attached to the instance.",
				Transform: transform.FromField("AttachedVolumes").Transform(func(ctx context.Context, d *transform.TransformData) (any, error) {
					if d.Value != nil {
						if volumes, ok := d.Value.([]servers.AttachedVolume); ok {
							result := []string{}
							for _, volume := range volumes {
								result = append(result, volume.ID)
							}
							return result, nil
						}
					}
					return nil, nil
				}),
			},
		},
		List: &plugin.ListConfig{
			Hydrate: listOpenStackInstance,
			KeyColumns: plugin.KeyColumnSlice{
				&plugin.KeyColumn{
					Name:    "name",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "host_id",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "status",
					Require: plugin.Optional,
				},
				// &plugin.KeyColumn{
				// 	Name: "image_name",
				// 	Require: plugin.Optional,
				// },
				&plugin.KeyColumn{
					Name:    "flavor_name",
					Require: plugin.Optional,
				},
				// &plugin.KeyColumn{
				// 	Name: "ipv4",
				// 	Require: plugin.Optional,
				// },
				// &plugin.KeyColumn{
				// 	Name: "ipv6",
				// 	Require: plugin.Optional,
				// },
				&plugin.KeyColumn{
					Name:    "project_id",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "user_id",
					Require: plugin.Optional,
				},
				&plugin.KeyColumn{
					Name:    "availability_zone",
					Require: plugin.Optional,
				},
				// TODO: add tags
			},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			Hydrate:    getOpenStackInstance,
		},
	}
}

/// LIST FUNCTION

func listOpenStackInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	plugin.Logger(ctx).Debug("retrieving openstack instance list", "query data", utils.ToPrettyJSON(d))

	client, err := getServiceClient(ctx, d, ComputeV2)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	opts := buildOpenStackInstanceFilter(ctx, d.EqualsQuals)

	allPages, err := servers.List(client, opts).AllPages()
	if err != nil {
		plugin.Logger(ctx).Error("error listing instances with options", "options", utils.ToPrettyJSON(opts), "error", err)
		return nil, err
	}
	allInstances := []*apiInstance{}
	err = servers.ExtractServersInto(allPages, &allInstances)
	if err != nil {
		plugin.Logger(ctx).Error("error extracting instances", "error", err)
		return nil, err
	}
	plugin.Logger(ctx).Debug("instances retrieved", "count", len(allInstances))

	for _, instance := range allInstances {
		if ctx.Err() != nil {
			plugin.Logger(ctx).Debug("context done, exit")
			break
		}
		instance := instance
		plugin.Logger(ctx).Debug("streaming instance", "data", utils.ToPrettyJSON(instance))
		d.StreamListItem(ctx, instance)
	}
	return nil, nil
}

//// HYDRATE FUNCTIONS

func getOpenStackInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	setLogLevel(ctx, d)

	id := d.EqualsQuals["id"].GetStringValue()
	plugin.Logger(ctx).Debug("retrieving openstack instance", "id", id)

	client, err := getServiceClient(ctx, d, ComputeV2)
	if err != nil {
		plugin.Logger(ctx).Error("error retrieving client", "error", err)
		return nil, err
	}

	result := servers.Get(client, id)
	plugin.Logger(ctx).Debug("API call complete", "result", utils.ToPrettyJSON(result))

	instance := &apiInstance{}
	if err := result.ExtractInto(instance); err != nil {
		plugin.Logger(ctx).Error("error retrieving instance", "error", err)
		return nil, err
	}

	plugin.Logger(ctx).Debug("returning instance", "data", utils.ToPrettyJSON(instance))
	return instance, nil
}

func buildOpenStackInstanceFilter(ctx context.Context, quals plugin.KeyColumnEqualsQualMap) servers.ListOpts {
	opts := servers.ListOpts{
		AllTenants: true,
	}

	if value, ok := quals["name"]; ok {
		opts.Name = value.GetStringValue()
	}
	if value, ok := quals["host_id"]; ok {
		opts.Host = value.GetStringValue()
	}
	if value, ok := quals["status"]; ok {
		opts.Status = value.GetStringValue()
	}
	// if value, ok := quals["ipv4_address"]; ok {
	// 	opts.IP = value.GetStringValue()
	// }
	// if value, ok := quals["ipv6_address"]; ok {
	// 	opts.IP6 = value.GetStringValue()
	// }
	if value, ok := quals["flavor_name"]; ok {
		opts.Flavor = value.GetStringValue()
	}
	// if value, ok := quals["image_name"]; ok {
	// 	opts.Image = value.GetStringValue()
	// }
	if value, ok := quals["project_id"]; ok {
		opts.TenantID = value.GetStringValue()
	}
	if value, ok := quals["user_id"]; ok {
		opts.UserID = value.GetStringValue()
	}
	if value, ok := quals["availability_zone"]; ok {
		opts.AvailabilityZone = value.GetStringValue()
	}
	plugin.Logger(ctx).Debug("returning", "filter", utils.ToPrettyJSON(opts))
	return opts
}

// apiInstance is an internal type used to unmarshal more datafrom the API
// response than would usually be possible through the ordinary gophercloud
// struct. OpenStack API microversions enable more response data that is not
// taken into account by the gophercloud library, which unmarshals only what
// is available at the base level for each API version, for backward compatibility.
// This is also why there is an ExtrctInto function that allows you to pass in
// an arbitrary struct to marshal the responsa data into.
type apiInstance struct {
	ID           string `json:"id"`
	TenantID     string `json:"tenant_id"`
	UserID       string `json:"user_id"`
	Name         string `json:"name"`
	CreatedAt    Time   `json:"created"`
	LaunchedAt   Time   `json:"OS-SRV-USG:launched_at"`
	UpdatedAt    Time   `json:"updated"`
	TerminatedAt Time   `json:"OS-SRV-USG:terminated_at"`
	HostID       string `json:"hostid"`
	Status       string `json:"status"`
	Progress     int    `json:"progress"`
	AccessIPv4   string `json:"accessIPv4"`
	AccessIPv6   string `json:"accessIPv6"`
	Image        any    `json:"image"`
	Flavor       struct {
		Disk       int `json:"disk"`
		Ephemeral  int `json:"ephemeral"`
		ExtraSpecs struct {
			CPUCores        string `json:"hw:cpu_cores"`
			CPUSockets      string `json:"hw:cpu_sockets"`
			RNGAllowed      string `json:"hw_rng:allowed"`
			WatchdogAction  string `json:"hw:watchdog_action"`
			VGPUs           string `json:"resources:VGPU"`
			TraitCustomVGPU string `json:"trait:CUSTOM_VGPU"`
		} `json:"extra_specs"`
		OriginalName string `json:"original_name"`
		RAM          int    `json:"ram"`
		Swap         int    `json:"swap"`
		VCPUs        int    `json:"vcpus"`
	} `json:"flavor"`
	Addresses map[string][]struct {
		MACAddress string `json:"OS-EXT-IPS-MAC:mac_addr"`
		IPType     string `json:"OS-EXT-IPS:type"`
		IPAddress  string `json:"addr"`
		IPVersion  int    `json:"version"`
	} `json:"addresses"`
	Metadata map[string]string `json:"metadata"`
	Links    []struct {
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"links"`
	KeyName        string `json:"key_name"`
	AdminPass      string `json:"adminPass"`
	SecurityGroups []struct {
		Name string `json:"name"`
	} `json:"security_groups"`
	AttachedVolumes []servers.AttachedVolume `json:"os-extended-volumes:volumes_attached"`
	// Fault              servers.Fault            `json:"fault"`
	Tags               *[]string `json:"tags"`
	ServerGroups       *[]string `json:"server_groups"`
	DiskConfig         string    `json:"OS-DCF:diskConfig"`
	AvailabilityZone   string    `json:"OS-EXT-AZ:availability_zone"`
	Host               string    `json:"OS-EXT-SRV-ATTR:host"`
	Hostname           string    `json:"OS-EXT-SRV-ATTR:hostname"`
	HypervisorHostname string    `json:"OS-EXT-SRV-ATTR:hypervisor_hostname"`
	InstanceName       string    `json:"OS-EXT-SRV-ATTR:instance_name"`
	KernelID           string    `json:"OS-EXT-SRV-ATTR:kernel_id"`
	LaunchIndex        int       `json:"OS-EXT-SRV-ATTR:launch_index"`
	RAMDiskID          string    `json:"OS-EXT-SRV-ATTR:ramdisk_id"`
	ReservationID      string    `json:"OS-EXT-SRV-ATTR:reservation_id"`
	RootDeviceName     string    `json:"OS-EXT-SRV-ATTR:root_device_name"`
	UserData           string    `json:"OS-EXT-SRV-ATTR:user_data"`
	PowerState         int       `json:"OS-EXT-STS:power_state"`
	VMState            string    `json:"OS-EXT-STS:vm_state"`
	ConfigDrive        string    `json:"config_drive"`
	Description        string    `json:"description"`
	//	TaskState          interface{}              `json:"OS-EXT-STS:task_state"`
}
