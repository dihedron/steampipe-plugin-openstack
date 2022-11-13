package openstack

import (
	"context"
	"errors"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

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
			{Name: "test", Type: proto.ColumnType_STRING, Description: "A test string to see how it works"},
		},
	}
}

func listOpenStackInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Debug("retrieving openstack instance list")
	plugin.Logger(ctx).Debug("plugin query data: %s", ToPrettyJSON(d))
	plugin.Logger(ctx).Debug("plugin hydrate data %s", ToPrettyJSON(h))
	return nil, ErrNotImplemented
}

func getOpenStackInstance(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// instanceID := d.KeyColumnQuals["id"].GetStringValue()
	// plugin.Logger(ctx).Debug("retrieving openstack instance %s", instanceID)
	plugin.Logger(ctx).Debug("retrieving openstack instance")
	plugin.Logger(ctx).Debug("plugin query data: %s", ToPrettyJSON(d))
	plugin.Logger(ctx).Debug("plugin hydrate data %s", ToPrettyJSON(h))

	// // create service
	// svc, err := EC2Client(ctx, d)
	// if err != nil {
	// 	plugin.Logger(ctx).Error("aws_ec2_instance.getEc2Instance", "connection_error", err)
	// 	return nil, err
	// }

	// params := &ec2.DescribeInstancesInput{
	// 	InstanceIds: []string{instanceID},
	// }

	// op, err := svc.DescribeInstances(ctx, params)
	// if err != nil {
	// 	plugin.Logger(ctx).Error("aws_ec2_instance.getEc2Instance", "api_error", err)
	// 	return nil, err
	// }

	// if op.Reservations != nil && len(op.Reservations) > 0 {
	// 	if op.Reservations[0].Instances != nil && len(op.Reservations[0].Instances) > 0 {
	// 		return op.Reservations[0].Instances[0], nil
	// 	}
	// }
	return nil, ErrNotImplemented
}

var ErrNotImplemented = errors.New("not implemented")
