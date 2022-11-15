package openstack

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

const (
	// defaults currently referring to Train
	DefaultComputeV2Microversion      = "2.79"
	DefaultIdentityV3Microversion     = "3.13"
	DefaultBlockStorageV3Microversion = "3.59"
)

type ServiceType string

const (
	// AuthenticatedClient is the cache key for the openStack authenticated client.
	AuthenticatedClient = "openstack_authenticated_client"

	// IdentityV3 identifies the OpenStack Identity V3 service (Keystone).
	IdentityV3 ServiceType = "openstack_identity_v3"
	// Compute identifies the penStack Compute V2 service (Nova).
	ComputeV2 = "openstack_compute_v2"
	// NetworkV2 identifies the OpenStack Network V2 service (Neutron).
	NetworkV2 = "openstack_network_v2"
	// BlockStorageV3 identifies the OpenStack Block Storage V3 service (Cinder).
	BlockStorageV3 = "openstack_blockstorage_v3"
)

type serviceConfig struct {
	newClient       func(client *gophercloud.ProviderClient, eo gophercloud.EndpointOpts) (*gophercloud.ServiceClient, error)
	getMicroversion func(config *openstackConfig) string
}

var serviceConfigMap = map[ServiceType]serviceConfig{
	IdentityV3: {
		newClient: openstack.NewIdentityV3,
		getMicroversion: func(config *openstackConfig) string {
			microversion := DefaultIdentityV3Microversion
			if config.IdentityV3Microversion != nil {
				microversion = *config.IdentityV3Microversion
			}
			return microversion
		},
	},
	ComputeV2: {
		newClient: openstack.NewComputeV2,
		getMicroversion: func(config *openstackConfig) string {
			microversion := DefaultComputeV2Microversion
			if config.ComputeV2Microversion != nil {
				microversion = *config.ComputeV2Microversion
			}
			return microversion
		},
	},
	NetworkV2: {
		newClient: openstack.NewNetworkV2,
		getMicroversion: func(config *openstackConfig) string {
			// TODO: check if we need to leverage/support micro-versions
			return ""
		},
	},
	BlockStorageV3: {
		newClient: openstack.NewBlockStorageV3,
		getMicroversion: func(config *openstackConfig) string {
			microversion := DefaultBlockStorageV3Microversion
			if config.BlockStorageV3Microversion != nil {
				microversion = *config.BlockStorageV3Microversion
			}
			return microversion
		},
	},
}

func getServiceClient(ctx context.Context, d *plugin.QueryData, key ServiceType) (*gophercloud.ServiceClient, error) {
	plugin.Logger(ctx).Debug("returning service client", "type", key)

	// load connection from cache, which preserves throttling protection etc
	if cachedData, ok := d.ConnectionManager.Cache.Get(string(key)); ok {
		plugin.Logger(ctx).Debug("returning service client from cache")
		return cachedData.(*gophercloud.ServiceClient), nil
	}

	if _, ok := serviceConfigMap[key]; !ok {
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

	client, err := serviceConfigMap[key].newClient(api, gophercloud.EndpointOpts{Region: region})

	if err != nil {
		plugin.Logger(ctx).Error("error creating service client", "type", key, "error", err)
		return nil, err
	}
	client.Microversion = serviceConfigMap[key].getMicroversion(&openstackConfig)

	// save to cache
	plugin.Logger(ctx).Debug("saving service client to cache", "type", key)
	d.ConnectionManager.Cache.Set(string(key), client)

	return client, nil
}

// Create the OpenStack REST API client.
func getAuthenticatedClient(ctx context.Context, d *plugin.QueryData) (*gophercloud.ProviderClient, error) {

	// load connection from cache, which preserves throttling protection etc
	if cachedData, ok := d.ConnectionManager.Cache.Get(AuthenticatedClient); ok {
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
	d.ConnectionManager.Cache.Set(AuthenticatedClient, client)

	return client, nil
}
