package openstack

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func Plugin(ctx context.Context) *plugin.Plugin {
	p := &plugin.Plugin{
		Name:             "steampipe-plugin-openstack",
		DefaultTransform: transform.FromGo().NullIfZero(),
		TableMap: map[string]*plugin.Table{
			"openstack_instance": tableOpenStackInstance(ctx),
			// "openstack_project":           tableOpenStackProject(ctx),
			// "openstack_volume":            tableOpenStackVolume(ctx),
			// "openstack_port":              tableOpenStackPort(ctx),
			// "openstack_attachment":        tableOpenStackAttachment(ctx),
			// "openstack_securitygroup":     tableOpenStackSecurityGroup(ctx),
			// "openstack_securitygrouprule": tableOpenStackSecurityGroupRule(ctx),
		},
		ConnectionConfigSchema: &plugin.ConnectionConfigSchema{
			NewInstance: ConfigInstance,
			Schema:      ConfigSchema,
		},
	}
	return p
}