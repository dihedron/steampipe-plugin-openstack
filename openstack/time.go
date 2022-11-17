package openstack

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
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
