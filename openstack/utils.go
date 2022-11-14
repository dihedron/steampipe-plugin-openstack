package openstack

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/hashicorp/go-hclog"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

const (
	DefaultComputeV2Microversion  = "2.79"
	DefaultIdentityV3Microversion = "3.13"
)

var ErrNotImplemented = errors.New("not implemented")

type serviceClientConfig struct {
	newClient       func(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error)
	getMicroversion func(config *openstackConfig) string
}

var serviceClientConfigs = map[string]serviceClientConfig{
	"openstack_identity_v3": {
		newClient: openstack.NewIdentityV3,
		getMicroversion: func(config *openstackConfig) string {
			microversion := DefaultIdentityV3Microversion
			if config.IdentityV3Microversion != nil {
				microversion = *config.IdentityV3Microversion
			}
			return microversion
		},
	},
	"openstack_compute_v2": {
		newClient: openstack.NewComputeV2,
		getMicroversion: func(config *openstackConfig) string {
			microversion := DefaultComputeV2Microversion
			if config.ComputeV2Microversion != nil {
				microversion = *config.ComputeV2Microversion
			}
			return microversion
		},
	},
	"openstack_network_v2": {
		newClient: openstack.NewNetworkV2,
		getMicroversion: func(config *openstackConfig) string {
			// no microversion support for networking
			return ""
		},
	},
}

func getServiceClient(ctx context.Context, d *plugin.QueryData, key string) (*gophercloud.ServiceClient, error) {
	plugin.Logger(ctx).Debug("returning service client", "type", key)

	// load connection from cache, which preserves throttling protection etc
	if cachedData, ok := d.ConnectionManager.Cache.Get(key); ok {
		plugin.Logger(ctx).Debug("returning service client from cache")
		return cachedData.(*gophercloud.ServiceClient), nil
	}

	if _, ok := serviceClientConfigs[key]; !ok {
		plugin.Logger(ctx).Error("invalid service client type", "type", key)
		panic(fmt.Sprintf("invalid service type: %q", key))
	}

	plugin.Logger(ctx).Info("creating new service client", "type", key)
	api, err := getAuthenticatedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("no valid authenticated provider client available", "error", err)
		return nil, err
	}

	openstackConfig := GetConfig(d.Connection)
	region := ""
	if openstackConfig.Region != nil {
		region = *openstackConfig.Region
	}

	client, err := serviceClientConfigs[key].newClient(api, gophercloud.EndpointOpts{Region: region})

	if err != nil {
		plugin.Logger(ctx).Error("error creating service client", "type", key, "error", err)
		return nil, err
	}
	client.Microversion = serviceClientConfigs[key].getMicroversion(&openstackConfig)

	// save to cache
	plugin.Logger(ctx).Debug("saving service client to cache", "type", key)
	d.ConnectionManager.Cache.Set(key, client)

	return client, nil
}

// Create the OpenStack REST API client.
func getAuthenticatedClient(ctx context.Context, d *plugin.QueryData) (*gophercloud.ProviderClient, error) {

	// load connection from cache, which preserves throttling protection etc
	cacheKey := "openstack_authenticated_client"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		plugin.Logger(ctx).Debug("returning the authenticated client from cache")
		return cachedData.(*gophercloud.ProviderClient), nil
	}

	plugin.Logger(ctx).Info("creating new authenticated client")

	// try with the environment first
	auth, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		plugin.Logger(ctx).Info("no auth info available in environment, filling with defaults", "error", err)

		// fill the auth info from the configuration
		auth.AllowReauth = true
		openstackConfig := GetConfig(d.Connection)
		if openstackConfig.EndpointUrl != nil {
			auth.IdentityEndpoint = *openstackConfig.EndpointUrl
		}
		if openstackConfig.UserID != nil {
			auth.UserID = *openstackConfig.UserID
		}
		if openstackConfig.Username != nil {
			auth.Username = *openstackConfig.Username
		}
		if openstackConfig.Password != nil {
			auth.Password = *openstackConfig.Password
		}
		if openstackConfig.ProjectID != nil {
			auth.TenantID = *openstackConfig.ProjectID
		}
		if openstackConfig.ProjectName != nil {
			auth.TenantName = *openstackConfig.ProjectName
		}
		if openstackConfig.DomainID != nil {
			auth.DomainID = *openstackConfig.DomainID
		}
		if openstackConfig.DomainName != nil {
			auth.DomainName = *openstackConfig.DomainName
		}
		if openstackConfig.AccessToken != nil {
			auth.TokenID = *openstackConfig.AccessToken
		}
		if openstackConfig.AppCredentialID != nil {
			auth.ApplicationCredentialID = *openstackConfig.AppCredentialID
		}
		if openstackConfig.AppCredentialSecret != nil {
			auth.IdentityEndpoint = *openstackConfig.AppCredentialSecret
		}
		if openstackConfig.AllowReauth != nil {
			auth.AllowReauth = *openstackConfig.AllowReauth
		}
	}

	client, err := openstack.AuthenticatedClient(auth)
	if err != nil {
		plugin.Logger(ctx).Error("error creating authenticated client", "error", err)
		return nil, err
	}

	// save to cache
	plugin.Logger(ctx).Debug("saving authenticated client to cache")
	d.ConnectionManager.Cache.Set(cacheKey, client)

	return client, nil
}

// setLogLevel changes the current HCLog level; this seems necessary as the
// STEAMPIPE_LOG_LEVEL variable does not seem to be properly read by the plugins.
func setLogLevel(ctx context.Context, d *plugin.QueryData) {
	openstackConfig := GetConfig(d.Connection)
	if openstackConfig.TraceLevel != nil {
		level := *openstackConfig.EndpointUrl
		plugin.Logger(ctx).SetLevel(hclog.LevelFromString(level))
	}
}

// toPrettyJSON dumps the input object to JSON.
func toPrettyJSON(v any) string {
	s, _ := json.MarshalIndent(v, "", "  ")
	return string(s)
}

// pointerTo returns a pointer to a given value.
func pointerTo[T any](value T) *T {
	return &value
}
