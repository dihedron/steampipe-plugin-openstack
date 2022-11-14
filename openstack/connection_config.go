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
	Region              *string `cty:"region"`
	ProjectID           *string `cty:"project_id"`
	ProjectName         *string `cty:"project_name"`
	DomainID            *string `cty:"domain_id"`
	DomainName          *string `cty:"domain_name"`
	AccessToken         *string `cty:"access_token"`
	AppCredentialID     *string `cty:"app_credential_id"`
	AppCredentialName   *string `cty:"app_credential_name"`
	AppCredentialSecret *string `cty:"app_credential_secret"`
	AllowReauth         *bool   `cty:"allow_reauth"`
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
	"region": {
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
	"app_credential_id": {
		Type: schema.TypeString,
	},
	"app_credential_name": {
		Type: schema.TypeString,
	},
	"app_credential_secret": {
		Type: schema.TypeString,
	},
	"allow_reauth": {
		Type: schema.TypeBool,
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
