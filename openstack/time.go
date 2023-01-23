package openstack

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type Time time.Time

var layouts = []string{
	"2006-01-02T15:04:05Z0700",
	"2006-01-02T15:04:05.000000",
}

func (t *Time) Format(format string) string {
	if t != nil {
		return time.Time(*t).Format(format)
	}
	return ""
}

//const layout2 = "2006-01-02T15:04:05Z"
// CCYY-MM-DDThh:mm:ss±hh:mm

// updated: 2022-10-07T18:56:03Z
// create: 2022-10-07T18:55:42Z
// launched at: 2022-10-07T18:56:02.000000
// terminated at: 2022-10-07T18:56:02.000000

func (t *Time) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		*t = Time(time.Time{})
		return nil
	}
	var errors error
	for _, layout := range layouts {
		//var t0 time.Time
		if t0, err := time.Parse(layout, s); err != nil {
			errors = multierror.Append(errors, err)
			continue
		} else {
			*t = Time(t0)
			return nil
		}
	}
	return errors
}

func (t *Time) MarshalJSON() ([]byte, error) {
	if (time.Time(*t)).IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", (time.Time(*t)).Format(layouts[0]))), nil
}

func (t *Time) String() string {
	return (time.Time(*t)).String()
}

func (t *Time) IsZero() bool {
	if t == nil {
		return true
	}
	return time.Time(*t).IsZero()
}

func ToTime(ctx context.Context, d *transform.TransformData) (any, error) {
	var err error
	if d.Value == nil {
		return nil, nil
	}
	switch t := d.Value.(type) {
	case *Time:
		if t == nil || t.IsZero() {
			return nil, nil
		}
		return t.String(), nil
	case Time:
		if t.IsZero() {
			return nil, nil
		}
		return t.String(), nil
	case time.Time:
		if t.IsZero() {
			return nil, nil
		}
		return t.String(), nil
	default:
		err = fmt.Errorf("invalid type: %T", d.Value)
	}
	return nil, err
}
