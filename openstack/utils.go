package openstack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"gopkg.in/yaml.v3"
)

var ErrNotImplemented = errors.New("not implemented")

// FromStringerField picks a field (given its name) and then calls the type-specific
// String() method on the value.
func FromStringerField(field string) *transform.ColumnTransforms {
	return &transform.ColumnTransforms{
		Transforms: []*transform.TransformCall{
			{
				Transform: transform.FieldValue,
				Param:     field,
			},
			{
				Transform: AsString,
				Param:     nil,
			},
		},
	}
}

// AsString invokes the String() function on the hydrate item coming from transform data.
func AsString(ctx context.Context, d *transform.TransformData) (interface{}, error) {
	if t, ok := d.Value.(fmt.Stringer); ok {
		return t.String(), nil
	}
	return nil, fmt.Errorf("invalid data type: %T (%s)", d.Value, toJSON(d.Value))
}

// setLogLevel changes the current HCLog level; this seems necessary as the
// STEAMPIPE_LOG_LEVEL variable does not seem to be properly read by the plugins.
func setLogLevel(ctx context.Context, d *plugin.QueryData) {
	openstackConfig := GetConfig(d.Connection)
	if openstackConfig.TraceLevel != nil {
		level := *openstackConfig.TraceLevel
		plugin.Logger(ctx).SetLevel(hclog.LevelFromString(level))
	}
}

// toJSON dumps the input object to JSON.
func toJSON(v any) string {
	s, _ := json.Marshal(v)
	return string(s)
}

// toPrettyJSON dumps the input object to JSON.
func toPrettyJSON(v any) string {
	s, _ := json.MarshalIndent(v, "", "  ")
	return string(s)
}

// toYAML dumps the input object to YAML.
func toYAML(v any) string {
	s, _ := yaml.Marshal(v)
	return string(s)
}

// pointerTo returns a pointer to a given value.
func pointerTo[T any](value T) *T {
	return &value
}
