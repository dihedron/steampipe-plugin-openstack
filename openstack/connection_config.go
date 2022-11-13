package openstack

import (
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/schema"
)

type openstackConfig struct {
	EndpointUrl         *string `cty:"endpoint_url"`
	UserID              *string `cty:"userid"`
	Username            *string `cty:"username"`
	Password            *string `cty:"password"`
	ProjectID           *string `cty:"project_id"`
	ProjectName         *string `cty:"project_name"`
	DomainID            *string `cty:"domain_id"`
	DomainName          *string `cty:"domain_name"`
	AccessToken         *string `cty:"access_token"`
	AppCredentialID     *string `cty:"app_credential_key"`
	AppCredentialSecret *string `cty:"app_credential_secret"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"endpoint_url": {
		Type: schema.TypeString,
	},
	"userid": {
		Type: schema.TypeString,
	},
	"username": {
		Type: schema.TypeString,
	},
	"password": {
		Type: schema.TypeString,
	},
	"project_id": {
		Type: schema.TypeString,
	},
	"project_name": {
		Type: schema.TypeString,
	},
	"domain_id": {
		Type: schema.TypeString,
	},
	"domain_name": {
		Type: schema.TypeString,
	},
	"access_token": {
		Type: schema.TypeString,
	},
	"app_credential_key": {
		Type: schema.TypeString,
	},
	"app_credential_secret": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &openstackConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) openstackConfig {
	if connection == nil || connection.Config == nil {
		return openstackConfig{}
	}
	config, _ := connection.Config.(openstackConfig)

	return config
}