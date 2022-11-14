package openstack

import (
	"encoding/json"
	"testing"
)

func TestOpenStackTime(t *testing.T) {

	type Test struct {
		Time Time
	}

	var tests = []string{
		`{"Time": "2022-09-24T13:53:23Z"}`,
		`{"Time": "2022-10-11T14:17:48.000000"}`,
	}

	for _, test := range tests {
		a := Test{}
		err := json.Unmarshal([]byte(test), &a)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("after unmarshalling: %q", a.Time.String())
	}
}
