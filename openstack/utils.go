package openstack

import (
	"context"
	"encoding/json"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

// ToPrettyJSON dumps the input object to JSON.
func ToPrettyJSON(v any) string {
	s, _ := json.MarshalIndent(v, "", "  ")
	return string(s)
}

// Create Rest API (v3) client
func connect(ctx context.Context, d *plugin.QueryData) *gophercloud.ProviderClient {

	// Load connection from cache, which preserves throttling protection etc
	cacheKey := "openstack"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		return cachedData.(*gophercloud.ProviderClient)
	}

	// TODO: do the real connection logic here
	// Get connection config for plugin
	// openstackConfig := GetConfig(d.Connection)
	// if openstackConfig.Token != nil {
	// 	token = *openstackConfig.Token
	// }
	// if openstackConfig.BaseURL != nil {
	// 	baseURL = *openstackConfig.BaseURL
	// }
	// token := os.Getenv("OPENSTACK_PROJECT")
	// baseURL := os.Getenv("OPENSTACK_ENDPOINT_URL")
	auth, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		panic("no authing available in environment")
	}

	provider, err := openstack.AuthenticatedClient(auth)
	if err != nil {
		panic("error creating authenticated client")
	}

	// save to cache
	d.ConnectionManager.Cache.Set(cacheKey, provider)

	return provider
}
