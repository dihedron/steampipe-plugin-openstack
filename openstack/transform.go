package openstack

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func TransformFromStringToInt(field string) *transform.ColumnTransforms {
	return &transform.ColumnTransforms{
		Transforms: []*transform.TransformCall{
			{
				Transform: transform.FieldValue,
				Param:     field,
			},
			{
				Transform: func(ctx context.Context, d *transform.TransformData) (interface{}, error) {
					if t, ok := d.Value.(string); ok {
						if t == "" {
							return 0, nil
						}
						if v, err := strconv.Atoi(t); err != nil {
							return nil, fmt.Errorf("error converting %q to int: %w", t, err)
						} else {
							return v, nil
						}
					}
					return 0, nil
				},
				Param: nil,
			},
		},
	}
}

func TransformFromTimeField(field string) *transform.ColumnTransforms {
	return &transform.ColumnTransforms{
		Transforms: []*transform.TransformCall{
			{
				Transform: transform.FieldValue,
				Param:     field,
			},
			{
				Transform: func(ctx context.Context, d *transform.TransformData) (interface{}, error) {
					var err error
					switch t := d.Value.(type) {
					case Time:
						if !t.IsZero() {
							return t.String(), nil
						} else {
							return nil, nil
						}
					case time.Time:
						if !t.IsZero() {
							return t.String(), nil
						} else {
							return nil, nil
						}
					default:
						err = fmt.Errorf("invalid type: %T", d.Value)
					}
					return nil, err
				},
				Param: nil,
			},
		},
	}
}
