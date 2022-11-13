.PHONY: plugin
plugin:
	@go build

.PHONY: clean
clean:
	@rm -rf steampipe-plugin-openstack

.PHONY: install
install: plugin
	@mkdir -p ~/.steampipe/plugins/local/openstack
	@cp steampipe-plugin-openstack ~/.steampipe/plugins/local/openstack/openstack.plugin
	@cp openstack.tf ~/.steampipe/config/openstack.spc