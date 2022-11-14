package openstack

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/hashicorp/go-hclog"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

const (
	COMPUTEv2_MICROVERSION  = "2.79"
	IDENTITYv3_MICROVERSION = "3.13"
)

var ErrNotImplemented = errors.New("not implemented")

func getComputeV2Client(ctx context.Context, d *plugin.QueryData) (*gophercloud.ServiceClient, error) {
	// load connection from cache, which preserves throttling protection etc
	cacheKey := "openstack_computev2"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		plugin.Logger(ctx).Debug("returning compute v2 client from cache")
		return cachedData.(*gophercloud.ServiceClient), nil
	}

	plugin.Logger(ctx).Info("creating new compute v2 client")
	api, err := getAuthenticatedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("no valid authenticated client available", "error", err)
		return nil, err
	}

	openstackConfig := GetConfig(d.Connection)
	region := ""
	if openstackConfig.Region != nil {
		region = *openstackConfig.Region
	}

	client, err := openstack.NewComputeV2(api, gophercloud.EndpointOpts{Region: region})

	if err != nil {
		plugin.Logger(ctx).Error("error creating compute v2 client", "error", err)
		return nil, err
	}
	// see https://docs.openstack.org/nova/latest/reference/api-microversion-history.html
	client.Microversion = COMPUTEv2_MICROVERSION

	// save to cache
	plugin.Logger(ctx).Debug("saving compute v2 client to cache")
	d.ConnectionManager.Cache.Set(cacheKey, client)

	return client, nil
}

func getIdentityV3Client(ctx context.Context, d *plugin.QueryData) (*gophercloud.ServiceClient, error) {
	// load connection from cache, which preserves throttling protection etc
	cacheKey := "openstack_identityv2"
	if cachedData, ok := d.ConnectionManager.Cache.Get(cacheKey); ok {
		plugin.Logger(ctx).Debug("returning identity v3 client from cache")
		return cachedData.(*gophercloud.ServiceClient), nil
	}

	plugin.Logger(ctx).Info("creating new identity v3 client")
	api, err := getAuthenticatedClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("no valid authenticated client available", "error", err)
		return nil, err
	}

	openstackConfig := GetConfig(d.Connection)
	region := ""
	if openstackConfig.Region != nil {
		region = *openstackConfig.Region
	}

	client, err := openstack.NewIdentityV3(api, gophercloud.EndpointOpts{Region: region})

	if err != nil {
		plugin.Logger(ctx).Error("error creating identity v3 client", "error", err)
		return nil, err
	}
	// see https://docs.openstack.org/nova/latest/reference/api-microversion-history.html
	client.Microversion = IDENTITYv3_MICROVERSION

	// save to cache
	plugin.Logger(ctx).Debug("saving compute v2 client to cache")
	d.ConnectionManager.Cache.Set(cacheKey, client)

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

func setLogLevel(ctx context.Context, d *plugin.QueryData) {
	openstackConfig := GetConfig(d.Connection)
	if openstackConfig.TraceLevel != nil {
		level := *openstackConfig.EndpointUrl
		plugin.Logger(ctx).SetLevel(hclog.LevelFromString(level))
	}
}

// ToPrettyJSON dumps the input object to JSON.
func toPrettyJSON(v any) string {
	s, _ := json.MarshalIndent(v, "", "  ")
	return string(s)
}
