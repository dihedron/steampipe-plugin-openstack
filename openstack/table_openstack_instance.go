package openstack

import (
	"context"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableOpenStackInstance(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "openstack_instance",
		Description: "OpenStack Virtual Machine Instance",
		List: &plugin.ListConfig{
			Hydrate: listOpenStackInstance,
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("id"),
			// IgnoreConfig: &plugin.IgnoreConfig{
			// 	ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidInstanceID.NotFound", "InvalidInstanceID.Unavailable", "InvalidInstanceID.Malformed"}),
			// },
			Hydrate: getOpenStackInstance,
		},
		Columns: []*plugin.Column{
			{Name: "id", Type: proto.ColumnType_STRING, Description: "The instance id"},
			{Name: "name", Type: proto.ColumnType_STRING, Description: "The name of the instance"},
			{Name: "project_id", Type: proto.ColumnType_STRING, Description: "The ID of the instance's project (aka tenant)"},
			{Name: "user_id", Type: proto.ColumnType_STRING, Description: "The ID of the instance's user"},
			{Name: "created_at", Type: proto.ColumnType_STRING, Description: "The creation time of the instance"},
			{Name: "launched_at", Type: proto.ColumnType_STRING, Description: "The launch time of the instance"},
			{Name: "updated_at", Type: proto.ColumnType_STRING, Description: "The update time of the instance"},
			{Name: "terminated_at", Type: proto.ColumnType_STRING, Description: "The termintaion time of the instance"},
			{Name: "host_id", Type: proto.ColumnType_STRING, Description: "The ID of the hypervisor (host) the instance is running on"},
			{Name: "status", Type: proto.ColumnType_STRING, Description: "The status of the instance"},
			{Name: "progress", Type: proto.ColumnType_INT, Description: "Progress information about the instance"},
			// AccessIPv4   string                 `json:"accessIPv4"`
			// AccessIPv6   string                 `json:"accessIPv6"`

		},
	}
}

//// LIST FUNCTION

func listOpenStackInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("retrieving openstack instance list")
	plugin.Logger(ctx).Debug("plugin query data: %s", toPrettyJSON(d))
	plugin.Logger(ctx).Debug("plugin hydrate data %s", toPrettyJSON(h))
	return nil, ErrNotImplemented
}

//// HYDRATE FUNCTIONS

func getOpenStackInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	id := d.KeyColumnQuals["id"].GetStringValue()
	plugin.Logger(ctx).Debug("retrieving openstack instance", "id", id)

	client, err := getComputeV2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("error creating compute v2 client", "error", err)
		return nil, err
	}

	result := servers.Get(client, id)
	instance := &apiInstance{}
	if err := result.ExtractInto(instance); err != nil {
		plugin.Logger(ctx).Error("error retrieving instance", "error", err)
		return nil, err
	}

	return &openstackInstance{
		ID:           id,
		Name:         instance.Name,
		ProjectID:    instance.TenantID,
		UserID:       instance.UserID,
		CreatedAt:    instance.CreatedAt.String(),
		LaunchedAt:   instance.LaunchedAt.String(),
		UpdatedAt:    instance.UpdatedAt.String(),
		TerminatedAt: instance.TerminatedAt.String(),
		HostID:       instance.HostID,
		Status:       instance.Status,
		Progress:     instance.Progress,
	}, nil
}

type apiInstance struct {
	ID           string                 `json:"id"`
	TenantID     string                 `json:"tenant_id"`
	UserID       string                 `json:"user_id"`
	Name         string                 `json:"name"`
	CreatedAt    Time                   `json:"created"`
	LaunchedAt   Time                   `json:"OS-SRV-USG:launched_at"`
	UpdatedAt    Time                   `json:"updated"`
	TerminatedAt Time                   `json:"OS-SRV-USG:terminated_at"`
	HostID       string                 `json:"hostid"`
	Status       string                 `json:"status"`
	Progress     int                    `json:"progress"`
	AccessIPv4   string                 `json:"accessIPv4"`
	AccessIPv6   string                 `json:"accessIPv6"`
	Image        map[string]interface{} `json:"-"`
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
	// AttachedVolumes    []servers.AttachedVolume `json:"os-extended-volumes:volumes_attached"`
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

type openstackInstance struct {
	ID           string
	Name         string
	ProjectID    string
	UserID       string
	CreatedAt    string
	LaunchedAt   string
	UpdatedAt    string
	TerminatedAt string
	HostID       string
	Status       string
	Progress     int
}
