connection "openstack" {
    # the path to the plugin
    plugin    = "local/openstack"
    # the OpenStack API endpoint; can also be set with the 
    # ... environment variable
    endpoint_url = "http://keystone.example.com:8080"
    userid = "<userid>"
    username = "<username>"
    password = "<password>"
    region = "<region>"
    project_id = "<project id>"
    project_name = "<project name>"
    domain_id = "<domain id>"
    domain_name = "<domain name>"
    # only use when loggin in via app credentials
    access_token = "<token id>"
    app_credential_id = "<application credential id>"
    app_credential_name = "<application credential id>"
    app_credential_secret = "application credential secret>"
    trace_level = "TRACE"
}