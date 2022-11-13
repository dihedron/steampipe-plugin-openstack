package openstack

import "encoding/json"

// ToPrettyJSON dumps the input object to JSON.
func ToPrettyJSON(v any) string {
	s, _ := json.MarshalIndent(v, "", "  ")
	return string(s)
}
